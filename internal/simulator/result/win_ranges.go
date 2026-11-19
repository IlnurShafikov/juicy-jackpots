//nolint:all // it's ok for simulator
package result

import (
	"fmt"
	"strconv"
	"strings"
)

// RangeData holds the aggregated data for a single win range bucket.
type RangeData struct {
	// TotalSum is the sum of all wins that fell into this range.
	TotalSum uint64
	// Range is the win range this data belongs to.
	Range WinRange
	// Count is the number of rounds that produced a win within this range.
	Count int
}

// WinRanges tracks how simulation wins are distributed across predefined multiplier buckets.
// Both wins and their counts are stored per bucket for RTP-analysis reports.
//
// Hot-path optimizations (AddWinRange is called once per simulated round):
//  1. Bucket lookup uses binary search O(log N) instead of linear scan O(N).
//  2. Sums and counts are plain slices indexed by bucket position, eliminating
//     the per-call string hashing that map[WinRange] keys would incur.
type WinRanges struct {
	sums   []uint64 // cumulative win sum per bucket, index matches rangesList
	counts []int    // number of wins per bucket, index matches rangesList
	// rangesList holds the buckets sorted by lower bound.
	rangesList []WinRange
}

// WinRange defines a single win-multiplier bucket used in simulation reports.
type WinRange struct {
	// Name is the human-readable label, e.g. "[1x-2x)".
	Name string
	// MinBound is the inclusive lower boundary in base units.
	MinBound uint64
	// MaxBound is the exclusive upper boundary in base units.
	MaxBound uint64
}

// NewWinRanges creates a WinRanges instance with buckets derived from costMultiplier.
// The bucket boundaries are expressed as multiples of the cost multiplier so that
// the same relative ranges apply regardless of the round price.
//
// Parameters:
//   - costMultiplier: the cost of a single round; used to scale range boundaries
//
// Returns:
//   - pointer to an initialised WinRanges instance
func NewWinRanges(costMultiplier int) *WinRanges {
	rangesList := generateWinRanges(costMultiplier)
	n := len(rangesList)

	return &WinRanges{
		sums:       make([]uint64, n),
		counts:     make([]int, n),
		rangesList: rangesList,
	}
}

// GetSortedRanges returns all range buckets sorted from the smallest to the largest boundary.
// Each element contains the range definition, the cumulative win sum, and the count.
//
// Returns:
//   - slice of RangeData ordered by ascending MinBound
func (w *WinRanges) GetSortedRanges() []RangeData {
	result := make([]RangeData, 0, len(w.rangesList))

	for i, r := range w.rangesList {
		result = append(result, RangeData{
			Range:    r,
			TotalSum: w.sums[i],
			Count:    w.counts[i],
		})
	}

	return result
}

// AddWinRange adds the given win amount to the appropriate bucket.
// If the win does not fall into any known bucket, it is silently ignored.
//
// Parameters:
//   - win: the win amount in base units (0 is recorded in the zero-win bucket)
func (w *WinRanges) AddWinRange(win uint64) {
	idx := w.findWinRangeIdx(win)
	if idx >= 0 {
		w.sums[idx] += win
		w.counts[idx]++
	}
}

// Merge adds another WinRanges' bucket counts and sums into this one.
// Both instances must have been created with the same costMultiplier.
func (w *WinRanges) Merge(other *WinRanges) {
	for i := range w.sums {
		w.sums[i] += other.sums[i]
		w.counts[i] += other.counts[i]
	}
}

// findWinRangeIdx returns the bucket index for win using binary search, or -1.
// rangesList[0] is the zero-win bucket (MinBound == MaxBound == 0); all other
// buckets are sorted by ascending MinBound, making binary search correct.
func (w *WinRanges) findWinRangeIdx(win uint64) int {
	if win == 0 {
		if len(w.rangesList) > 0 {
			return 0
		}

		return -1
	}

	// Binary search for the last bucket where MinBound <= win.
	// Bucket 0 is the zero-win bucket; start from index 1.
	lo, hi := 1, len(w.rangesList)-1

	for lo <= hi {
		mid := (lo + hi) / 2
		if w.rangesList[mid].MinBound <= win {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	if hi >= 1 && win < w.rangesList[hi].MaxBound {
		return hi
	}

	return -1
}

// generateWinRanges builds the win buckets scaled by the cost multiplier.
// The first bucket captures zero wins, the rest follow the predefined multiplier bounds.
//
// Parameters:
//   - costMultiplier: cost multiplier of the game mode
//
// Returns:
//   - list of buckets sorted by lower bound
func generateWinRanges(costMultiplier int) []WinRange {
	// The first bucket is special: it captures the exact zero-win case.
	ranges := []WinRange{
		{
			MinBound: 0,
			MaxBound: 0,
			Name:     "[0x-0x)",
		},
	}

	// Remaining buckets start at 0.01x instead of 0.05x for finer granularity at low wins.
	baseBounds := []float64{
		0.01, 0.05, 0.1, 0.2, 0.4, 0.6, 0.8, 1, 1.5, 2, 2.5, 5, 10,
		25, 50, 100, 200, 300, 400, 500, 1000, 1500, 2000, 3000, 4000, 5000, 7000, 10000, 15000,
	}

	for i := 0; i < len(baseBounds)-1; i++ {
		ranges = append(ranges, WinRange{
			MinBound: uint64(baseBounds[i] * float64(costMultiplier)),
			MaxBound: uint64(baseBounds[i+1] * float64(costMultiplier)),
			Name:     formatRangeName(baseBounds[i], baseBounds[i+1]),
		})
	}

	return ranges
}

// formatRangeName returns the label of a bucket, e.g. "[1x-2x)".
func formatRangeName(minimum, maximum float64) string {
	return fmt.Sprintf("[%vx-%vx)", formatMultiplier(minimum), formatMultiplier(maximum))
}

// formatMultiplier formats a multiplier without trailing zeros, e.g. 2.50 becomes "2.5".
func formatMultiplier(val float64) string {
	if val == float64(int(val)) {
		return strconv.Itoa(int(val))
	}

	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", val), "0"), ".")
}
