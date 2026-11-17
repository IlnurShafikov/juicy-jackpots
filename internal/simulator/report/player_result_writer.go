//nolint:all // it's ok for simulator
package report

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
)

const (
	PlayerOutputFileExtension = ".txt"
	PlayerPathPrefix          = "player_simulation"
)

// PlayerResultsWriter handles writing player simulation results to files
type PlayerResultsWriter struct{}

// NewPlayerResultsWriter creates a new PlayerResultsWriter instance
func NewPlayerResultsWriter() *PlayerResultsWriter {
	return &PlayerResultsWriter{}
}

// WritePlayerResults writes player simulation statistics to a timestamped file
// under the output directory. The file path encodes RTP version, round cost,
// game mod ID, player count, and seed for easy identification.
//
// Parameters:
//   - stats: per-balance-group simulation statistics produced by PlayerSimulationRunner
//   - config: simulation configuration used to build the output path and report header
//
// Returns:
//   - error if directory creation, file creation, or writing fails
func (prw *PlayerResultsWriter) WritePlayerResults(
	stats []*PlayerStats,
	config *simulator.SimulationConfig,
) error {
	writer, closer, err := prw.createWriter(config)
	if err != nil {
		return simulator.NewSimulationError("create player writer", err)
	}

	defer prw.safeCloseWriter(closer)

	return prw.writeReport(writer, stats, config.RoundCost, config.BetValue)
}

// createWriter creates a writer with file output
func (prw *PlayerResultsWriter) createWriter(
	config *simulator.SimulationConfig,
) (io.Writer, io.Closer, error) {
	outputPath := prw.buildOutputPath(config)

	if err := prw.ensureDirectoryExists(outputPath); err != nil {
		return nil, nil, fmt.Errorf("ensure directory exists: %w", err)
	}

	filePath := prw.buildFilePath(outputPath)

	file, err := prw.createOutputFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	return file, file, nil
}

// safeCloseWriter safely closes the writer and logs any errors
func (prw *PlayerResultsWriter) safeCloseWriter(closer io.Closer) {
	if closer == nil {
		return
	}

	if err := closer.Close(); err != nil {
		log.Printf("failed to close player writer: %v", err)
	}
}

// buildOutputPath constructs the output directory path for player simulation
func (prw *PlayerResultsWriter) buildOutputPath(config *simulator.SimulationConfig) string {
	return filepath.Join(
		simulator.OutputDir,
		PlayerPathPrefix,
		fmt.Sprintf("v%d", config.RtpVersion),
		fmt.Sprintf("cost%d", config.RoundCost),
		config.GameModID,
		fmt.Sprintf("players%d", config.PlayersCount),
		fmt.Sprintf("seed%d", config.MainSeed),
	)
}

// formatBalances formats balance slice for path
func (prw *PlayerResultsWriter) formatBalances(balances []uint64) string {
	strBalances := make([]string, len(balances))
	for i, b := range balances {
		strBalances[i] = fmt.Sprintf("%d", b)
	}

	return strings.Join(strBalances, "_")
}

// buildFilePath constructs the full file path with timestamp
func (prw *PlayerResultsWriter) buildFilePath(outputPath string) string {
	fileName := prw.generateFileName()

	return filepath.Join(outputPath, fileName)
}

// generateFileName creates a timestamped filename
func (prw *PlayerResultsWriter) generateFileName() string {
	return time.Now().Format("2006-01-02") + PlayerOutputFileExtension
}

// ensureDirectoryExists creates the directory if it doesn't exist
func (prw *PlayerResultsWriter) ensureDirectoryExists(path string) error {
	err := os.MkdirAll(path, simulator.FilePermissions)
	if err != nil {
		return fmt.Errorf("create directory %q: %w", path, err)
	}

	return nil
}

// createOutputFile creates and opens the output file
func (prw *PlayerResultsWriter) createOutputFile(filePath string) (*os.File, error) {
	//nolint:gosec // Potential file inclusion via variable
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file %q: %w", filePath, err)
	}

	return file, nil
}

// writeReport writes the player simulation report
func (prw *PlayerResultsWriter) writeReport(
	writer io.Writer,
	stats []*PlayerStats,
	roundCost uint64,
	betValue int,
) error {
	reportBuilder := NewPlayerReportBuilder()
	report := reportBuilder.WritePlayerReport(stats, roundCost, betValue)

	if _, err := writer.Write([]byte(report)); err != nil {
		return fmt.Errorf("write player report: %w", err)
	}

	return nil
}
