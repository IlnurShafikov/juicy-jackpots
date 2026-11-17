// Package free contains data transfer objects (DTOs) for free spins phase.
// These structures represent free spin results sent to the client.
package free

// Step represents a complete free spin step for client communication.
type Step struct {
	// Payload contains all the game data (symbols, events, counters)
	Payload Payload `json:"payload"`

	// TotalWin is the cumulative win amount up to and including this step
	TotalWin uint64 `json:"totalWin"`

	// StepWin is the win amount for this specific free spin
	StepWin uint64 `json:"stepWin"`

	// Type indicates the step type (should be TYPE_FREE_SPIN)
	Type int `json:"type"`

	// MaxWinReached indicates if the maximum win limit was reached
	MaxWinReached bool `json:"maxWinReached"`
}
