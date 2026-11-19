package result

import (
	"math"

	"github.com/IlnurShafikov/slot-engine/pkg/tools/simtools"
)

// xiSquaring is the power to which every RTP value is raised for the variance.
const xiSquaring = 2

// Xi accumulates squared RTP values for volatility calculation.
// Kept internal to simulator/result so it can expose Sum for parallel merge.
type Xi struct {
	// roundCost is the cost of one round, used to convert a win to RTP.
	roundCost uint64
	// sum is the accumulated sum of squared RTP values.
	sum float64
}

// newXi creates an Xi accumulator for the given round cost.
func newXi(roundCost uint64) *Xi {
	return &Xi{roundCost: roundCost, sum: 0}
}

// Add records a single win into the xi accumulator.
func (x *Xi) Add(win uint64) {
	x.sum += math.Pow(simtools.CalculateRTP(x.roundCost, win), xiSquaring)
}

// Sd returns the standard deviation using the accumulated xi sum.
func (x *Xi) Sd(totalRTP, iterationsCount float64) float64 {
	return simtools.CalculateStandardDeviation(totalRTP, x.sum, iterationsCount)
}

// Sum returns the raw accumulated sum for parallel merge.
func (x *Xi) Sum() float64 { return x.sum }

// AddSum merges another worker's xi sum into this one.
func (x *Xi) AddSum(s float64) { x.sum += s }
