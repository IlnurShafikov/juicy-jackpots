//nolint:all // it's ok for simulator
package simulator

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// SimulationConfig holds all parameters that control a simulation run.
// Values are populated from CLI flags by ParseConfigs.
type SimulationConfig struct {
	// PlayerBalances is the list of initial player balances used in player-mode simulation.
	// Each value is later multiplied by the simple cost multiplier to get the actual balance.
	PlayerBalances []uint64

	// GameModID identifies the game mode to simulate (e.g. "simple").
	GameModID string

	// FeatureID identifies the optional game feature (e.g. bonus-buy ID). Empty means none.
	FeatureID string

	// Description is an optional human-readable label attached to the simulation output.
	Description string

	// RngType selects the random number generator backend: "mt" or "fortuna".
	RngType string

	// RoundCost is the cost of a single round in base units, populated after the game mod is loaded.
	RoundCost uint64

	// MainSeed is the seed used to initialise the RNG. Defaults to time.Now().UnixNano() when 0.
	MainSeed int64

	// SpinsCount is the total number of spins to execute in a simulation run.
	SpinsCount int

	// RtpVersion is the RTP percentage version to load (e.g. 96 for 96%).
	RtpVersion int

	// PlayersCount is the number of players per balance group in player-mode simulation.
	PlayersCount int

	// BetValue is the bet amount per spin (in the game's smallest unit).
	BetValue int

	// SimpleCostMultiplier is the cost multiplier for the simple game mode.
	// Used to scale player balances to actual in-game amounts.
	SimpleCostMultiplier int

	// WorkersCount is the number of parallel goroutines used in simulation.
	// Values > 1 enable parallel mode; results are reproducible for the same seed + worker count.
	WorkersCount int

	// PlayMod enables interactive (step-by-step) simulation mode when true.
	PlayMod bool

	// PlayerMod enables player balance simulation mode when true.
	PlayerMod bool
}

// ParseConfigs parses CLI flags and returns a populated SimulationConfig.
// If the main seed flag is not provided or is zero, the current Unix nanosecond
// timestamp is used as the seed.
//
// Example usage:
//
//	cfg := simulator.ParseConfigs()
//	if err := runners.NewSimulationManager().RunSimulation(cfg); err != nil {
//	    log.Fatal(err)
//	}
func ParseConfigs() *SimulationConfig {
	spinsCount := flag.Int("spins", DefaultSpinsCount, "spins count")
	rtpVersion := flag.Int("rtp", DefaultRTPVersion, "slot rtpVersion")
	modID := flag.String("mod", "simple", "need '' or '/x100FS'")
	featureID := flag.String("feature", "", "feature ID, e.g. 'x100FS'")
	description := flag.String("description", "", "simulation description")
	mainSeed := flag.Int64("seed", 0, "main seed")
	interactiveMod := flag.Bool("play", false, "play mode")
	rngType := flag.String("rng", DefaultRngType, "mt or fortuna")
	playerMod := flag.Bool("playerMod", false, "player simulation mode")
	balancesStr := flag.String("balance", "", "comma-separated player balances, e.g. '10,20,30'")
	playersCount := flag.Int("playersCount", DefaultPlayersCount, "players count per balance group")
	betValue := flag.Int("bet", DefaultBetValue, "bet value per spin")
	simpleCostMultiplier := flag.Int("costMultiplier", DefaultSimpleCostMultiplier, "simple mode cost multiplier")
	workersCount := flag.Int("workers", DefaultWorkersCount, "number of parallel workers (1 = sequential)")

	flag.Parse()

	balances, err := parseBalances(*balancesStr)
	if err != nil && *playerMod {
		fmt.Printf("Error parsing balances: %v\n", err)

		balances = []uint64{}
	}

	resp := &SimulationConfig{
		Description:          *description,
		SpinsCount:           *spinsCount,
		RtpVersion:           *rtpVersion,
		GameModID:            *modID,
		FeatureID:            *featureID,
		MainSeed:             *mainSeed,
		PlayMod:              *interactiveMod,
		RngType:              *rngType,
		RoundCost:            0,
		PlayerMod:            *playerMod,
		PlayerBalances:       balances,
		PlayersCount:         *playersCount,
		BetValue:             *betValue,
		SimpleCostMultiplier: *simpleCostMultiplier,
		WorkersCount:         *workersCount,
	}

	if resp.MainSeed == 0 {
		log.Print("main seed is not set, using time.Now().UnixNano()")

		resp.MainSeed = time.Now().UnixNano()
	}

	return resp
}

// parseBalances parses a comma-separated string of balance values into a uint64 slice.
// Empty string returns an empty slice without error.
//
// Parameters:
//   - balancesStr: comma-separated balance values, e.g. "100,200,500"
//
// Returns:
//   - slice of parsed uint64 balance values
//   - error if any individual value cannot be parsed as a non-negative integer
func parseBalances(balancesStr string) ([]uint64, error) {
	if balancesStr == "" {
		return []uint64{}, nil
	}

	parts := strings.Split(balancesStr, ",")
	balances := make([]uint64, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		value, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid balance value '%s': %w", trimmed, err)
		}

		balances = append(balances, value)
	}

	return balances, nil
}
