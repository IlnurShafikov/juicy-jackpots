package common

import (
	"fmt"
	"strconv"

	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/pairs"
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
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (*paytable.PayTable, error) {
	const fileName = "paytable"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read paytable: %w", err)
	}

	resp, err := paytable.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse paytable: %w", err)
	}

	return resp, nil
}

// ReadJackpotPairs reads the jackpot-values CSV file and parses it into a
// pairs.Pairs  that associates each coin count with its payout multiplier.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - Pairs mapping jackpot type ID to multiplier value
//   - error if the file cannot be read or parsed
func ReadJackpotPairs(
	csvReader *reader.Adapter,
	generalPath *tools.PathBuilder,
) (pairs.Pairs, error) {
	const fileName = "jp_pairs"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read jackpot data: %w", err)
	}

	jackpot, err := pairs.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse jackpot data: %w", err)
	}

	return jackpot, nil
}

// ReadJpTypes reads and parses the jackpot types configuration from a CSV file.
// This configuration maps jackpot values to their corresponding jackpot type identifiers.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - generalPath: path builder pointing to the general configuration directory
//
// Returns:
//   - map where keys are jackpot values (uint64) and values are jackpot type IDs (int)
//   - error if reading or parsing fails
func ReadJpTypes(csvReader *reader.Adapter, generalPath *tools.PathBuilder) (map[uint64]int, error) {
	const fileName = "jp_types"

	data, err := csvReader.ReadCSV(generalPath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read jp types data: %w", err)
	}

	jpTypes, err := parseJpTypes(data)
	if err != nil {
		return nil, fmt.Errorf("parse jp types: %w", err)
	}

	return jpTypes, nil
}

// parseJpTypes parses jackpot types configuration from CSV data into a map.
// Each row must contain exactly 2 values: jackpot type ID and jackpot value.
//
// Parameters:
//   - data: 2D array of strings where each row contains [jpType, jpValue]
//
// Returns:
//   - map where keys are jackpot values (uint64) and values are jackpot type IDs (int)
//   - error if data format is invalid or conversion fails
func parseJpTypes(data [][]string) (map[uint64]int, error) {
	const dataLen = 2

	jpTypes := make(map[uint64]int, len(data))

	for _, row := range data {
		if len(row) != dataLen {
			return nil, ErrInvalidData
		}

		jpType, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("parse jp id '%s': %w", row[0], err)
		}

		jpValue, err := strconv.ParseUint(row[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse jp type '%s': %w", row[1], err)
		}

		jpTypes[jpValue] = jpType
	}

	return jpTypes, nil
}
