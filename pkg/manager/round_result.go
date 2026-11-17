package manager

import (
	"context"
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/step"
)

// GetRoundResult is the main entry point for processing a game round request.
// It validates the request parameters, generates the game result, transforms it
// into a client-facing format, and returns the complete bet result.
//
// This method is typically called by the API layer when a player initiates a spin.
//
// Parameters:
//   - ctx: context for request cancellation and tracing (currently unused)
//   - params: request parameters (mod ID, feature ID, bet, RNG seed, RTP version)
//
// Returns:
//   - pointer to BetResult containing all game steps and results
//   - error if validation fails, game execution fails, or transformation fails
func (gm *GameManager) GetRoundResult(
	_ context.Context,
	modID, featureID string,
	rtpVersion int,
	seed int64,
	betValue uint64,
) (*step.BetResult, error) {
	result, err := gm.getRoundResult(modID, featureID, rtpVersion, seed, betValue)
	if err != nil {
		return nil, fmt.Errorf("get result: %w", err)
	}

	return result, nil
}
