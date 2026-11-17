package readers_test

import (
	"testing"
	"time"

	"github.com/IlnurShafikov/juicy-jackpots/internal/configs/readers"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools/rngtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadGameMods(t *testing.T) {
	t.Parallel()

	// Test reading game for RTP 96
	gameMods, err := readers.ReadGameMods(readers.RTP96)
	require.NoError(t, err)
	assert.NotNil(t, gameMods)
}

func TestSpin(t *testing.T) {
	t.Parallel()

	const (
		modID  = "simple"
		featID = ""
	)

	rng := rngtest.New()

	// Test reading game for RTP 96
	gameMods, err := readers.ReadGameMods(readers.RTP96)
	require.NoError(t, err)
	assert.NotNil(t, gameMods)

	gameMod, err := gameMods.GetGameMod(modID, featID)
	require.NoError(t, err)

	for range 10000 {
		rng.Reseed(int64(time.Now().Nanosecond()))
		result, err := gameMod.Spin(rng)
		require.NoError(t, err)
		assert.NotNil(t, result)
	}
}
