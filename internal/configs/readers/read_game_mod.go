package readers

import (
	"fmt"
	"strconv"

	"github.com/IlnurShafikov/juicy-jackpots/internal/configs"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/rules"
	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

const (
	// rtpFileName is the name of the CSV file containing the expected RTP value
	rtpFileName = "exp_rtp"

	// simpleModPrefix is the directory name for simple mode configurations
	simpleModPrefix = "simple"
)

// gameModsData holds all game mode configurations for a specific RTP version.
type gameModsData struct {
	// simple is the simple (single bet) game mode.
	simple *game.Mod
}

// ReadGameMods reads and constructs all game modes for a specific RTP version.
// It loads configurations for each mode, creates game rules, and packages everything together.
//
// Parameters:
//   - rtpVersion: the RTP version number (e.g., 96 for 96%)
//
// Returns:
//   - pointer to Mods containing all game modes and rules for this RTP version
//   - error if reading or constructing fails
func ReadGameMods(rtpVersion int) (*game.Mods, error) {
	csvReader := reader.NewAdapter(configs.ConfigFS, configs.RootPath)

	strRtpVersion := strconv.Itoa(rtpVersion)
	rtpPath := tools.NewPathBuilder("").RTP(strRtpVersion)

	generalCfg, err := Read(csvReader)
	if err != nil {
		return nil, fmt.Errorf("get generalcfg config: %w", err)
	}

	allGameMods, err := readAllGameMods(csvReader, rtpPath, generalCfg)
	if err != nil {
		return nil, fmt.Errorf("get all game game: %w", err)
	}

	expRTP, err := csvReader.ReadFloat64(rtpPath.File(rtpFileName))
	if err != nil {
		return nil, fmt.Errorf("get exp RTP: %w", err)
	}

	gameRules, err := rules.MakeGameRules(generalCfg, expRTP)
	if err != nil {
		return nil, fmt.Errorf("make game rules: %w", err)
	}

	gameMods := game.NewMods(
		allGameMods.simple,
		gameRules,
		rtpVersion,
	)

	return gameMods, nil
}

// readAllGameMods reads configurations for all available game modes.
// Currently supports the simple mode, with placeholders for additional modes.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - rtpPath: path builder pointing to the RTP version directory
//   - generalCfg: general configuration shared across all modes
//
// Returns:
//   - pointer to gameModsData containing all mode configurations
//   - error if reading any mode fails
func readAllGameMods(
	csvReader *reader.Adapter,
	rtpPath *tools.PathBuilder,
	generalCfg *configs.GeneralConfig,
) (*gameModsData, error) {
	// TODO: Read other mods here if needed
	simpleMod, err := readGameMod(csvReader, rtpPath.Mod(simpleModPrefix),
		generalCfg.SimpleCostMultiplier(), generalCfg)
	if err != nil {
		return nil, fmt.Errorf("get simple mod config: %w", err)
	}

	// TODO: Example for other mods
	// buy3ScatterFsMod, err := readGameMod(csvReader, rtpPath.Mod(buyFs3ScatterModPrefix),
	// 	generalCfg.BuyFs3scatterCostMultiplier, generalCfg)
	// if err != nil {
	// 	return nil, fmt.Errorf("get buyFs3Scatter mod config: %w", err)
	// }

	return &gameModsData{
		simple: simpleMod,
	}, nil
}

// readGameMod reads the configuration for a single game mode.
// It loads both base and free phase configurations and creates a complete Mod.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - modPath: path builder pointing to the mode directory
//   - costMultiplier: cost multiplier for this mode
//   - generalCfg: general configuration shared across all modes
//
// Returns:
//   - pointer to Mod containing the complete mode configuration
//   - error if reading any phase fails
func readGameMod(
	csvReader *reader.Adapter,
	modPath *tools.PathBuilder,
	costMultiplier uint64,
	generalCfg *configs.GeneralConfig,
) (*game.Mod, error) {
	basePhase, err := readBasePhase(csvReader, modPath, generalCfg)
	if err != nil {
		return nil, fmt.Errorf("read base config: %w", err)
	}

	gameMod := game.NewMod(basePhase, costMultiplier, generalCfg.MaxWinLimit())

	return gameMod, nil
}
