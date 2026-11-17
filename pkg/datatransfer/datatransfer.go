package datatransfer

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/step"
)

// Transfer is responsible for transforming internal game results into client-facing DTOs.
// It handles the conversion of all game phases (base, free, respins) into API responses.
type Transfer struct {
	// costMultiplier is the base cost multiplier for the game mode
	costMultiplier uint64
}

// NewTransfer creates a new Transfer instance with the specified cost multiplier.
//
// Parameters:
//   - costMultiplier: the base cost multiplier for game rounds
//
// Returns:
//   - pointer to a new Transfer instance
func NewTransfer(
	costMultiplier uint64,
) *Transfer {
	resp := &Transfer{
		costMultiplier: costMultiplier,
	}

	return resp
}

// Transfer converts a complete round result into a BetResult for the API.
// It orchestrates the transformation of base spin, free spins, and respins into client steps.
//
// Parameters:
//   - res: the complete round result from the game engine
//   - betValue: the player's bet amount
//
// Returns:
//   - pointer to BetResult containing all transformed steps
//   - error if transformation fails
func (t *Transfer) Transfer(res *result.Round, betValue uint64) (*step.BetResult, error) {
	var err error

	// Pre-allocate slice with estimated capacity for all steps
	allCliSteps := make([]*step.Step, 0, 1)
	roundCost := res.CostMultiplier * betValue

	// Transform the base game spin
	allCliSteps, _, err = t.TransformBaseStep(
		res.BasePhase, allCliSteps, roundCost, betValue)
	if err != nil {
		return nil, fmt.Errorf("transfer base step: %w", err)
	}

	// Construct the final bet result
	clientRoundResp := &step.BetResult{
		Steps:          allCliSteps,
		UserBet:        betValue,
		CostMultiplier: t.costMultiplier,
		Seed:           res.Seed,
	}

	return clientRoundResp, nil
}
