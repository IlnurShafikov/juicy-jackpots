package manager

import (
	"context"
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/step"
)

// EnterGame handles the game entry/load request from the client.
// It retrieves the game configuration for the specified RTP version and creates
// the initial game state including rules, initial board, and configuration.
//
// Parameters:
//   - ctx: context for request cancellation and tracing (currently unused)
//   - betValue: the player's initial bet amount
//   - rtpVersion: the RTP version to load (e.g., 96 for 96%)
//
// Returns:
//   - pointer to LoadGameResult containing game rules and initial state
//   - error if RTP version is invalid or result generation fails
func (gm *GameManager) EnterGame(
	_ context.Context, betValue uint64, rtpVersion int,
) (*step.LoadGameResult, error) {
	gameMods, err := gm.getGameByRTP(rtpVersion)
	if err != nil {
		return nil, fmt.Errorf("get game by rtp version '%d': %w", rtpVersion, err)
	}

	resp, err := datatransfer.GetEnterResult(betValue, gameMods)
	if err != nil {
		return nil, fmt.Errorf("get enter result: %w", err)
	}

	return resp, nil
}
