package result

// BaseGame accumulates statistics for the base game phase across a simulation run.
type BaseGame struct {
	// TotalWins is the sum of all win amounts from base game spins.
	TotalWins uint64
	// LineWins is the cumulative win amount attributed to line wins only.
	LineWins      uint64
	CoinWins      uint64
	CoinWinsCount uint64
}

func newBaseGame() *BaseGame {
	return &BaseGame{
		TotalWins: 0,
	}
}

// AddWins adds the given win amount to the base game win total.
func (b *BaseGame) AddWins(allWins uint64) {
	b.TotalWins += allWins
}

// AddLineWins adds the given amount to the cumulative line-win total.
//
// Parameters:
//   - lineWins: line win amount from the current round
func (b *BaseGame) AddLineWins(lineWins uint64) {
	b.LineWins += lineWins
}

func (b *BaseGame) AddCoinWins(coinWins uint64) {
	b.CoinWins += coinWins
	b.CoinWinsCount++
}

// GetWins returns the total win amount accumulated from base game spins.
func (b *BaseGame) GetWins() uint64 {
	return b.TotalWins
}

// Merge adds another BaseGame's statistics into this one.
func (b *BaseGame) Merge(other *BaseGame) {
	b.TotalWins += other.TotalWins
	b.CoinWins += other.CoinWins
	b.LineWins += other.LineWins
	b.CoinWinsCount += other.CoinWinsCount
}
