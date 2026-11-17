// Package report provides report generation and file-writing utilities for the simulator.
// It includes builders for human-readable simulation summaries and player-mode statistics,
// as well as writers that persist the output to the filesystem.
//
//nolint:all // it's ok for simulator
package report

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/result"
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// OutputFileExtension is the file extension used for simulation report files.
	OutputFileExtension = ".txt"
	// PathVersionPrefix is the directory prefix for RTP version segments in the output path.
	PathVersionPrefix = "v"
	// PathCostPrefix is the directory prefix for the round cost segment in the output path.
	PathCostPrefix = "cost"
	// PathSeedPrefix is the directory prefix for the seed segment in the output path.
	PathSeedPrefix = "seed"
)

// ResultsWriter creates and writes simulation result files to the output directory.
// Each simulation run produces a timestamped text file under a path that encodes
// the RTP version, round cost, game mod ID, and seed for easy identification.
type ResultsWriter struct{}

// NewResultsWriter creates a new ResultsWriter instance.
func NewResultsWriter() *ResultsWriter {
	return &ResultsWriter{}
}

// WriteResults writes the complete simulation result to a file.
// It creates any necessary output directories, opens a new file,
// and delegates report formatting to GeneralWriter.
//
// Parameters:
//   - results: the aggregated simulation statistics to write
//   - config: simulation configuration used to determine the output path
//   - startTime: the time the simulation started (included in the report header)
//
// Returns:
//   - error if directory creation, file creation, or writing fails
func (rw *ResultsWriter) WriteResults(
	results *result.Result,
	config *simulator.SimulationConfig,
	startTime time.Time,
) error {
	writer, closer, err := rw.createWriter(config)
	if err != nil {
		return simulator.NewSimulationError("create writer", err)
	}

	defer rw.safeCloseWriter(closer)

	return writer.WriteResult(results, startTime, config)
}

// createWriter creates a GeneralWriter with file output
func (rw *ResultsWriter) createWriter(config *simulator.SimulationConfig) (*GeneralWriter, io.Closer, error) {
	outputPath := rw.buildOutputPath(config)

	if err := rw.ensureDirectoryExists(outputPath); err != nil {
		return nil, nil, fmt.Errorf("ensure directory exists: %w", err)
	}

	filePath := rw.buildFilePath(outputPath)

	file, err := rw.createOutputFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	return NewGeneralWriter(file), file, nil
}

// safeCloseWriter safely closes the writer and logs any errors
func (rw *ResultsWriter) safeCloseWriter(closer io.Closer) {
	if closer == nil {
		return
	}

	if err := closer.Close(); err != nil {
		log.Printf("failed to close writer: %v", err)
	}
}

// buildOutputPath constructs the output directory path
func (rw *ResultsWriter) buildOutputPath(config *simulator.SimulationConfig) string {
	return buildOutputPath(
		config.RtpVersion,
		config.RoundCost,
		config.GameModID,
		config.MainSeed,
	)
}

// buildFilePath constructs the full file path with timestamp
func (rw *ResultsWriter) buildFilePath(outputPath string) string {
	fileName := rw.generateFileName()

	return filepath.Join(outputPath, fileName)
}

// generateFileName creates a timestamped filename
func (rw *ResultsWriter) generateFileName() string {
	return time.Now().Format("2006-01-02") + OutputFileExtension
}

// ensureDirectoryExists creates the directory if it doesn't exist
func (rw *ResultsWriter) ensureDirectoryExists(path string) error {
	err := os.MkdirAll(path, simulator.FilePermissions)
	if err != nil {
		return fmt.Errorf("create directory %q: %w", path, err)
	}

	return nil
}

// createOutputFile creates and opens the output file
func (rw *ResultsWriter) createOutputFile(filePath string) (*os.File, error) {
	//nolint:gosec // Potential file inclusion via variable
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file %q: %w", filePath, err)
	}

	return file, nil
}

// buildOutputPath constructs directory path from simulation parameters
func buildOutputPath(rtpVersion int, roundCost uint64, featureID string, mainSeed int64) string {
	return filepath.Join(
		simulator.OutputDir,
		fmt.Sprintf("%s%d", PathVersionPrefix, rtpVersion),
		fmt.Sprintf("%s%d", PathCostPrefix, roundCost),
		featureID,
		fmt.Sprintf("%s%d", PathSeedPrefix, mainSeed),
	)
}
