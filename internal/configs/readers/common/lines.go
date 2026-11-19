package common

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/lines"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// ReadLines reads and parses win-line definitions from the general CSV directory.
// Win lines specify which cell positions on the grid form a paying combination
// and are shared across all RTP versions and game modes.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - Lines containing all defined win-line patterns
//   - error if reading or parsing the CSV file fails
func ReadLines(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (lines.Lines, error) {
	const fileName = "win_lines"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read win lines data: %w", err)
	}

	winLines, err := lines.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse win lines: %w", err)
	}

	return winLines, nil
}
