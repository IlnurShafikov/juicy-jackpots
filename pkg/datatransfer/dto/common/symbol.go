package common

import "github.com/IlnurShafikov/slot-engine/pkg/model"

// Symbol represents a single symbol on the game board with its position and optional properties.
// This structure is sent to the client to display the game board.
type Symbol struct {
	// Property contains optional symbol-specific properties (e.g., multipliers, special behaviors)
	Property *SymbolProperties `json:"properties,omitempty"`

	// Position indicates where this symbol appears on the board (reel and row indices)
	Position model.Position `json:"position"`

	// SymbolID is the identifier of the symbol type
	SymbolID int `json:"symbol"`
}
