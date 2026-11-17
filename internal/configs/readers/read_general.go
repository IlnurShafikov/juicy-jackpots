package readers

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/configs"
	"github.com/IlnurShafikov/juicy-jackpots/internal/configs/readers/common"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

const (
	// generalPrefix defines the directory name for general configuration files
	generalPrefix = "general"
	// symbolOnReelName is the CSV file name for the per-reel visible symbol counts.
	symbolOnReelName = "symbols_on_reel"
)

// Read loads and parses the general game configuration.
// It reads the paytable and multipliers from CSV files in the general directory.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//
// Returns:
//   - pointer to GeneralConfig containing paytable and multipliers
//   - error if reading or parsing fails
func Read(reader *reader.Adapter) (*configs.GeneralConfig, error) {
	generalPath := tools.NewPathBuilder(generalPrefix)

	paytable, err := common.ReadPaytable(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read paytable: %w", err)
	}

	symbolsOnReel, err := common.ReadIntSlices(reader, generalPath.File(symbolOnReelName))
	if err != nil {
		return nil, fmt.Errorf("read symbols on reel: %w", err)
	}

	winLines, err := common.ReadLines(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read win lines: %w", err)
	}

	multipliers, err := common.ReadMultipliers(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read multipliers: %w", err)
	}

	volatility, err := common.ReadVolatility(reader, generalPath)
	if err != nil {
		return nil, fmt.Errorf("read volatility: %w", err)
	}

	return &configs.GeneralConfig{
		PayTable:       paytable,
		Multipliers:    multipliers,
		SymbolsOnReels: symbolsOnReel,
		WinLines:       winLines,
		Volatility:     volatility,
	}, nil
}
