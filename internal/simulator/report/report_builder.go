//nolint:all // it's ok for simulator
package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/result"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/simtools"
)

// Constants for report formatting
const (
	SectionSeparator = "\n"
	TabIndent        = "\t"
)

// Section headers
const (
	GameModId  = "Game Mod = "
	RtpVersion = "RTP Version = "

	HeaderMainSeed   = "Main seed "
	HeaderStartedAt  = "Started at "
	HeaderFinishedAt = "Finished at "
	HeaderRounds     = "Rounds "

	SectionTotalRTP = "TOTAL RTP = "
	SectionBaseRTP  = "BASE RTP = "
	SectionFreeRTP  = "FS RTP = "
	SectionSD       = "SD = "
	SectionVI       = "VI = "

	LabelByWins         = "By anywhere wins = "
	LabelFreeTriggered  = "Free games triggered = "
	LabelAllFreeSteps   = "All free game steps = "
	LabelRetriggers     = "Retriggers = "
	LabelRetriggerSteps = "Retriggers steps count = "
	LabelMaxWinCount    = "MaxWin reached count = "
	LabelMaxTotalWin    = "Max total win = "
	LabelStepsWithWins  = "Steps with wins = "
	LabelHitRate        = "Hit rate = "

	TitleBaseWinRanges = "Only Base Phase TotalWin ranges"
	TitleFreeWinRanges = "Only Free Phase TotalWin ranges"
)

// Config carries all metadata required to generate a simulation summary report.
type Config struct {
	// StartTime is when the simulation run began (used for elapsed-time display).
	StartTime time.Time
	// GameModId identifies the game mod that was simulated.
	GameModId string
	// MainSeed is the RNG seed used for the simulation.
	MainSeed int64
	// RoundCost is the cost of a single round in base units.
	RoundCost uint64
	// TotalBet is the sum of all bets placed during the simulation (used for RTP calculation).
	TotalBet uint64
	// Spins is the total number of simulated rounds.
	Spins int
	// RtpVersion is the RTP percentage version that was simulated (e.g. 96).
	RtpVersion int
}

// GeneralWriter wraps an io.Writer and provides helper methods
// for writing formatted simulation output.
type GeneralWriter struct {
	writer io.Writer
}

// NewGeneralWriter creates a new GeneralWriter that writes to the given io.Writer.
//
// Parameters:
//   - writer: the underlying writer (e.g. os.File or bytes.Buffer)
func NewGeneralWriter(writer io.Writer) *GeneralWriter {
	return &GeneralWriter{writer: writer}
}

// Builder accumulates formatted simulation report text using a strings.Builder.
// Call WriteReport to populate it, then GetFormattedReport to retrieve the result.
type Builder struct {
	builder strings.Builder
}

// NewReportBuilder creates a new, empty Builder.
func NewReportBuilder() *Builder {
	return &Builder{
		builder: strings.Builder{},
	}
}

// addLine adds a line with line break
func (rb *Builder) addLine(str string) {
	rb.builder.WriteString(str + "\n")
}

// addTabLine adds a line with tab indentation and line break
func (rb *Builder) addTabLine(str string) {
	rb.builder.WriteString(TabIndent + str + "\n")
}

// addSection adds section separator
func (rb *Builder) addSection() {
	rb.builder.WriteString(SectionSeparator)
}

// GetFormattedReport returns the accumulated report text as a string.
func (rb *Builder) GetFormattedReport() string {
	return rb.builder.String()
}

// WriteReport generates a complete simulation summary report and appends it to the builder.
// The report includes a header, RTP breakdown by phase, volatility metrics,
// maximum win information, hit rate, and win-range distributions.
//
// Parameters:
//   - res: the aggregated simulation result to report on
//   - config: metadata for the report header (game mod, seed, timing, etc.)
func (rb *Builder) WriteReport(res *result.Result, config Config) {
	rb.writeReportHeader(config)
	rb.writeTotalRTPSection(config, res)
	rb.writeBaseGameSection(config, res.Base())
	rb.writeVolatilitySection(res, config)
	rb.writeMaxWinSection(res, config.RoundCost)
	rb.writeHitRateSection(res, config.Spins)
	rb.writeWinRangesSection(res, config)
}

// writeReportHeader writes report header
func (rb *Builder) writeReportHeader(config Config) {
	rb.addLine(GameModId + config.GameModId)
	rb.addLine(RtpVersion + simtools.IntToString(config.RtpVersion))
	rb.addSection()
	rb.addLine(HeaderMainSeed + simtools.Int64ToString(config.MainSeed))
	rb.addLine(HeaderStartedAt + config.StartTime.Format(time.RFC3339))
	rb.addLine(HeaderFinishedAt + time.Now().Format(time.RFC3339))
	rb.addLine(HeaderRounds + simtools.IntToString(config.Spins))
}

// writeTotalRTPSection writes total RTP section
func (rb *Builder) writeTotalRTPSection(config Config, res *result.Result) {
	rb.addSection()

	totalWins := res.Base().GetWins()
	totalRTP := simtools.CalculateRTP(config.TotalBet, totalWins)
	rb.addLine(SectionTotalRTP + simtools.Float64ToString(totalRTP))
}

