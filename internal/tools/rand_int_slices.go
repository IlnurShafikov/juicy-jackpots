package tools

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/tools/cast"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// RandomIntSlices returns a randomly shuffled slice of needCount unique integers
// drawn from the range [0, maxValue). The Fisher-Yates algorithm is used to
// shuffle all candidate indexes and the first needCount elements are returned.
//
// Parameters:
//   - rng: random number generator used for index swapping
//   - maxValue: exclusive upper bound of the candidate index range (must be > 0)
//   - needCount: number of unique indexes to return (must satisfy 0 < needCount <= maxValue)
//
// Returns:
//   - slice of needCount unique random integers in [0, maxValue)
//   - [ErrInvalidArguments] if needCount <= 0, maxValue == 0, or needCount > maxValue
//   - error if the internal safe-cast from int to uint32 fails
func RandomIntSlices(
	rng model.RNG,
	maxValue int,
	needCount int,
) ([]int, error) {
	if needCount <= 0 || maxValue == 0 || needCount > maxValue {
		return nil, ErrInvalidArguments
	}

	resp := make([]int, 0, maxValue)

	for row := range maxValue {
		resp = append(resp, row)
	}

	for i := len(resp) - 1; i > 0; i-- {
		safeValue, err := cast.SafeCast[uint32](i)
		if err != nil {
			return nil, fmt.Errorf("safe cast int to uint32: %w", err)
		}

		j := rng.Random(0, safeValue)
		resp[i], resp[j] = resp[j], resp[i]
	}

	return resp[:needCount], nil
}
