package filler

import (
	"fmt"
	"slices"

	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/coin"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/pairs"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/rules"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools/cast"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/weight"
)

// JpCoins places coins on the board for the jackpot hunt feature.
// It draws a random number of coins from a weighted distribution, positions them on
// randomly selected reels, and assigns a jackpot type and value based on the coin count.
type JpCoins struct {
	// coinCount holds the weighted distribution of how many coins to place per spin.
	coinCount *weight.Weights
	// jpTypes maps a jackpot value to its corresponding jackpot type ID.
	jpTypes map[uint64]int
	// jpPairs maps a coin count to its corresponding jackpot payout value.
	jpPairs pairs.Pairs
}

// NewJpCoins creates a Hunt filler with coin distribution and jackpot configuration.
//
// Parameters:
//   - coinCount: weighted distribution controlling how many coins appear per spin
//   - jpPairs: mapping from coin count to jackpot payout value
//   - jpTypes: mapping from jackpot value to jackpot type ID
//
// Returns:
//   - pointer to a new Hunt filler
func NewJpCoins(
	coinCount *weight.Weights,
	jpPairs pairs.Pairs,
	jpTypes map[uint64]int,
) *JpCoins {
	return &JpCoins{
		coinCount: coinCount,
		jpPairs:   jpPairs,
		jpTypes:   jpTypes,
	}
}

// Fill places coins on the board for a single spin and returns the coin result.
// Returns a zero-value Coins result (no positions) when the drawn coin count is zero.
//
// Parameters:
//   - rng: random number generator used for count sampling and reel selection
//   - board: current game board used to determine reel count
//
// Returns:
//   - pointer to Coins containing positions, jackpot value, and type
//   - error if weight sampling, reel selection, or type cast fails
func (h *JpCoins) Fill(
	rng model.RNG,
	board model.Board,
) (*coin.Result, error) {
	coinsRes := coin.NewResult()

	coinCount, err := h.coinCount.GetRandomValue(rng)
	if err != nil {
		return nil, fmt.Errorf("get random value: %w", err)
	}

	if coinCount == 0 {
		return coinsRes, nil
	}

	reelCount := board.ReelCount()

	randReels, err := tools.RandomIntSlices(rng, reelCount, coinCount)
	if err != nil {
		return nil, fmt.Errorf("random int slices: %w", err)
	}

	slices.Sort(randReels)

	coinPos := buildCoinPositions(rng, randReels)
	coinsRes.Positions = coinPos

	err = h.coinsProperties(coinsRes, coinCount)
	if err != nil {
		return nil, fmt.Errorf("coins properties: %w", err)
	}

	return coinsRes, nil
}

// coinsProperties resolves the jackpot value and type for the given coin count
// and writes them into coinsRes.
//
// Parameters:
//   - coinsRes: coin result to populate with value and jackpot type
//   - coinCount: number of coins placed this spin, used to look up the jackpot pair
//
// Returns:
//   - [ErrInvalidPairsCoinCount] if coinCount has no matching jackpot pair
//   - error if the jackpot value cannot be cast to uint64
func (h *JpCoins) coinsProperties(coinsRes *coin.Result, coinCount int) error {
	value, ok := h.jpPairs.GetValueByKey(coinCount)
	if !ok {
		return ErrInvalidPairsCoinCount
	}

	jpValue, err := cast.SafeCast[uint64](value)
	if err != nil {
		return fmt.Errorf("safe cast int to uint64: %w", err)
	}

	jpType := h.jpTypes[jpValue]

	coinsRes.Value = jpValue
	coinsRes.JpType = jpType

	return nil
}

// buildCoinPositions assigns a random row to each reel in reels.
// If a reel has a nudge position the coin is placed on the middle row instead.
//
// Parameters:
//   - rng: random number generator used to select a row per reel
//   - reels: 0-based reel indices where coins are placed
//
// Returns:
//   - positions of the placed coins, one per reel
func buildCoinPositions(
	rng model.RNG,
	reels []int,
) model.Positions {
	coinPos := make(model.Positions, 0, len(reels))

	for _, reel := range reels {
		randRow := rng.Random(model.MinRNGValue, rules.MaxRowIndex)

		coinPos = append(coinPos, model.NewPosition(reel, int(randRow)))
	}

	return coinPos
}