// writeBaseGameSection writes base game section
func (rb *Builder) writeBaseGameSection(config Config, baseRes *result.BaseGame) {
	rb.addSection()

	baseRTP := simtools.CalculateRTP(config.TotalBet, baseRes.GetWins())

	rb.addLine(SectionBaseRTP + simtools.Float64ToString(baseRTP))
	rb.addTabLine(LabelByWins + simtools.Float64ToString(baseRTP))
}

// writeVolatilitySection writes volatility section
func (rb *Builder) writeVolatilitySection(res *result.Result, config Config) {
	rb.addSection()

	totalWins := res.Base().GetWins()
	totalRTP := simtools.CalculateRTP(config.TotalBet, totalWins)

	sd := res.Xi().Sd(totalRTP, float64(config.Spins))
	volatility := simtools.CalculateVolatility(sd)

	rb.addLine(SectionSD + simtools.Float64ToString(sd))
	rb.addLine(SectionVI + simtools.Float64ToString(volatility))
}

// writeMaxWinSection writes max win section
func (rb *Builder) writeMaxWinSection(res *result.Result, roundCost uint64) {
	rb.addSection()

	maxWinMultiplier := res.MaxTotalWin() / roundCost

	rb.addLine(LabelMaxWinCount + simtools.Uint64ToString(res.MaxWinReached()))
	rb.addLine(LabelMaxTotalWin + simtools.Uint64ToString(maxWinMultiplier) + "x")
	rb.addLine(LabelMaxTotalWin + simtools.Uint64ToString(res.MaxTotalWin()))
}

// writeHitRateSection writes hit rate section
func (rb *Builder) writeHitRateSection(res *result.Result, spins int) {
	rb.addSection()

	stepsWithWins := res.PositiveBets().GetSpins()
	hitRate := res.PositiveBets().CalculateHitRate(spins)

	rb.addLine(LabelStepsWithWins + simtools.IntToString(stepsWithWins))
	rb.addLine(LabelHitRate + simtools.Float64ToString(hitRate))
}

// writeWinRangesSection writes win ranges section
func (rb *Builder) writeWinRangesSection(res *result.Result, repCfg Config) {
	rb.addSection()
	rb.writeWinRangeDetails(TitleBaseWinRanges, res.BaseWinRanges, repCfg.Spins, repCfg.TotalBet)
	rb.addSection()
}

func (rb *Builder) writeWinRangeDetails(
	title string, winRanges *result.WinRanges, stepsCount int, totalBet uint64,
) {
	rb.addSection()
	rb.addLine(title)

	rb.addLine(fmt.Sprintf("%-18s | %-13s | %-10s | %-8s |",
		"Win Range", "Phase Count", "Once In", "RTP"))

	allRanges := winRanges.GetSortedRanges()

	for _, curRange := range allRanges {
		currOnce := float64(1) / (float64(curRange.Count) / float64(stepsCount))
		rangeRtp := simtools.CalculateRTP(totalBet, curRange.TotalSum)

		line := fmt.Sprintf("%-18s | %-13d | %-10.2f | %-8.4f |",
			curRange.Range.Name, curRange.Count, currOnce, rangeRtp)

		rb.addLine(line)
	}
}

// WriteResult formats and writes the full simulation report to the underlying io.Writer.
// It builds the Config from the simulation config, generates the report text,
// and flushes it to the writer in a single call.
//
// Parameters:
//   - res: the aggregated simulation result
//   - startedTime: the time the simulation started (included in the report header)
//   - simCfg: simulation configuration used for report metadata
//
// Returns:
//   - error if writing to the underlying writer fails
func (g *GeneralWriter) WriteResult(
	res *result.Result, startedTime time.Time, simCfg *simulator.SimulationConfig,
) error {
	config := g.buildReportConfig(simCfg, startedTime)
	report := g.generateReport(res, config)

	return g.writeReport(report)
}

// buildReportConfig constructs a Config from a SimulationConfig and start time.
func (g *GeneralWriter) buildReportConfig(
	simCfg *simulator.SimulationConfig, startedTime time.Time,
) Config {
	//nolint:gosec // integer overflow is not a concern here
	totalBet := simCfg.RoundCost * uint64(simCfg.SpinsCount)

	return Config{
		MainSeed:   simCfg.MainSeed,
		StartTime:  startedTime,
		RoundCost:  simCfg.RoundCost,
		Spins:      simCfg.SpinsCount,
		TotalBet:   totalBet,
		GameModId:  simCfg.GameModID,
		RtpVersion: simCfg.RtpVersion,
	}
}

// generateReport creates and returns the formatted report string using Builder.
func (g *GeneralWriter) generateReport(res *result.Result, config Config) string {
	builder := NewReportBuilder()
	builder.WriteReport(res, config)

	return builder.GetFormattedReport()
}

// writeReport flushes the report string to the underlying io.Writer.
func (g *GeneralWriter) writeReport(report string) error {
	if _, err := g.writer.Write([]byte(report)); err != nil {
		return fmt.Errorf("writer.Write(): %w", err)
	}

	return nil
}
