package datatransfer

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result/events"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/rules"
	baseDto "github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/base"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/common"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/step"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

// TransformBaseStep transforms a base game result into a client-facing step.
// It converts the internal base result into a DTO and adds it to the steps collection.
//
// Parameters:
//   - step: the base game result to transform
//   - cost: the cost of this spin
//   - bet: the player's bet amount
//   - allCliSteps: accumulated list of client steps
//
// Returns:
//   - updated slice of client steps
//   - total win amount (including this step)
//   - error if transformation fails
func (t *Transfer) TransformBaseStep(
	step *result.Base, allCliSteps []*step.Step, cost, bet uint64,
) ([]*step.Step, uint64, error) {
	clientStep, err := buildBaseStep(step, cost, bet, 0, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("build base step: %w", err)
	}

	data, err := jsoniter.Marshal(clientStep)
	if err != nil {
		return nil, 0, fmt.Errorf("marshal base step: %w", err)
	}

	id, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, 0, fmt.Errorf("marshal step id: %w", err)
	}

	allCliSteps = append(allCliSteps, makeBaseClientStep(id, data, cost, clientStep, 0))

	return allCliSteps, clientStep.TotalWin, nil
}

// buildBaseStep converts an internal base result into a base DTO step.
// This function orchestrates the transformation of all game state components
// (symbols, events, wins) into a structured DTO that will be serialized and sent to the client.
//
// Parameters:
//   - step: the internal base game result containing all game state
//   - cost: the cost of this spin (usually equal to bet)
//   - betValue: the player's bet amount used to calculate monetary values
//   - initFsCount: number of free spins triggered by this step (0 if none)
//   - initRsCount: number of respins triggered by this step (0 if none)
//
// Returns:
//   - pointer to baseDto.Step containing all transformed data
//   - error if any transformation fails
func buildBaseStep(
	step *result.Base, cost, betValue uint64, initFsCount, initRsCount uint32,
) (*baseDto.Step, error) {
	// Transform the initial board state before any modifications
	symbols, err := TransformSymbols(step.Board, betValue)
	if err != nil {
		return nil, fmt.Errorf("transform before symbols: %w", err)
	}

	// Extract and transform all game events (wins, triggers, features, etc.)
	eventList, err := ExtractAndTransformEvents(step, betValue)
	if err != nil {
		return nil, fmt.Errorf("transform events: %w", err)
	}

	// Calculate the total win by multiplying base result by bet value
	totalWin := step.TotalWin() * betValue

	// Build the payload containing all game state information
	payload := buildBasePayload(eventList, symbols, nil, initFsCount, initRsCount)

	// Construct and return the complete step DTO
	return &baseDto.Step{
		Payload:       payload,
		TotalWin:      totalWin,
		StepWin:       totalWin, // For base spins, step win equals total win
		Cost:          cost,
		Type:          rules.Step_TYPE_BASE_SPIN,
		MaxWinReached: step.MaxWinReached(),
	}, nil
}

// buildBasePayload constructs the payload for a base game step.
// The payload contains all game state data that will be sent to the client,
// including board symbols, events, and feature counters.
//
// Parameters:
//   - eventList: all game events that occurred during this spin
//   - symbols: the board state before game logic was applied
//   - afterSymbols: the board state after game logic was applied
//   - freeSpinsAdded: number of free spins triggered (0 if none)
//   - initRsCount: number of respins triggered (0 if none)
//
// Returns:
//   - baseDto.Payload with all game state data populated
func buildBasePayload(
	eventList []events.Event, symbols, afterSymbols []*common.Symbol,
	freeSpinsAdded, initRsCount uint32,
) baseDto.Payload {
	// Initialize the payload with core game data
	payload := baseDto.Payload{
		Events:       eventList,
		Symbols:      symbols,
		AfterSymbols: afterSymbols,
		NextStepType: nextStepTypeFromBase(freeSpinsAdded, initRsCount),
		Counter:      nil, // Will be populated only if free spins were triggered
	}

	// Only add counter if free spins were triggered
	// Counter tracks how many spins were added, remain, and have been played
	if freeSpinsAdded > 0 {
		payload.Counter = &common.Counter{
			Added:  freeSpinsAdded, // Number of free spins just triggered
			Left:   0,              // Will be updated by the free spin handler
			Played: 0,              // No spins played yet at this point
		}
	}

	return payload
}

// makeBaseClientStep constructs the protobuf Step message for the base spin.
// This function wraps the serialized DTO data into a protobuf message that conforms
// to the game API protocol and can be sent to the client.
//
// Parameters:
//   - id: unique binary identifier for this step (UUID)
//   - data: JSON-serialized baseDto.Step containing all game state
//   - cost: the cost of this spin
//   - clientStep: the DTO step (used to extract win amounts and metadata)
//   - stepsAdded: number of bonus steps triggered (free spins or respins)
//
// Returns:
//   - pointer to gameV1.Step protobuf message ready to be sent to client
func makeBaseClientStep(
	id, data []byte, cost uint64, clientStep *baseDto.Step, stepsAdded uint32,
) *step.Step {
	return &step.Step{
		Id:             id,                       // Unique step identifier
		ClientResponse: data,                     // JSON payload with complete game state
		Cost:           cost,                     // How much this spin cost the player
		StepWin:        clientStep.StepWin,       // Win amount for this specific step
		TotalWin:       clientStep.TotalWin,      // Cumulative win across all steps
		Type:           clientStep.Type,          // Step type (TYPE_BASE_SPIN)
		StepsAdded:     stepsAdded,               // Number of bonus steps triggered
		MaxWinReached:  clientStep.MaxWinReached, // Whether max win limit was hit
	}
}

// nextStepTypeFromBase determines what type of step follows the base spin.
// The client uses this information to know what UI state to transition to
// and what type of step to expect next in the sequence.
//
// Logic priority:
//  1. If free spins were triggered → next step is a free spin
//  2. Else if respins were triggered → next step is a respin
//  3. Otherwise → next step is another base spin (round complete)
//
// Parameters:
//   - freeSpinsCount: number of free spins triggered by this base spin
//   - reSpinCount: number of respins triggered by this base spin
//
// Returns:
//   - TYPE_FREE_SPIN if free spins were triggered
//   - TYPE_RE_SPIN if respins were triggered
//   - TYPE_BASE_SPIN if neither feature was triggered (normal flow)
func nextStepTypeFromBase(freeSpinsCount, reSpinCount uint32) int {
	switch {
	case freeSpinsCount > 0:
		return rules.Step_TYPE_FREE_SPIN
	case reSpinCount > 0:
		return rules.Step_TYPE_RE_SPIN
	default:
		return rules.Step_TYPE_BASE_SPIN
	}
}
