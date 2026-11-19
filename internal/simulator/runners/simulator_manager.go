// Package runners provides the execution layer of the simulator.
// It contains three simulation modes — simple (batch), interactive (step-by-step),
// and player (balance-depletion) — coordinated by SimulationManager.
//
// Typical usage:
//
//	cfg := simulator.ParseConfigs()
//	if err := runners.NewSimulationManager().RunSimulation(cfg); err != nil {
//	    log.Fatal(err)
//	}
//
//nolint:all // it's ok for simulator
package runners

import (
	"fmt"
	"time"

	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/parser"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/report"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/result"
)

// SimulationManager orchestrates the entire simulation process
type SimulationManager struct {
	// gameModFinder looks up the game mod of the requested RTP version.
	gameModFinder *simulator.GameModFinder
	// resultsWriter writes the simulation report.
	resultsWriter *report.ResultsWriter
	// interactiveRunner runs the step-by-step interactive mode.
	interactiveRunner *InteractiveRunner
	// playerResultsWriter writes the player simulation report.
	playerResultsWriter *report.PlayerResultsWriter
}

// NewSimulationManager creates a new SimulationManager with all dependencies
func NewSimulationManager() *SimulationManager {
	return &SimulationManager{
		gameModFinder:       simulator.NewGameModFinder(),
		resultsWriter:       report.NewResultsWriter(),
		interactiveRunner:   NewInteractiveRunner(),
		playerResultsWriter: report.NewPlayerResultsWriter(),
	}
}

// RunSimulation runs the complete simulation process based on configuration
func (sm *SimulationManager) RunSimulation(cfg *simulator.SimulationConfig) error {
	gameMod, err := sm.findAndConfigureGameMod(cfg)
	if err != nil {
		return fmt.Errorf("run simulation: %w",
			simulator.NewSimulationError("find and configure game mod", err))
	}

	cfg.RoundCost = gameMod.CostMultiplier()

	if cfg.PlayMod {
		return sm.runInteractiveMode(gameMod, cfg)
	}

	if cfg.PlayerMod {
		return sm.runPlayerMode(gameMod, cfg)
	}

	return sm.runSimpleMode(gameMod, cfg)
}

// findAndConfigureGameMod finds the appropriate game mod for the simulation
func (sm *SimulationManager) findAndConfigureGameMod(
	cfg *simulator.SimulationConfig,
) (*game.Mod, error) {
	gameMod, err := sm.gameModFinder.FindGameMod(cfg.RtpVersion, cfg.GameModID, cfg.FeatureID)
	if err != nil {
		return nil, fmt.Errorf("find game mod (RTP: %d, ModID: %s, FeatureID: %s): %w",
			cfg.RtpVersion, cfg.GameModID, cfg.FeatureID, err)
	}

	return gameMod, nil
}

// runSimpleMode executes simulation in standard mode and writes results to file
func (sm *SimulationManager) runSimpleMode(
	gameMod *game.Mod,
	cfg *simulator.SimulationConfig,
) error {
	startTime := time.Now()

	results, err := sm.executeSimulation(gameMod, cfg)
	if err != nil {
		return fmt.Errorf("run simple mode: %w",
			simulator.NewSimulationError("execute batch simulation", err))
	}

	if err := sm.writeSimulationResults(results, cfg, startTime); err != nil {
		return fmt.Errorf("run simple mode: %w",
			simulator.NewSimulationError("write simulation results", err))
	}

	return nil
}

// runInteractiveMode executes simulation in interactive mode
func (sm *SimulationManager) runInteractiveMode(
	gameMod *game.Mod,
	cfg *simulator.SimulationConfig,
) error {
	return sm.interactiveRunner.runInteractiveSimulation(gameMod, cfg)
}

// runPlayerMode executes simulation in player mode
func (sm *SimulationManager) runPlayerMode(
	gameMod *game.Mod,
	cfg *simulator.SimulationConfig,
) error {
	if err := sm.validatePlayerModeConfig(cfg); err != nil {
		return fmt.Errorf("run player mode: %w",
			simulator.NewSimulationError("validate config", err))
	}

	stats, err := sm.executePlayerSimulation(gameMod, cfg)
	if err != nil {
		return fmt.Errorf("run player mode: %w",
			simulator.NewSimulationError("execute player simulation", err))
	}

	if err := sm.writePlayerResults(stats, cfg); err != nil {
		return fmt.Errorf("run player mode: %w",
			simulator.NewSimulationError("write player results", err))
	}

	return nil
}

// executeSimulation runs the core simulation logic and returns results
func (sm *SimulationManager) executeSimulation(
	gameMod *game.Mod,
	cfg *simulator.SimulationConfig,
) (*result.Result, error) {
	roundParser := sm.createSimulationComponents(gameMod, cfg)

	if err := sm.runSimulationLoop(roundParser, gameMod, cfg); err != nil {
		return nil, fmt.Errorf("run simulation loop: %w", err)
	}

	return roundParser.GetResult(), nil
}

// executePlayerSimulation runs player simulation and returns stats
func (sm *SimulationManager) executePlayerSimulation(
	gameMod *game.Mod,
	cfg *simulator.SimulationConfig,
) ([]*report.PlayerStats, error) {
	playerRunner := NewPlayerSimulationRunner(cfg.SimpleCostMultiplier)

	stats, err := playerRunner.RunPlayerSimulation(gameMod, cfg)
	if err != nil {
		return nil, fmt.Errorf("run player simulation: %w", err)
	}

	return stats, nil
}

// createSimulationComponents creates the round parser for a simulation run.
func (sm *SimulationManager) createSimulationComponents(
	gameMod *game.Mod,
	cfg *simulator.SimulationConfig,
) *parser.RoundParser {
	return parser.NewRoundParser(gameMod.CostMultiplier(), cfg.SimpleCostMultiplier)
}

// runSimulationLoop executes the simulation loop, choosing sequential or parallel runner.
func (sm *SimulationManager) runSimulationLoop(
	roundParser *parser.RoundParser,
	gameMod *game.Mod,
	cfg *simulator.SimulationConfig,
) error {
	if cfg.WorkersCount > 1 {
		runner := NewParallelRunner(cfg.SpinsCount, cfg.WorkersCount)
		return runner.RunSimulation(roundParser, gameMod, cfg)
	}

	runner := NewSimulationRunner(cfg.SpinsCount)
	return runner.RunSimulation(roundParser, gameMod, cfg)
}

// writeSimulationResults writes the simulation results to output
func (sm *SimulationManager) writeSimulationResults(
	results *result.Result,
	cfg *simulator.SimulationConfig,
	startTime time.Time,
) error {
	err := sm.resultsWriter.WriteResults(results, cfg, startTime)
	if err != nil {
		return fmt.Errorf("write results: %w", err)
	}

	return nil
}

// writePlayerResults writes player simulation results to output
func (sm *SimulationManager) writePlayerResults(
	stats []*report.PlayerStats,
	cfg *simulator.SimulationConfig,
) error {
	err := sm.playerResultsWriter.WritePlayerResults(stats, cfg)
	if err != nil {
		return fmt.Errorf("write player results: %w", err)
	}

	return nil
}

// validatePlayerModeConfig validates player mode configuration
func (sm *SimulationManager) validatePlayerModeConfig(cfg *simulator.SimulationConfig) error {
	if len(cfg.PlayerBalances) == 0 {
		return fmt.Errorf("no player balances provided")
	}

	if cfg.PlayersCount <= 0 {
		return fmt.Errorf("players count must be positive")
	}

	return nil
}
