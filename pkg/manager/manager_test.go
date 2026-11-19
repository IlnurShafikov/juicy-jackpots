package manager_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/rules"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/step"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/manager"
)

const (
	rtp96 = 96
	bet   = 2
)

// TestGameManager_GetRoundResult tests the GetRoundResult method of the GameManager.
func TestGameManager_GetRoundResult(t *testing.T) {
	gameManager, err := manager.NewGameManager()
	require.NoError(t, err)

	for range 1000 {
		spinRes, err := gameManager.GetRoundResult(
			context.Background(), rules.SimpleModID, rules.EmptyFeatureID, rtp96, seed(), bet)
		require.NoError(t, err)
		require.NotNil(t, spinRes)
	}
}

// TestGameManager_AllRtpGetRoundResult tests the GetRoundResult method for all rtp versions
// of the GameManager.
func TestGameManager_AllRtpGetRoundResult(t *testing.T) {
	rtps := []int{rtp96}

	gameManager, err := manager.NewGameManager()
	require.NoError(t, err)

	for _, curRtp := range rtps {
		for range 100 {
			spinRes, err := gameManager.GetRoundResult(
				context.Background(), rules.SimpleModID, rules.EmptyFeatureID, curRtp, seed(), bet)
			require.NoError(t, err)
			require.NotNil(t, spinRes)
		}
	}
}

// TestEnterGame tests the EnterGame method of the GameManager.
// It ensures that the method can be called without errors and returns a valid response.
func TestEnterGame(t *testing.T) {
	gm, err := manager.NewGameManager()
	require.NoError(t, err)

	ctx := context.Background()
	enterGame, err := gm.EnterGame(ctx, bet, rtp96)
	require.NoError(t, err)
	require.NotNil(t, enterGame)
}

// TestAllGameMods tests all game to ensure they can be processed without errors.
func TestAllGameMods(t *testing.T) {
	gameManager, err := manager.NewGameManager()
	require.NoError(t, err)

	allGameMods := []string{
		rules.SimpleModID,
	}

	for _, modID := range allGameMods {
		for range 1000 {
			spinRes, err := gameManager.GetRoundResult(
				context.Background(), modID, rules.EmptyFeatureID, rtp96, seed(), bet)
			require.NoError(t, err)
			require.NotNil(t, spinRes)
		}
	}
}

func TestGameManager_GetRoundResult_Payout(t *testing.T) {
	gameManager, err := manager.NewGameManager()
	require.NoError(t, err)

	for range 1000 {
		gameSeed := seed()

		beforeDT, err := gameManager.GenerateRoundResult(rtp96, rules.SimpleModID, rules.EmptyFeatureID, gameSeed)
		require.NoError(t, err)
		require.NotNil(t, beforeDT)

		afterDT, err := gameManager.GetRoundResult(
			context.Background(), rules.SimpleModID, rules.EmptyFeatureID, rtp96, gameSeed, bet)
		require.NoError(t, err)
		require.NotNil(t, afterDT)

		rawTotalWin := calculateRawTotalWin(beforeDT)
		transferTotalWin := calculateTransferTotalWin(afterDT.Steps)

		resMsg := fmt.Sprintf(
			"mismatch win %d need %d, seed=%d",
			transferTotalWin, rawTotalWin*bet, afterDT.Seed,
		)

		require.Equal(t, rawTotalWin*bet, transferTotalWin, resMsg)
	}
}

// calculateRawTotalWin calculates the total win from the internal game result
// before it's been transformed or scaled by bet value.
//
// Parameters:
//   - roundRes: the internal round result from the game engine
//
// Returns:
//   - total win amount (in base units, not scaled by bet)
func calculateRawTotalWin(roundRes *result.Round) uint64 {
	var totalWin uint64

	totalWin += roundRes.BasePhase.TotalWin()

	return totalWin
}

// calculateTransferTotalWin calculates the total win from the transformed
// client-facing steps. This value should be scaled by the bet value.
//
// Parameters:
//   - allSteps: slice of transformed steps sent to the client
//
// Returns:
//   - total win amount (already scaled by bet value)
func calculateTransferTotalWin(allSteps []*step.Step) uint64 {
	var totalWin uint64

	for _, step := range allSteps {
		totalWin += step.TotalWin
	}

	return totalWin
}

// seed generates a non-zero seed based on the current time.
// This ensures each test run uses different random seeds.
//
// Returns:
//   - a non-zero int64 seed value for RNG
func seed() int64 {
	var respSeed int64

	for {
		respSeed = int64(time.Now().Nanosecond())

		if respSeed != 0 {
			break
		}
	}

	return respSeed
}
