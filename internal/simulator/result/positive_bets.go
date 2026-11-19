package result

// PositiveBets counts rounds where the player won anything.
// Kept internal to simulator/result so it can expose AddN for parallel merge.
type PositiveBets struct {
	// spins is the number of rounds with a win.
	spins int
}

// newPositiveBets creates a PositiveBets counter set to zero.
func newPositiveBets() *PositiveBets { return &PositiveBets{spins: 0} }

// Add records one positive-bet round.
func (p *PositiveBets) Add() { p.spins++ }

// AddN merges another worker's positive-bet count into this one.
func (p *PositiveBets) AddN(n int) { p.spins += n }

// GetSpins returns the total number of positive-bet rounds.
func (p *PositiveBets) GetSpins() int { return p.spins }

// CalculateHitRate returns the fraction of rounds that had a win.
func (p *PositiveBets) CalculateHitRate(totalSpinsCount int) float64 {
	return float64(p.spins) / float64(totalSpinsCount)
}
