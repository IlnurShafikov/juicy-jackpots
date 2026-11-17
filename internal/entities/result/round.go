// Package result contains structures for representing complete game round results
// that are used to communicate between game logic and the API layer.
package result

// Round represents the complete result of a game round.
// It includes the base spin, any triggered free spins, and round metadata.
type Round struct {
	// BasePhase contains the result of the base game spin
	BasePhase *Base

	// CostMultiplier is the cost multiplier for this round, not including the player's bet level
	CostMultiplier uint64

	// Seed is the RNG seed used for this round (for replay and verification purposes)
	Seed int64
}

// NewRound creates a new Round result with the specified cost multiplier.
//
// Parameters:
//   - costMultiplier: the cost multiplier for this round
//
// Returns:
//   - pointer to a new Round instance
func NewRound(
	costMultiplier uint64,
) *Round {
	return &Round{
		BasePhase:      nil,
		CostMultiplier: costMultiplier,
		Seed:           0,
	}
}
