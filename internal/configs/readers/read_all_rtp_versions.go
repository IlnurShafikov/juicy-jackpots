package readers

import (
	"errors"

	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
)

const (
	// RTP96 represents the 96% RTP version
	RTP96 = 96
	RTP94 = 94
	RTP92 = 92
	RTP90 = 90
)

// getSupportedRTPVersions returns a list of all supported RTP versions.
//
// Returns:
//   - slice of integers representing supported RTP percentages
func getSupportedRTPVersions() []int {
	return []int{
		RTP96, RTP94, RTP92, RTP90,
	}
}

// ReadGameRtpVersions reads and returns all available RTP versions of the game.
// It loads configurations for each supported RTP version and creates game mods for them.
//
// Returns:
//   - RtpGameVersions containing all configured RTP versions
//   - error if reading or parsing any RTP version fails
func ReadGameRtpVersions() (game.RtpGameVersions, error) {
	rtpVersions := getSupportedRTPVersions()
	gameVersions := make([]*game.Mods, len(rtpVersions))

	for i, rtpVersion := range rtpVersions {
		gameMods, err := ReadGameMods(rtpVersion)
		if err != nil {
			return nil, errors.Join(ErrReadGameConfigs, err)
		}

		gameVersions[i] = gameMods
	}

	return gameVersions, nil
}
