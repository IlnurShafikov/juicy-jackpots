package simulator

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/configs/readers"
	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
)

// GameModFinder is responsible for locating a specific game mod by RTP version,
// mod ID, and feature ID. It reads all available game configurations on each call.
type GameModFinder struct{}

// NewGameModFinder creates and returns a new GameModFinder instance.
func NewGameModFinder() *GameModFinder {
	return &GameModFinder{}
}

// FindGameMod loads all RTP game versions and returns the game mod that matches
// the given RTP percentage, mod ID, and feature ID.
//
// Parameters:
//   - rtpPercent: the RTP version to look up (e.g. 96 for 96%)
//   - modID: the game mod identifier (e.g. "simple")
//   - featureID: the feature identifier (empty string means no feature)
//
// Returns:
//   - pointer to the matching game.Mod
//   - ErrGameNotFound if no game version matches the given RTP
//   - error if configuration loading or mod lookup fails
func (gmf *GameModFinder) FindGameMod(
	rtpPercent int, modID, featureID string,
) (*game.Mod, error) {
	gameMods, err := readers.ReadGameRtpVersions()
	if err != nil {
		return nil, fmt.Errorf("factory.MakeGame: %w", err)
	}

	gameModsByRTP := gmf.findGameModsByRTP(gameMods, rtpPercent)
	if gameModsByRTP == nil {
		return nil, ErrGameNotFound
	}

	gameMod, err := gameModsByRTP.GetGameMod(modID, featureID)
	if err != nil {
		return nil, fmt.Errorf("find game mod: %w", err)
	}

	return gameMod, nil
}

// findGameModsByRTP searches the provided list of game mods for one matching rtpPercent.
// Returns nil if no match is found.
func (gmf *GameModFinder) findGameModsByRTP(
	gameMods []*game.Mods, rtpPercent int,
) *game.Mods {
	for _, gameMod := range gameMods {
		if gameMod.GetRTPVersion() == rtpPercent {
			return gameMod
		}
	}

	return nil
}
