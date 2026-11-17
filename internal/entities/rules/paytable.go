package rules

import "github.com/IlnurShafikov/slot-engine/pkg/paytable"

// PayTable represents the game's payout structure.
// It defines how much each symbol pays for different match counts.
type PayTable struct {
	// Symbols maps symbol IDs to their payout values for different match counts
	Symbols paytable.SymbolsPayouts `json:"symbols"`

	// Ranges defines the valid count ranges for symbol matches (e.g., 3-5 symbols)
	Ranges paytable.CountRanges `json:"ranges"`
}

// NewPaytable creates a new PayTable by copying data from a source paytable.
// This ensures the game has its own copy of the paytable data.
//
// Parameters:
//   - payTable: source paytable to copy from
//
// Returns:
//   - pointer to a new PayTable with copied data
func NewPaytable(payTable *paytable.PayTable) *PayTable {
	return &PayTable{
		Symbols: payTable.CopySymbolsPayouts(),
		Ranges:  payTable.CopyCountRanges(),
	}
}
