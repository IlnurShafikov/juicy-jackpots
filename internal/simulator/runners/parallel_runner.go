package runners

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/parser"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools/rngtest"
)

const (
	// progressTickMs is how often the progress reporter samples the atomic spin counter.
	progressTickMs = 100

	// progressUpdateMask is the bitmask used by bar.DefaultConfigs to gate renders.
	// Update(n) only renders when n&progressUpdateMask == 0, so we must align the value.
	progressUpdateMask = 1023
)

// parallelRunner executes the simulation loop across multiple goroutines.
// Each goroutine owns an independent xoshiro256+ generator positioned 2^128
// steps apart from its neighbours, guaranteeing non-overlapping random streams.
// Reproducibility holds for identical (mainSeed, workersCount) pairs.
type parallelRunner struct {
	progressBar *result.Progress
	numWorkers  int
}

// NewParallelRunner creates a parallelRunner for the given spin count and worker count.
func NewParallelRunner(spinsCount, numWorkers int) *parallelRunner {
	return &parallelRunner{
		progressBar: result.NewProgressBar(result.DefaultConfigs(spinsCount)),
		numWorkers:  numWorkers,
	}
}

// RunSimulation executes spins in parallel and merges worker results into roundParser.
func (pr *parallelRunner) RunSimulation(
	roundParser *parser.RoundParser,
	gameMod *game.Mod,
	config *simulator.SimulationConfig,
) error {
	results, err := pr.runWorkers(gameMod, config)
	if err != nil {
		return fmt.Errorf("parallel simulation: %w", err)
	}

	for _, res := range results {
		roundParser.MergeResult(res)
	}

	pr.progressBar.Finish()

	return nil
}

// runWorkers launches all worker goroutines and collects their results.
func (pr *parallelRunner) runWorkers(
	gameMod *game.Mod,
	config *simulator.SimulationConfig,
) ([]*result.Result, error) {
	var (
		completed atomic.Int64
		wg        sync.WaitGroup
	)

	results := make([]*result.Result, pr.numWorkers)
	errs := make([]error, pr.numWorkers)

	stopProgress := pr.startProgressReporter(&completed)

	for i := range pr.numWorkers {
		wg.Add(1)

		go func(workerIdx int) {
			defer wg.Done()

			start, end := workerRange(workerIdx, pr.numWorkers, config.SpinsCount)

			res, err := pr.runWorkerSpins(gameMod, config, start, end, &completed)
			if err != nil {
				errs[workerIdx] = fmt.Errorf("worker %d: %w", workerIdx, err)

				return
			}

			results[workerIdx] = res
		}(i)
	}

	wg.Wait()
	stopProgress()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// runWorkerSpins executes the spin loop for one worker and returns its accumulated result.
// Each spin gets a fresh generator seeded from (mainSeed, spinIdx), so spin i always
// receives the same random numbers regardless of worker count or ordering.
func (pr *parallelRunner) runWorkerSpins(
	gameMod *game.Mod,
	config *simulator.SimulationConfig,
	start, end int,
	completed *atomic.Int64,
) (*result.Result, error) {
	rng := rngtest.New()
	rng.Reseed(config.MainSeed)

	workerParser := parser.NewRoundParser(gameMod.CostMultiplier(), config.SimpleCostMultiplier)

	for i := range end - start {
		rng.Reset(config.MainSeed, start+i)

		spinResult, err := gameMod.Spin(rng)
		if err != nil {
			return nil, fmt.Errorf("game mod spin: %w", err)
		}

		if err := workerParser.Parse(spinResult); err != nil {
			return nil, fmt.Errorf("parse spin: %w", err)
		}

		completed.Add(1)
	}

	return workerParser.GetResult(), nil
}

// startProgressReporter starts a goroutine that polls the atomic counter and updates the
// progress bar on a fixed interval. Returns a stop function that the caller must invoke
// after all workers finish to terminate the reporter goroutine.
func (pr *parallelRunner) startProgressReporter(completed *atomic.Int64) func() {
	stop := make(chan struct{})

	go func() {
		ticker := time.NewTicker(progressTickMs * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Align to the nearest lower multiple of progressUpdateMask+1 so that
				// bar.Progress.Update() passes its internal bitmask check and renders.
				count := int(completed.Load()) &^ progressUpdateMask
				pr.progressBar.Update(count)
			case <-stop:
				return
			}
		}
	}()

	return func() { close(stop) }
}

// workerRange returns the [start, end) spin indices for the given worker.
// The last worker absorbs any remainder from integer division.
func workerRange(workerIdx, numWorkers, totalSpins int) (start, end int) {
	perWorker := totalSpins / numWorkers
	start = workerIdx * perWorker
	end = start + perWorker

	if workerIdx == numWorkers-1 {
		end = totalSpins
	}

	return start, end
}
