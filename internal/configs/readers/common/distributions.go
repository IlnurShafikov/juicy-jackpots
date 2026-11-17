package common

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

// ReadDistributions reads and parses symbol distribution configurations from a CSV file.
// Distributions define the weighted probabilities for symbols appearing on reels.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - phasePath: path builder pointing to the game phase configuration directory
//
// Returns:
//   - pointer to Distributions containing weighted symbol distributions
//   - error if reading or parsing fails
func ReadDistributions(
	reader *reader.Adapter,
	phasePath *tools.PathBuilder,
) (*weight.Weights, error) {
	const fileName = "distributions"

	data, err := reader.ReadCSV(phasePath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read distributions data: %w", err)
	}

	resp, err := weight.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse weight: %w", err)
	}

	return resp, nil
}
