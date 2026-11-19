package common

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// ReadVolatility reads the game volatility setting from a CSV configuration file.
// The volatility indicates the risk level of the game (e.g., "low", "medium", "high").
//
// Parameters:
//   - csvReader: CSV adapter for reading configuration files
//   - generalPath: path builder for locating the volatility configuration file
//
// Returns:
//   - volatility string value
//   - error if file reading fails
func ReadVolatility(csvReader *reader.Adapter, generalPath *tools.PathBuilder) (string, error) {
	const volatilityFileName = "volatility"

	data, err := csvReader.ReadCSV(generalPath.File(volatilityFileName))
	if err != nil {
		return "", fmt.Errorf("read volatility file: %w", err)
	}

	return data[0][0], nil
}
