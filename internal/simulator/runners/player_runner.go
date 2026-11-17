//nolint:all // it's ok for simulator
package runners

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/report"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools/rngtest"
)

// PlayerSimulationRunner executes player-balance simulation.
// Each player starts with a given balance and spins until funds are exhausted.
// Statistics are collected per balance group and returned for report generation.
type PlayerSimulationRunner struct {
	// simpleCostMultiplier scales the CLI-provided balance values to in-game units.
	simpleCostMultiplier int
}

// NewPlayerSimulationRunner creates a new PlayerSimulationRunner with the specified cost multiplier.
//
// Parameters:
//   - simpleCost: the simple-mode cost multiplier used to scale player balances
//
// Returns:
//   - pointer to a new PlayerSimulationRunner instance
func NewPlayerSimulationRunner(simpleCost int) *PlayerSimulationRunner {
	return &PlayerSimulationRunner{
		simpleCostMultiplier: simpleCost,
	}
}

// RunPlayerSimulation runs a complete player-balance simulation for every configured balance group.
// For each balance value, a group of players is simulated; each player spins until broke.
//
// Parameters:
//   - gameMod: the game mod to use for spins
//   - config: simulation configuration including balances, players count, and RNG seed
//
// Returns:
//   - slice of PlayerStats, one entry per balance group, sorted by the order of config.PlayerBalances
//   - error if any individual spin or group simulation fails
func (psr *PlayerSimulationRunner) RunPlayerSimulation(
	gameMod *game.Mod, config *simulator.SimulationConfig,
) ([]*report.PlayerStats, error) {
	rng := psr.initializeRNG(config.MainSeed)

	costMultiplier := gameMod.CostMultiplier()

	roundCost := costMultiplier * uint64(config.BetValue)

	actualBalances := psr.multiplyBalances(config.PlayerBalances, uint64(psr.simpleCostMultiplier))

	allStats := make([]*report.PlayerStats, 0, len(config.PlayerBalances))

	for _, balance := range actualBalances {
		groupStats, err := psr.simulatePlayerGroup(
			gameMod, rng, roundCost, balance, config)
		if err != nil {
			return nil, fmt.Errorf("simulate group with balance %d: %w", balance, err)
		}

		allStats = append(allStats, groupStats)
	}

	return allStats, nil
}

// simulatePlayerGroup simulates one group of players sharing the same initial balance.
// Returns aggregated statistics for the group after all players have exhausted their funds.
func (psr *PlayerSimulationRunner) simulatePlayerGroup(
	gameMod *game.Mod, rng *rngtest.Mt19937,
	roundCost, initialBalance uint64, simConfig *simulator.SimulationConfig,
) (*report.PlayerStats, error) {
	group := report.NewPlayerGroup(initialBalance, simConfig.PlayersCount)

	for i, player := range group.Players {
		if err := psr.simulatePlayer(
			gameMod, rng, player, roundCost, simConfig.BetValue,
		); err != nil {
			return nil, fmt.Errorf("simulate player %d: %w", i, err)
		}
	}

	return report.NewPlayerStats(group), nil
}

// simulatePlayer simulates a single player, executing spins until the balance drops below roundCost.
func (psr *PlayerSimulationRunner) simulatePlayer(
	gameMod *game.Mod, rng *rngtest.Mt19937, player *report.Player, roundCost uint64, betValue int,
) error {
	for player.CanSpin(roundCost) {
		if err := psr.executePlayerSpin(gameMod, rng, player, roundCost, betValue); err != nil {
			return fmt.Errorf("execute player spin: %w", err)
		}
	}

	return nil
}

// executePlayerSpin deducts the bet, performs one spin, and credits any winnings.
func (psr *PlayerSimulationRunner) executePlayerSpin(
	gameMod *game.Mod, rng *rngtest.Mt19937, player *report.Player, roundCost uint64, betValue int,
) error {
	player.MakeBet(roundCost)
	player.IncrementSpins()

	spinResult, err := gameMod.Spin(rng)
	if err != nil {
		return fmt.Errorf("game spin: %w", err)
	}

	totalWin := psr.calculateTotalWin(spinResult) * uint64(betValue)
	player.AddWin(totalWin)

	return nil
}

// calculateTotalWin sums the base-phase win and the last free-phase cumulative win.
// The last free-phase step carries the cumulative total, so earlier steps are not added separately.
func (psr *PlayerSimulationRunner) calculateTotalWin(spinResult *result.Round) uint64 {
	totalWin := spinResult.BasePhase.TotalWin()

	return totalWin
}

// initializeRNG creates a Mersenne Twister RNG seeded with mainSeed.
func (psr *PlayerSimulationRunner) initializeRNG(mainSeed int64) *rngtest.Mt19937 {
	rng := rngtest.New()
	rng.Reseed(mainSeed)

	return rng
}

// multiplyBalances scales each balance value by multiplier to convert display units to game units.
func (psr *PlayerSimulationRunner) multiplyBalances(
	balances []uint64, multiplier uint64,
) []uint64 {
	updatedBalances := make([]uint64, len(balances))

	for i, val := range balances {
		updatedBalances[i] = val * multiplier
	}

	return updatedBalances
}
