//nolint:all // it's ok for enter step
package datatransfer

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/rules"
	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/base"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/common"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/step"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	jsoniter "github.com/json-iterator/go"
)

// GetEnterResult creates the initial game load result sent to the client.
// This includes the game rules, initial board state, and configuration.
//
// Parameters:
//   - betValue: the player's bet amount
//   - gameMods: the game modes configuration for the selected RTP
//
// Returns:
//   - pointer to LoadGameResult containing rules and initial state
//   - error if construction fails
func GetEnterResult(betValue uint64, gameMods *game.Mods) (*step.LoadGameResult, error) {
	simpleCostMultiplier := gameMods.Simple.CostMultiplier()
	roundCost := betValue * simpleCostMultiplier

	clientEnterStep, err := enterStep(roundCost, betValue)
	if err != nil {
		return nil, fmt.Errorf("get enter step: %w", err)
	}

	jsonStep, err := jsoniter.ConfigFastest.Marshal(clientEnterStep)
	if err != nil {
		return nil, fmt.Errorf("marshal enter step: %w", err)
	}

	resp := &step.LoadGameResult{
		GameRules:      gameMods.Rules,
		Step:           clientInitialStep(jsonStep, roundCost),
		CostMultiplier: simpleCostMultiplier,
	}

	return resp, nil
}

// enterStep constructs the initial step shown when the player loads the game.
//
// Parameters:
//   - cost: the cost of a spin at the current bet level
//   - betValue: the player's bet amount
//
// Returns:
//   - pointer to base Step DTO representing the initial state
//   - error if construction fails
func enterStep(cost, betValue uint64) (*base.Step, error) {
	clientSymbols, err := TransformSymbols(initBoard(), betValue)
	if err != nil {
		return nil, fmt.Errorf("transfer symbols: %w", err)
	}

	resp := initialStep(cost, clientSymbols)

	return resp, err
}

// initialStep creates a Step DTO for the initial game state.
//
// Parameters:
//   - cost: the cost of a spin
//   - clientSymbols: the initial symbol board to display
//
// Returns:
//   - pointer to base Step DTO
func initialStep(cost uint64, clientSymbols []*common.Symbol) *base.Step {
	return &base.Step{
		Payload:       initialPayload(clientSymbols),
		TotalWin:      0,
		StepWin:       0,
		Cost:          cost,
		Type:          rules.Step_TYPE_BASE_SPIN,
		MaxWinReached: false,
	}
}

// clientInitialStep creates a protobuf Step message for the initial state.
//
// Parameters:
//   - jsonStep: JSON-serialized step data
//   - roundCost: the cost of a spin
//
// Returns:
//   - pointer to protobuf Step message
func clientInitialStep(jsonStep []byte, roundCost uint64) *step.Step {
	return &step.Step{
		Id:             []byte(""), // empty ID for initial step (not for analytics)
		ClientResponse: jsonStep,
		Cost:           roundCost,
		StepWin:        0,
		TotalWin:       0,
		Type:           rules.Step_TYPE_BASE_SPIN,
		StepsAdded:     0,
		MaxWinReached:  false,
	}
}

// initialPayload creates the payload for the initial game state.
//
// Parameters:
//   - clientSymbols: the initial symbol board
//
// Returns:
//   - base Payload DTO
func initialPayload(clientSymbols []*common.Symbol) base.Payload {
	return base.Payload{
		Counter:      nil,
		Events:       nil,
		Symbols:      clientSymbols,
		NextStepType: rules.Step_TYPE_BASE_SPIN,
	}
}

func initBoard() model.Board {
	return model.Board{
		{rules.Ace, rules.High1, rules.High2},
		{rules.Jack, rules.King, rules.Med1},
		{rules.Ace, rules.Queen, rules.Med2},
	}
}
