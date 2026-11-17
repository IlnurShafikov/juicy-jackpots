package result

// BaseGame accumulates statistics for the base game phase across a simulation run.
type BaseGame struct {
	// Wins is the sum of all win amounts from base game spins.
	Wins uint64
}

func newBaseGame() *BaseGame {
	return &BaseGame{
		Wins: 0,
	}
}

// AddWins adds the given win amount to the base game win total.
func (b *BaseGame) AddWins(waysWins uint64) {
	b.Wins += waysWins
}

// GetWins returns the total win amount accumulated from base game spins.
func (b *BaseGame) GetWins() uint64 {
	return b.Wins
}

// Merge adds another BaseGame's statistics into this one.
func (b *BaseGame) Merge(other *BaseGame) {
	b.Wins += other.Wins
}
