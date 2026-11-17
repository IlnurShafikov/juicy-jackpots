// Package base contains data transfer objects (DTOs) for base game phase.
// These structures represent the base spin result sent to the client.
package base

// Step represents a complete base game spin step for client communication.
type Step struct {
	// Payload contains all the game data (symbols, events, counters)
	Payload Payload `json:"payload"`

	// TotalWin is the cumulative win amount up to and including this step
	TotalWin uint64 `json:"totalWin"`

	// StepWin is the win amount for this specific step
	StepWin uint64 `json:"stepWin"`

	// Cost is the cost of this spin
	Cost uint64 `json:"cost"`

	// Type indicates the step type (base spin, free spin, etc.)
	Type int `json:"type"`

	// MaxWinReached indicates if the maximum win limit was reached
	MaxWinReached bool `json:"maxWinReached"`
}
