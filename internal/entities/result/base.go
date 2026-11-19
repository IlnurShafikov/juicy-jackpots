// Package result contains data structures for representing game results
// across different game phases (base, free spins, respins).
package result

import (
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/coin"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result/events"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// Base represents the result of a base game spin.
// It contains all information about the spin outcome including
// events, wins, and triggered features.
type Base struct {
	// Board stores the game board configuration/state after the spin
	Board model.Board

	// CoinsResult holds the jackpot coins placement outcome for this spin.
	CoinsResult *coin.Result

	// events store all events that occurred during this spin
	events []events.Event

	// totalWin is the total win amount from this base spin (excluding respins/free spins)
	totalWin uint64

	// maxWinReached indicates if the maximum win limit was reached
	maxWinReached bool
}

// NewBase creates a new Base result with default values.
//
// Returns:
//   - pointer to a new Base instance
func NewBase() *Base {
	return &Base{
		Board:         nil,
		CoinsResult:   nil,
		events:        nil,
		totalWin:      0,
		maxWinReached: false,
	}
}

// SetEvents sets the events that occurred during this spin.
//
// Parameters:
//   - e: slice of events to set
func (b *Base) SetEvents(e []events.Event) {
	b.events = e
}

// SetTotalWin sets the total win amount for this spin.
//
// Parameters:
//   - win: total win amount
func (b *Base) SetTotalWin(win uint64) {
	b.totalWin = win
}

// SetMaxWinReached sets whether the maximum win limit was reached.
//
// Parameters:
//   - reached: true if max win was reached, false otherwise
func (b *Base) SetMaxWinReached(reached bool) {
	b.maxWinReached = reached
}

// Events returns all events that occurred during this spin.
//
// Returns:
//   - slice of events
func (b *Base) Events() []events.Event {
	return b.events
}

// TotalWin returns the total win amount for this spin.
//
// Returns:
//   - total win amount
func (b *Base) TotalWin() uint64 {
	return b.totalWin
}

// MaxWinReached returns whether the maximum win limit was reached.
//
// Returns:
//   - true if max win was reached, false otherwise
func (b *Base) MaxWinReached() bool {
	return b.maxWinReached
}
