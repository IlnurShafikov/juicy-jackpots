// Package result contains data structures that accumulate simulation statistics
// across all simulated rounds, including win totals, hit rates, RTP metrics,
// and win-range distributions for both base and free game phases.
package result

// Result holds the aggregated statistics collected during a simulation run.
// It is populated by the parser package as each round result is processed.
type Result struct {
	// BaseWinRanges tracks the distribution of base-phase total wins.
	BaseWinRanges *WinRanges

	base *BaseGame

	positiveBets  *PositiveBets
	xi            *Xi
	maxTotalWin   uint64
	maxWinReached uint64
}

// New creates a new Result with all counters zeroed and win-range buckets initialised.
func New(roundCost uint64, simpleCostMultiplier int) *Result {
	return &Result{
		base:          newBaseGame(),
		positiveBets:  newPositiveBets(),
		xi:            newXi(roundCost),
		maxTotalWin:   0,
		maxWinReached: 0,
		BaseWinRanges: NewWinRanges(simpleCostMultiplier),
	}
}

// Merge incorporates another Result into this one.
// Used by the parallel runner to combine per-worker accumulators.
func (r *Result) Merge(other *Result) {
	r.base.Merge(other.base)
	r.positiveBets.AddN(other.positiveBets.GetSpins())
	r.xi.AddSum(other.xi.Sum())
	r.maxWinReached += other.maxWinReached

	if other.maxTotalWin > r.maxTotalWin {
		r.maxTotalWin = other.maxTotalWin
	}

	r.BaseWinRanges.Merge(other.BaseWinRanges)
}

// Base returns the base game statistics collector.
func (r *Result) Base() *BaseGame { return r.base }

// Xi returns the Xi statistics tracker.
func (r *Result) Xi() *Xi { return r.xi }

// MaxTotalWin returns the maximum total win across all rounds.
func (r *Result) MaxTotalWin() uint64 { return r.maxTotalWin }

// PositiveBets returns the positive bets statistics collector.
func (r *Result) PositiveBets() *PositiveBets { return r.positiveBets }

// AddMaxTotalWin updates the maximum total win if the provided win amount is greater.
func (r *Result) AddMaxTotalWin(win uint64) {
	if r.maxTotalWin < win {
		r.maxTotalWin = win
	}
}

// AddMaxWinReached increments the counter for maximum wins reached.
func (r *Result) AddMaxWinReached() {
	r.maxWinReached++
}

// MaxWinReached returns the total count of maximum wins reached.
func (r *Result) MaxWinReached() uint64 {
	return r.maxWinReached
}
