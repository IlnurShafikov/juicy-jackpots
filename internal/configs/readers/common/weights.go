package common

import (
	"fmt"

	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

// ReadWeight reads and parses a generic weight distribution configuration from a CSV file.
// This is a utility function used by other weight reading functions to parse weight data.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - filePath: complete path to the CSV file containing weight distribution
//
// Returns:
//   - pointer to Weights containing the parsed probability distribution
//   - error if reading or parsing fails
func ReadWeight(
	csvReader *reader.Adapter,
	filePath string,
) (*weight.Weights, error) {
	data, err := csvReader.ReadCSV(filePath)
	if err != nil {
		return nil, fmt.Errorf("read weights data: %w", err)
	}

	resp, err := weight.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse weights: %w", err)
	}

	return resp, nil
}
