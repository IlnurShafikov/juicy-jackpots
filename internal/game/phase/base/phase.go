// Package base implements the base game phase logic.
// The base phase represents the regular paid spins in the game.
package base

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/components/collectors"
	"github.com/IlnurShafikov/juicy-jackpots/internal/components/filler"
	"github.com/IlnurShafikov/juicy-jackpots/internal/components/limiter"
	"github.com/IlnurShafikov/juicy-jackpots/internal/components/manager"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/coin"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result/events"
	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/generator"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

// SpinPhase represents the base game spin phase with all its components.
// It contains the game mechanics and configurations needed to execute a base spin.
type SpinPhase struct {
	// gen produces raw reel boards for a given group ID.
	gen *generator.Generator

	// dist define the weighted symbol distributions for reels
	dist *weight.Weights

	// calc evaluates the line wins on the board.
	calc *calculators.Lines

	// jpCoinFiller places coin tokens on the board for the jackpot coins feature
	jpCoinFiller *filler.JpCoins

	// jpCoinWinManager evaluates coin results and emits jackpot win events
	jpCoinWinManager *manager.JpCoinWin
}

// NewSpinPhase creates a new base game spin phase with the specified components.
//
// Parameters:
//   - dist: symbol distributions for the reels
//   - gen: board generator built from the reel strips
//   - calc: line wins calculator
//   - jpCoinFiller: draws the jackpot coins of the spin
//   - jpCoinWinManager: converts the coins into a jackpot win
//
// Returns:
//   - pointer to a new SpinPhase instance
func NewSpinPhase(
	dist *weight.Weights,
	gen *generator.Generator,
	// calc evaluates the line wins on the board.
	calc *calculators.Lines,
	jpCoinFiller *filler.JpCoins,
	jpCoinWinManager *manager.JpCoinWin,
) *SpinPhase {
	return &SpinPhase{
		dist:             dist,
		gen:              gen,
		calc:             calc,
		jpCoinFiller:     jpCoinFiller,
		jpCoinWinManager: jpCoinWinManager,
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

	coinRes, err := b.jpCoinFiller.Fill(rng, board)
	if err != nil {
		return nil, fmt.Errorf("fill coins: %w", err)
	}

	totalWin, err := b.processWin(board, coinRes, eventer, winLimiter)
	if err != nil {
		return nil, fmt.Errorf("process win: %w", err)
	}

	res.Board = board
	res.CoinsResult = coinRes
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
//   - board: generated board to evaluate
//   - coinRes: jackpot coins drawn for this spin
//   - eventer: event collector that accumulates game events for client delivery
//   - winLimiter: win cap enforcing the maximum allowed payout
//
// Returns:
//   - combined total of line wins and coin wins, capped by the limiter if reached
//   - error if line calculation or event creation fails
func (b *SpinPhase) processWin(
	board model.Board,
	coinRes *coin.Result,
	eventer *collectors.EventsCollector,
	winLimiter limiter.WinLimiter,
) (uint64, error) {
	lineWins, err := b.calc.CalculateWithStrategies(
		board, calculators.WithLinesWinLimiter(winLimiter),
	)
	if err != nil {
		return 0, fmt.Errorf("calculate win lines: %w", err)
	}

	if len(lineWins) != 0 {
		lineWinEvent, err := events.NewLineWinsEvent(lineWins)
		if err != nil {
			return 0, fmt.Errorf("create line wins event: %w", err)
		}

		eventer.Add(lineWinEvent)
	}

	b.jpCoinWinManager.Calculate(coinRes, winLimiter, eventer)

	return winLimiter.TotalWin(), nil
}
