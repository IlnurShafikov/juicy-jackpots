// Package manager contains the components that turn game results into wins.
package manager

import (
	"github.com/IlnurShafikov/juicy-jackpots/internal/components/collectors"
	"github.com/IlnurShafikov/juicy-jackpots/internal/components/limiter"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/coin"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result/events"
)

// JpCoinWin handles the win calculation and event emission for the jackpot coins of a spin.
type JpCoinWin struct{}

// NewJpCoinWin creates a JpCoinWin manager.
//
// Returns:
//   - pointer to a new JpCoinWin
func NewJpCoinWin() *JpCoinWin { return &JpCoinWin{} }

// Calculate applies the jackpot coin win to the win limiter and emits a JackpotWin event.
// Returns zero immediately when the coin carries no value.
//
// Parameters:
//   - coinsRes: coin result holding the jackpot value and type to process
//   - winLimiter: win cap tracker updated with the jackpot value
//   - eventer: event collector that receives the JackpotWin event
//
// Returns:
//   - jackpot win value, or 0 if coinsRes.Value is zero
func (c *JpCoinWin) Calculate(
	coinsRes *coin.Result,
	winLimiter limiter.WinLimiter,
	eventer *collectors.EventsCollector,
) uint64 {
	winValue := coinsRes.Value

	if winValue == 0 {
		return 0
	}

	cutWin, _ := winLimiter.Add(winValue)

	eventer.Add(events.NewJackpotWinEvent(coinsRes.Copy()))

	return cutWin
}
