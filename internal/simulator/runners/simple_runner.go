package runners

import (
	"fmt"

	resp "github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/parser"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools/rngtest"
)

// SimulationRunner handles the execution of simulation with progress tracking
type SimulationRunner struct {
	progressBar *result.Progress
	rngPool     *rngtest.Mt19937
}

// NewSimulationRunner creates a new SimulationRunner with progress bar
func NewSimulationRunner(spinsCount int) *SimulationRunner {
	return &SimulationRunner{
		progressBar: result.NewProgressBar(result.DefaultConfigs(spinsCount)),
		rngPool:     rngtest.New(),
	}
}

// RunSimulation executes the complete simulation with progress tracking
func (sr *SimulationRunner) RunSimulation(
	roundParser *parser.RoundParser,
	gameMod *game.Mod,
	config *simulator.SimulationConfig,
) error {
	if err := sr.executeSpinLoop(roundParser, gameMod, config); err != nil {
		return fmt.Errorf("run simulator: %w",
			simulator.NewSimulationError("execute spin loop", err))
	}

	sr.progressBar.Finish()

	return nil
}

// executeSpinLoop runs the main simulation loop.
// For the "mt" rng type, each spin is seeded by (mainSeed, spinIndex) via xoshiro.Reset,
// matching the parallel runner's behaviour exactly.
func (sr *SimulationRunner) executeSpinLoop(
	roundParser *parser.RoundParser,
	gameMod *game.Mod,
	config *simulator.SimulationConfig,
) error {
	rng := rngtest.New()
	rng.Reseed(config.MainSeed)

	for i := range config.SpinsCount {
		spinResult, err := sr.executeSpin(gameMod, rng, config, i)
		if err != nil {
			return fmt.Errorf("process spin %d: %w", i+1, err)
		}

		if err := roundParser.Parse(spinResult); err != nil {
			return fmt.Errorf("parse spin %d: %w", i+1, err)
		}

		sr.progressBar.Update(i)
	}

	return nil
}

// executeSpin performs a single game spin using the configured rng type.
// For "mt", it reseeds rng with (mainSeed, spinIdx) before each spin so that
// spin i always receives the same random numbers in both sequential and parallel runs.
func (sr *SimulationRunner) executeSpin(
	gameMod *game.Mod,
	rng *rngtest.Mt19937,
	config *simulator.SimulationConfig,
	spinIdx int,
) (*resp.Round, error) {
	switch config.RngType {
	case "mt":
		rng.Reset(config.MainSeed, spinIdx)

		spinResult, err := gameMod.Spin(rng)
		if err != nil {
			return nil, fmt.Errorf("game mod spin: %w", err)
		}

		return spinResult, nil

	default:
		return nil, simulator.ErrInvalidRngType
	}
}
