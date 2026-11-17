package manager

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/configs/readers"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools/rngtest"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/step"
)

// GameManager is the central coordinator for game operations.
// It manages game configurations, RNG, result generation, and data transformation
// for all supported RTP versions and game modes.
type GameManager struct {
	// dt handles transformation of game results into client-facing DTOs
	dt *datatransfer.Transfer

	// rngPool provides pooled random number generators for deterministic gameplay
	rngPool *rngtest.Mt19937

	// rtpGames holds all game configurations for different RTP versions
	rtpGames game.RtpGameVersions
}

// NewGameManager creates and initializes a new GameManager instance.
// It loads all game configurations, initializes the RNG pool, and validates API errors.
//

// Returns:
//   - pointer to initialized GameManager
//   - error if configuration loading
func NewGameManager() (*GameManager, error) {
	rtpGameVersions, err := readers.ReadGameRtpVersions()
	if err != nil {
		return nil, fmt.Errorf("read rtp game versions: %w", err)
	}

	dataTransfer := datatransfer.NewTransfer(
		rtpGameVersions[0].Simple.CostMultiplier())

	res := &GameManager{
		rtpGames: rtpGameVersions,
		rngPool:  rngtest.New(),
		dt:       dataTransfer,
	}

	return res, nil
}

// getRoundResult generates and transforms a complete game round result.
// This is an internal method that coordinates result generation and transformation.
//
// Parameters:
//   - params: request parameters including mod ID, bet value, RNG seed, and RTP version
//
// Returns:
//   - pointer to BetResult ready for API response
//   - error if generation or transformation fails
func (gm *GameManager) getRoundResult(
	modID, featureID string,
	rtpVersion int,
	seed int64,
	betValue uint64,
) (*step.BetResult, error) {
	roundResult, err := gm.GenerateRoundResult(rtpVersion, modID, featureID, seed)
	if err != nil {
		return nil, fmt.Errorf("generate round result '%s': %w", modID, err)
	}

	resp, err := gm.dt.Transfer(roundResult, betValue)
	if err != nil {
		return nil, fmt.Errorf("transfer round result '%s': %w", modID, err)
	}

	return resp, nil
}

// GenerateRoundResult executes a game round using the specified parameters.
// It retrieves the appropriate game mode, gets an RNG from the pool,
// executes the game logic, and returns the internal result.
//
// Parameters:
//   - params: request parameters including mod ID, bet value, RNG seed, and RTP version
//
// Returns:
//   - pointer to Round containing the complete game result
//   - error if game mode retrieval, RNG acquisition, or spin execution fails
func (gm *GameManager) GenerateRoundResult(
	rtpVersion int,
	modID, featureID string,
	seed int64,
) (*result.Round, error) {
	gameMod, err := gm.GetGameMod(rtpVersion, modID, featureID)
	if err != nil {
		return nil, err
	}

	gm.rngPool.Reseed(seed)

	roundResult, err := gameMod.Spin(gm.rngPool)
	if err != nil {
		return nil, fmt.Errorf("get result:%s %w", modID, err)
	}

	roundResult.Seed = seed

	return roundResult, nil
}

// GetGameMod retrieves the appropriate game mode based on request parameters.
// It first finds the correct RTP version, then retrieves the specific game mode.
// Internal errors are wrapped into API-friendly errors before returning.
//
// Parameters:
//   - mod ID, feature ID, and RTP version
//
// Returns:
//   - pointer to the requested game Mod
//   - error if RTP version or game mode is not found
func (gm *GameManager) GetGameMod(rtpversion int, modID, featureID string) (*game.Mod, error) {
	gameMods, err := gm.getGameByRTP(rtpversion)
	if err != nil {
		return nil, err
	}

	gameMod, err := gameMods.GetGameMod(modID, featureID)
	if err != nil {
		return nil, fmt.Errorf("get game mod: %w", err)
	}

	return gameMod, nil
}

// getGameByRTP finds the game configuration for a specific RTP version.
// It searches through all loaded RTP versions to find a match.
//
// Parameters:
//   - rtpVersion: the RTP version to retrieve (e.g., 96 for 96%)
//
// Returns:
//   - pointer to Mods containing all game modes for this RTP
//   - ErrInvalidRTPVersion if the requested RTP version is not found
func (gm *GameManager) getGameByRTP(rtpVersion int) (*game.Mods, error) {
	for i := range gm.rtpGames {
		rtpGame := gm.rtpGames[i]

		if rtpGame.GetRTPVersion() == rtpVersion {
			return rtpGame, nil
		}
	}

	return nil, fmt.Errorf("rtp version %d: %w", rtpVersion, ErrInvalidRTPVersion)
}
