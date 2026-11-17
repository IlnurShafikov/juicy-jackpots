// Package base implements the base game phase logic.
// The base phase represents the regular paid spins in the game.
package base

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/components/collectors"
	"github.com/IlnurShafikov/juicy-jackpots/internal/components/limiter"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result/events"
	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/generator"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

// / SpinPhase represents the base game spin phase with all its components.
// It contains the game mechanics and configurations needed to execute a base spin.
type SpinPhase struct {
	// gen produces raw reel boards for a given group ID.
	gen *generator.Generator

	// dist define the weighted symbol distributions for reels
	dist *weight.Weights

	calc *calculators.Lines
}

// NewSpinPhase creates a new base game spin phase with the specified components.
//
// Parameters:
//   - dist: symbol distributions for the reels
//   - weightsGroup: weighted value groups for game mechanics
//
// Returns:
//   - pointer to a new SpinPhase instance
func NewSpinPhase(
	dist *weight.Weights,
	gen *generator.Generator,
	calc *calculators.Lines,
) *SpinPhase {
	return &SpinPhase{
		dist: dist,
		gen:  gen,
		calc: calc,
	}
}

// Spin executes the base game phase and returns the result.
// This is where all base game mechanics are implemented: reel spins, win calculation,
// feature triggers, etc.
//
// Parameters:
//   - rng: random number generator for deterministic randomness
//
// Returns:
//   - pointer to Base containing the spin result
//   - error if spin execution fails
func (b *SpinPhase) Spin(
	rng model.RNG,
	winLimiter limiter.WinLimiter,
) (*result.Base, error) {
	eventer := collectors.NewEvents()
	res := result.NewBase()

	groupID, err := b.dist.GetRandomValue(rng)
	if err != nil {
		return nil, fmt.Errorf("distributions get random value: %w", err)
	}

	board, err := b.gen.Generate(rng, groupID)
	if err != nil {
		return nil, fmt.Errorf("generate board: %w", err)
	}

	totalWin, err := b.processWin(board, eventer, winLimiter)
	if err != nil {
		return nil, fmt.Errorf("process win: %w", err)
	}

	res.Board = board
	res.SetTotalWin(totalWin)
	res.SetMaxWinReached(winLimiter.MaxWinReached())
	res.SetEvents(eventer.Get())

	return res, nil
}

// processWin evaluates line wins and coin wins for the generated board result,
// records corresponding events, and returns the combined payout.
// When the win limiter is hit, it returns the capped total immediately.
//
// Parameters:
//   - resp: board result containing the symbol grid and collector state
//   - eventer: event collector that accumulates game events for client delivery
//   - limiter: win cap enforcing the maximum allowed payout
//
// Returns:
//   - combined total of line wins and coin wins, capped by the limiter if reached
//   - error if line calculation or event creation fails
func (b *SpinPhase) processWin(
	board model.Board,
	eventer *collectors.EventsCollector,
	winLimiter limiter.WinLimiter,
) (uint64, error) {
	lineWins, err := b.calc.CalculateWithStrategies(
		board, calculators.WithLinesWinLimiter(winLimiter),
	)
	if err != nil {
		return 0, fmt.Errorf("calculate win lines: %w", err)
	}

	totalWin := calculateWinAmount(lineWins)

	if len(lineWins) != 0 {
		lineWinEvent, err := events.NewLineWinsEvent(lineWins)
		if err != nil {
			return 0, fmt.Errorf("create line wins event: %w", err)
		}

		eventer.Add(lineWinEvent)
	}

	if winLimiter.MaxWinReached() {
		return winLimiter.TotalWin(), nil
	}

	return totalWin, nil
}

func calculateWinAmount(
	wins calculators.LineWins,
) uint64 {
	if len(wins) == 0 {
		return 0
	}

	var totalWin uint64

	for _, win := range wins {
		totalWin += win.Payout
	}

	return totalWin
}
