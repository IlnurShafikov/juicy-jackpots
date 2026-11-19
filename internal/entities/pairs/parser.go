package pairs

import (
	"fmt"
	"strconv"
)

// Parse converts CSV rows into pairs. The first column of a row is the key,
// the second one is the value; extra columns are ignored.
//
// Parameters:
//   - data: CSV rows, each with at least two integer columns
//
// Returns:
//   - parsed pairs in the order of the rows
//   - error if a row is too short or contains a non-integer value
func Parse(data [][]string) (Pairs, error) {
	const minLength = 2

	mappings := make(Pairs, len(data))

	for i, row := range data {
		if len(row) < minLength {
			return nil, ErrInvalidData
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID '%s': %w", row[0], err)
		}

		value, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, fmt.Errorf("invalid value '%s': %w", row[1], err)
		}

		mappings[i] = Pair{
			Key:   id,
			Value: value,
		}
	}

	return mappings, nil
}
