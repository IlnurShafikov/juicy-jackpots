package datatransfer

import (
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/coin"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/common"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// TransformSymbols converts the internal game board representation into client DTO symbols.
// The board is converted cell by cell, reel by reel, into
// an array of Symbol DTOs that will be sent to the client.
//
// Parameters:
//   - gameBoard: the internal game board
//   - coinsResult: jackpot coins of the spin; their positions get JpToken (may be nil)
//
// Returns:
//   - slice of Symbol DTOs representing the board
//   - error if transformation fails
func TransformSymbols(gameBoard model.Board, coinsResult *coin.Result) ([]*common.Symbol, error) {
	if gameBoard == nil {
		return nil, ErrGameBoardIsNil
	}

	reelsCount := gameBoard.ReelCount()
	rowsCount := gameBoard.FirstReelSize()
	coinSet := buildCoinSet(coinsResult)

	resp := make([]*common.Symbol, 0, reelsCount*rowsCount)

	for reel := range reelsCount {
		for row := range rowsCount {
			pos := model.NewPosition(reel, row)
			sym := &common.Symbol{
				SymbolID: gameBoard[reel][row],
				Position: pos,
				Property: nil,
			}

			if _, ok := coinSet[pos]; ok {
				sym.Property = &common.SymbolProperties{JpToken: true}
			}

			resp = append(resp, sym)
		}
	}

	return resp, nil
}

// buildCoinSet converts coin positions into a set for O(1) membership checks.
//
// Parameters:
//   - coinRes: coin result with the positions of the jackpot coins; may be nil
//
// Returns:
//   - set of the positions that hold a coin, or nil if there are none
func buildCoinSet(coinRes *coin.Result) map[model.Position]struct{} {
	if coinRes == nil || len(coinRes.Positions) == 0 {
		return nil
	}

	set := make(map[model.Position]struct{}, len(coinRes.Positions))

	for _, pos := range coinRes.Positions {
		set[pos] = struct{}{}
	}

	return set
}
