package pairs

import (
	"fmt"
	"strconv"
)

// Parse converts raw CSV rows into a Pairs slice.
// Each row must have at least two columns: key (integer) and value (integer).
//
// Parameters:
//   - data: raw CSV content as a 2-D string slice, one row per pair
//
// Returns:
//   - Pairs with one entry per row
//   - [ErrInvalidData] if a row has fewer than two columns, or a parse error
//     if any column is not a valid integer
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
