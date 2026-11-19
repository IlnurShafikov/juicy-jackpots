package datatransfer

import (
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/common"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// TransformSymbols converts the internal game board representation into client DTO symbols.
// The board is converted cell by cell, reel by reel, into
// an array of Symbol DTOs that will be sent to the client.
//
// Parameters:
//   - gameBoard: the internal game board
//   - betValue: the player's bet value (not used, symbols have no properties)
//
// Returns:
//   - slice of Symbol DTOs representing the board
//   - error if transformation fails
func TransformSymbols(gameBoard model.Board, _ uint64) ([]*common.Symbol, error) {
	if gameBoard == nil {
		return nil, ErrGameBoardIsNil
	}

	reelsCount := len(gameBoard)
	rowsCount := len(gameBoard[0])

	resp := make([]*common.Symbol, 0, reelsCount*rowsCount)

	for reel := range reelsCount {
		for row := range rowsCount {
			resp = append(resp, &common.Symbol{
				SymbolID: gameBoard[reel][row],
				Position: model.NewPosition(reel, row),
				Property: nil,
			})
		}
	}

	return resp, nil
}
