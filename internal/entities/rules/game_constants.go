// Package rules contain game rules, constants, and configuration structures
// that define how the game operates, including symbol IDs, paytables, and game mechanics.
package rules

// Symbol ID constants define unique identifiers for each symbol in the game.
// TODO: implement symbols ids according to your game design
const (
	// VoidSymbol represents an empty or null symbol position
	VoidSymbol = 0

	// WildID represents the wild symbol that can substitute for other symbols
	WildID = 1

	// High1 represents the first high-paying symbol
	High1 = 2

	// High2 represents the second high-paying symbol
	High2 = 3

	// Med1 represents the first medium-paying symbol
	Med1 = 4

	// Med2 represents the second medium-paying symbol
	Med2 = 5

	// Ace represents the Ace card symbol (low-paying)
	Ace = 6

	// King represents the King card symbol (low-paying)
	King = 7

	// Queen represents the Queen card symbol (low-paying)
	Queen = 8

	// Jack represents the Jack card symbol (low-paying)
	Jack = 9
)

// Game mechanics constants
const (
	// InitRespinsCount defines the initial number of respins when the respin feature is triggered
	InitRespinsCount    = 3
	Step_TYPE_BASE_SPIN = 1
	Step_TYPE_FREE_SPIN = 2
	Step_TYPE_RE_SPIN   = 3
)
