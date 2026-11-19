// Package rngtest provides a deterministic Mersenne Twister based random number
// generator that is used by the simulator and tests only.
package rngtest

import (
	"math/rand"
	"sync"

	"github.com/seehuhn/mt19937"
)

// SplitMix64 constants — identical to xoshiro256+ so that Reset produces
// the same deterministic (mainSeed, spinIdx) derivation in both generators.
const (
	smGamma     = uint64(0x9e3779b97f4a7c15)
	smMixConst1 = uint64(0xbf58476d1ce4e5b9)
	smMixConst2 = uint64(0x94d049bb133111eb)
	smShift1    = 30 // SplitMix64 first shift
	smShift2    = 27 // SplitMix64 second shift
	smShift3    = 31 // SplitMix64 final shift
)

// Mt19937 is a goroutine-safe random number generator based on the Mersenne Twister.
type Mt19937 struct {
	// rand is the underlying random source.
	rand *rand.Rand
	// mu guards rand and seed.
	mu *sync.Mutex
	// seed is the last seed that was set.
	seed int64
}

// Reseed sets a new seed and restarts the random sequence.
//
// Parameters:
//   - seed: new seed value
func (m *Mt19937) Reseed(seed int64) {
	m.rand.Seed(seed) //nolint:staticcheck // existing API, kept for backward compat
	m.seed = seed
}

// New создает рандом, который предназначен только для тестов.
// Нельзя использовать его в prod коде,
// так как он не является криптографически надежным.
func New() *Mt19937 {
	//nolint:gosec // криптография не нужна
	m := rand.New(mt19937.New())

	return &Mt19937{
		rand: m,
		mu:   &sync.Mutex{},
		seed: 0,
	}
}

// int returns a non-negative pseudo-random int under the mutex.
func (m *Mt19937) int() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.rand.Int()
}

// ShiftUint32 maps a raw random value into the inclusive range [minVal, maxVal].
//
// Parameters:
//   - got: raw random value
//   - minVal: lower bound of the range (inclusive)
//   - maxVal: upper bound of the range (inclusive)
//
// Returns:
//   - value in the range [minVal, maxVal]
func ShiftUint32(got, minVal, maxVal uint32) uint32 {
	if got == 0 && minVal == 0 {
		return 0
	}

	return got%(maxVal-minVal+1) + minVal
}

// Random returns a pseudo-random number in the inclusive range [minVal, maxVal].
// It returns 0 when both bounds are 0.
//
// Parameters:
//   - minVal: lower bound of the range (inclusive)
//   - maxVal: upper bound of the range (inclusive)
//
// Returns:
//   - random value in the range [minVal, maxVal]
func (m *Mt19937) Random(minVal, maxVal uint32) uint32 {
	if minVal == 0 && maxVal == 0 {
		return 0
	}

	return ShiftUint32(uint32(m.int()), minVal, maxVal) //nolint:gosec // G115: bit-pattern
}

// Seed returns the last seed that was set.
func (m *Mt19937) Seed() int64 {
	return m.seed
}

// Int64 returns a non-negative pseudo-random 63-bit integer.
func (m *Mt19937) Int64() int64 {
	return m.rand.Int63()
}

// Reset reseeds the generator deterministically from the (mainSeed, spinIdx) pair.
// The same arguments always produce the same random sequence.
// Uses the same SplitMix64 derivation as xoshiro256+.Reset so that both
// generators behave identically when called with identical arguments.
func (m *Mt19937) Reset(mainSeed int64, spinIdx int) {
	//nolint:gosec // spinIdx is a bit pattern for RNG init, not a magnitude
	combined := uint64(mainSeed) ^ (uint64(spinIdx) * smGamma)
	derived := splitMix64(&combined)

	m.mu.Lock()
	m.rand.Seed(int64(derived)) //nolint:staticcheck,gosec // existing API; bit-pattern cast
	m.seed = int64(derived)     //nolint:gosec // G115: bit-pattern cast for RNG seed
	m.mu.Unlock()
}

// splitMix64 advances the SplitMix64 state and returns the next output.
func splitMix64(state *uint64) uint64 {
	*state += smGamma
	val := *state
	val = (val ^ (val >> smShift1)) * smMixConst1
	val = (val ^ (val >> smShift2)) * smMixConst2

	return val ^ (val >> smShift3)
}
