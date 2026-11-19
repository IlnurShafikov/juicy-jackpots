//nolint:all // it's ok for simulator
package report

import (
	"fmt"
	"strings"

	"github.com/IlnurShafikov/slot-engine/pkg/tools/simtools"
)

const (
	// PlayerTableSeparator separates the columns of the player table.
	PlayerTableSeparator = " | "
	// PlayerTablePadding is the padding added around column values.
	PlayerTablePadding = " "
)

// PlayerReportBuilder builds player simulation reports
type PlayerReportBuilder struct {
	// builder accumulates the report text.
	builder strings.Builder
}

// NewPlayerReportBuilder creates a new PlayerReportBuilder
func NewPlayerReportBuilder() *PlayerReportBuilder {
	return &PlayerReportBuilder{
		builder: strings.Builder{},
	}
}

// WritePlayerReport generates a complete player simulation report as a plain-text string.
// The report includes the bet value, round cost, and a table with median spin counts
// grouped by initial player balance.
//
// Parameters:
//   - stats: per-balance-group statistics from the simulation
//   - roundCost: the raw round cost (will be multiplied by betValue for display)
//   - betValue: the bet amount per spin
//
// Returns:
//   - formatted report string ready to be written to a file or stdout
func (prb *PlayerReportBuilder) WritePlayerReport(
	stats []*PlayerStats,
	roundCost uint64,
	betValue int,
) string {
	prb.writeRoundCostValue(roundCost, betValue)
	prb.addEmptyLine()
	prb.writePlayerTable(stats, roundCost)

	return prb.builder.String()
}

// writeRoundCostValue writes bet value header
func (prb *PlayerReportBuilder) writeRoundCostValue(roundCost uint64, betValue int) {
	prb.builder.WriteString(fmt.Sprintf("BetValue: %d\n", betValue))
	prb.builder.WriteString(fmt.Sprintf("RoundCost: %d\n", roundCost*uint64(betValue)))
}

// writePlayerTable writes the main statistics table
func (prb *PlayerReportBuilder) writePlayerTable(stats []*PlayerStats, roundCost uint64) {
	prb.writeBalanceRow(stats, roundCost)
	prb.writePlayersCountRow(stats)
	prb.writeMedianRow(stats)
}

// writeBalanceRow writes table header with balances
func (prb *PlayerReportBuilder) writeBalanceRow(stats []*PlayerStats, roundCost uint64) {
	prb.builder.WriteString("Start balance:")
	prb.builder.WriteString(PlayerTablePadding)
	prb.builder.WriteString(PlayerTableSeparator)

	for i, stat := range stats {
		prb.builder.WriteString(prb.formatColumn(
			simtools.Uint64ToString(stat.InitialBalance/roundCost) + " bet"))

		if i < len(stats)-1 {
			prb.builder.WriteString(PlayerTableSeparator)
		}
	}

	prb.builder.WriteString("\n")
}

// writePlayersCountRow writes players count row
func (prb *PlayerReportBuilder) writePlayersCountRow(stats []*PlayerStats) {
	prb.builder.WriteString("Players count:")
	prb.builder.WriteString(PlayerTablePadding)
	prb.builder.WriteString(PlayerTableSeparator)

	for i, stat := range stats {
		prb.builder.WriteString(prb.formatColumn(
			simtools.IntToString(stat.PlayersCount)))

		if i < len(stats)-1 {
			prb.builder.WriteString(PlayerTableSeparator)
		}
	}

	prb.builder.WriteString("\n")
}

// writeMedianRow writes median spins row
func (prb *PlayerReportBuilder) writeMedianRow(stats []*PlayerStats) {
	prb.builder.WriteString("Median:")
	prb.builder.WriteString(strings.Repeat(PlayerTablePadding, 8))
	prb.builder.WriteString(PlayerTableSeparator)

	for i, stat := range stats {
		prb.builder.WriteString(prb.formatColumn(
			simtools.Float64ToString(stat.MedianSpins)))

		if i < len(stats)-1 {
			prb.builder.WriteString(PlayerTableSeparator)
		}
	}

	prb.builder.WriteString("\n")
}

// formatColumn formats a value for table column with padding
func (prb *PlayerReportBuilder) formatColumn(value string) string {
	const columnWidth = 7

	if len(value) >= columnWidth {
		return value
	}

	padding := columnWidth - len(value)

	return value + strings.Repeat(PlayerTablePadding, padding)
}

// addEmptyLine adds an empty line
func (prb *PlayerReportBuilder) addEmptyLine() {
	prb.builder.WriteString("\n")
}
