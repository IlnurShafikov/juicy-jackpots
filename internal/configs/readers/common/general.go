package common

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/paytable"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// ReadPaytable reads and parses the paytable configuration from a CSV file.
// The paytable defines symbol payouts and count ranges for winning combinations.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - pointer to PayTable containing symbol payouts and ranges
//   - error if reading or parsing fails
func ReadPaytable(
	reader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (*paytable.PayTable, error) {
	const fileName = "paytable"

	data, err := reader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read paytable: %w", err)
	}

	resp, err := paytable.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse paytable: %w", err)
	}

	return resp, nil
}
