package common

// SymbolProperties contains optional properties that can be attached to symbols.
// These properties modify the display of symbols.
type SymbolProperties struct {
	// JpToken marks a symbol position that holds a jackpot coin.
	// Omitted from JSON when false.
	JpToken bool `json:"jpToken,omitempty"`
}
