package events

import "github.com/IlnurShafikov/juicy-jackpots/internal/entities/coin"

// JackpotWinEvent represents the event emitted when a jackpot coin is collected.
// It carries the winning coin result for delivery to the client.
type JackpotWinEvent struct {
	// Target holds the coin result containing the jackpot value and type.
	Target *coin.Result `json:"target"`
	// Type identifies this event as a JackpotWin event.
	Type Type `json:"type"`
}

// NewJackpotWinEvent creates a JackpotWinEvent for the given coin result.
//
// Parameters:
//   - target: coin result containing the jackpot value and type to report
//
// Returns:
//   - pointer to a new JackpotWinEvent
func NewJackpotWinEvent(target *coin.Result) *JackpotWinEvent {
	return &JackpotWinEvent{
		Type:   JackpotWinType,
		Target: target,
	}
}

// EventType returns the event type identifier, satisfying the [Event] interface.
//
// Returns:
//   - JackpotWinType constant
func (c *JackpotWinEvent) EventType() Type {
	return c.Type
}
