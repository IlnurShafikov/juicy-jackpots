package common

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/configs"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// ReadMultipliers reads and parses multiplier configurations from CSV files.
// It loads the maximum win multiplier and cost multiplier for the simple game mode.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - pointer to MultipliersConfigs containing all multiplier values
//   - error if reading or parsing fails
func ReadMultipliers(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (*configs.MultipliersConfigs, error) {
	const (
		maxWinName     = "max_win_limit"
		simpleCostName = "simple_cost_multiplier"
	)

	maxMultiplier, err := csvReader.ReadUint64(generalPath.File(maxWinName))
	if err != nil {
		return nil, fmt.Errorf("read max win multiplier: %w", err)
	}

	baseCostMultiplier, err := csvReader.ReadUint64(generalPath.File(simpleCostName))
	if err != nil {
		return nil, fmt.Errorf("read base cost multiplier: %w", err)
	}

	return &configs.MultipliersConfigs{
		MaxWinMultiplier:     maxMultiplier,
		SimpleCostMultiplier: baseCostMultiplier,
	}, nil
}
