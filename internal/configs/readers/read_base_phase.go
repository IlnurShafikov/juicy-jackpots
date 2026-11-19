package readers

import (
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/components/filler"
	"github.com/IlnurShafikov/juicy-jackpots/internal/components/manager"
	"github.com/IlnurShafikov/juicy-jackpots/internal/configs"
	"github.com/IlnurShafikov/juicy-jackpots/internal/configs/readers/common"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/rules"
	"github.com/IlnurShafikov/juicy-jackpots/internal/game/phase/base"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/generator"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// readBasePhase reads and constructs the base game phase configuration.
// It loads symbol distributions and weight groups from CSV files.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - modPath: path builder pointing to the game mode directory
//   - general: general configuration shared across all phases
//
// Returns:
//   - pointer to SpinPhase containing base game configuration
//   - error if reading or parsing fails
func readBasePhase(
	reader *reader.Adapter, modPath *tools.PathBuilder, general *configs.GeneralConfig,
) (*base.SpinPhase, error) {
	const (
		basePrefix = "base"
	)

	// PathBuilder for creating paths to files
	phasePath := modPath.Phase(basePrefix)

	dist, err := common.ReadDistributions(reader, phasePath)
	if err != nil {
		return nil, fmt.Errorf("read distributions: %w", err)
	}

	gen, err := readBaseReelsGenerator(reader, phasePath, general, dist.CopyValues())
	if err != nil {
		return nil, fmt.Errorf("read base reel generator: %w", err)
	}

	fillerCoin, err := readHuntFiller(reader, phasePath, general)
	if err != nil {
		return nil, fmt.Errorf("read filler hunt: %w", err)
	}

	lineCalculator := readLineCalculator(general)

	jpWinManager := manager.NewJpCoinWin()

	spinPhase := base.NewSpinPhase(dist, gen, lineCalculator, fillerCoin, jpWinManager)

	return spinPhase, nil
}

// readHuntFiller reads the jackpot hunt coin count weight and constructs the Hunt filler.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - phasePath: path builder pointing to the phase configuration directory
//   - general: general configuration containing jackpot pairs and type mappings
//
// Returns:
//   - pointer to Hunt filler configured with coin count weights and jackpot data
//   - error if reading the weight file fails
func readHuntFiller(
	reader *reader.Adapter,
	phasePath *tools.PathBuilder,
	general *configs.GeneralConfig,
) (*filler.JpCoins, error) {
	const jpCoinCountFileName = "jp_coin_count"

	jpHuntCoinCount, err := common.ReadWeight(reader, phasePath.File(jpCoinCountFileName))
	if err != nil {
		return nil, fmt.Errorf("read jp hunt coin count weight: %w", err)
	}

	huntFiller := filler.NewJpCoins(jpHuntCoinCount, general.JackpotPairs, general.JpTypes)

	return huntFiller, nil
}

// readLineCalculator constructs the line-win calculator from the shared general
// configuration. Wild and non-scoring symbols (coin, collector) are registered
// so the evaluator handles substitutions and skips correctly.
//
// Parameters:
//   - general: shared game configuration supplying paytable, win lines, and symbol IDs
//
// Returns:
//   - pointer to Calculator ready to evaluate line wins on base-game boards
func readLineCalculator(
	general *configs.GeneralConfig,
) *calculators.Lines {
	wild := model.Collection{rules.WildID}
	skipSymbol := model.Collection{}

	linesCalc := calculators.NewLines(
		general.PayTable,
		general.WinLines,
		wild,
		skipSymbol)

	return linesCalc
}

// readBaseReelsGenerator reads reel-strip data for every trigger group and
// builds the standard board generator used for base-game spins.
//
// Parameters:
//   - reader: CSV file adapter backed by the embedded config FS
//   - phasePath: path builder pointing to the base-game phase directory
//   - general: shared game configuration supplying the per-reel symbol counts
//   - groupIDs: trigger group IDs for which reel files must be loaded
//
// Returns:
//   - pointer to Generator with all group reel strips loaded
//   - error if any reel file cannot be read or the generator cannot be initialised
func readBaseReelsGenerator(
	reader *reader.Adapter,
	phasePath *tools.PathBuilder,
	general *configs.GeneralConfig,
	groupIDs []model.GroupID,
) (*generator.Generator, error) {
	reels, err := common.ReadReels(reader, phasePath, groupIDs)
	if err != nil {
		return nil, fmt.Errorf("read reels: %w", err)
	}

	reelGen, err := generator.NewReelsGenerator(reels, general.SymbolsOnReels)
	if err != nil {
		return nil, fmt.Errorf("make base board generator: %w", err)
	}

	return reelGen, nil
}
