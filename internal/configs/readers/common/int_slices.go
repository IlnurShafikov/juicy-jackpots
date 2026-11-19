package common

import (
	"fmt"
	"strconv"

	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// ReadIntSlices reads a single-column CSV file and returns its values
// as a flat int slice. Used to load symbol counts, coin values, and reel indexes.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - filePath: absolute path within the embedded FS to the target CSV file
//
// Returns:
//   - int slice with one element per non-empty CSV cell
//   - error if the file cannot be read or any value cannot be parsed as an integer
func ReadIntSlices(
	csvReader *reader.Adapter,
	filePath string,
) ([]int, error) {
	data, err := csvReader.ReadCSV(filePath)
	if err != nil {
		return nil, fmt.Errorf("read symbols on data: %w", err)
	}

	resp, err := parseDataToSliceInt(data)
	if err != nil {
		return nil, fmt.Errorf("parse symbols on data: %w", err)
	}

	return resp, nil
}

// parseDataToSliceInt flattens a 2-D string slice into a single int slice by
// converting every cell value with strconv.Atoi.
//
// Parameters:
//   - data: raw CSV content as a 2-D string slice
//
// Returns:
//   - flattened int slice in row-major order
//   - error wrapping [ErrEmptyData] if data or any row is empty, or a parse error
func parseDataToSliceInt(data [][]string) ([]int, error) {
	if len(data) == 0 {
		return nil, ErrEmptyData
	}

	resp := make([]int, 0, len(data))

	for _, row := range data {
		if len(row) == 0 {
			return nil, ErrEmptyData
		}

		for _, str := range row {
			val, err := strconv.Atoi(str)
			if err != nil {
				return nil, fmt.Errorf("strconv.Atoi: %w", err)
			}

			resp = append(resp, val)
		}
	}

	return resp, nil
}
