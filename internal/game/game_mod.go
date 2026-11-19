package game

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/components/limiter"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// BasePhase defines the interface for base game phase operations
type BasePhase interface {
	// Spin executes a single base game spin.
	Spin(rng model.RNG, winLimiter limiter.WinLimiter) (*result.Base, error)
}

// Mod represents a complete game modification: a single base phase that
// evaluates the line wins and the jackpot coins of every spin
type Mod struct {
	// baseSpinPhase handles base game spin logic.
	baseSpinPhase BasePhase
	// costMultiplier is the multiplier applied to bet values.
	costMultiplier uint64
	// maxWinLimit is the maximum win amount allowed per round.
	maxWinLimit uint64
}

// NewMod creates a new Mod instance with the specified phases and limits
func NewMod(
	basePhase BasePhase,
	costMultiplier, maxWinLimit uint64,
) *Mod {
	return &Mod{
		baseSpinPhase:  basePhase,
		costMultiplier: costMultiplier,
		maxWinLimit:    maxWinLimit,
	}
}

// CostMultiplier returns the cost multiplier for this game mod
func (g *Mod) CostMultiplier() uint64 {
	return g.costMultiplier
}

// Spin executes a complete game round and returns the result
func (g *Mod) Spin(rng model.RNG) (*result.Round, error) {
	winLimiter, err := limiter.New(g.maxWinLimit)
	if err != nil {
		return nil, fmt.Errorf("initialize maxWinCollector %w", err)
	}

	roundResult := result.NewRound(g.costMultiplier)

	if err := g.startBasePhase(rng, roundResult, winLimiter); err != nil {
		return nil, fmt.Errorf("execute base phase: %w", err)
	}

	return roundResult, nil
}

// startBasePhase executes the base game phase and updates the round result
func (g *Mod) startBasePhase(
	rng model.RNG, roundResult *result.Round, winLimiter limiter.WinLimiter,
) error {
	baseResult, err := g.baseSpinPhase.Spin(rng, winLimiter)
	if err != nil {
		return fmt.Errorf("spin base game: %w", err)
	}

	roundResult.BasePhase = baseResult

	return nil
}
