// Package rules contain game rules, constants, and configuration structures
// that define how the game operates, including symbol IDs, paytables, and game mechanics.
package rules

// Symbol ID constants define unique identifiers for each symbol in the game.
// The IDs match the strips and the paytable.
const (
	// VoidSymbol represents an empty or null symbol position (reserved, not used on the strips)
	VoidSymbol = 0

	// High1 is the most valuable symbol
	High1 = 1

	// High2 is the second high-pay symbol
	High2 = 2

	// Med1 is the first medium-pay symbol
	Med1 = 3

	// Med2 is the second medium-pay symbol
	Med2 = 4

	// Low1 is the first low-pay symbol
	Low1 = 5

	// Low2 is the second low-pay symbol
	Low2 = 6

	// Low3 is the third low-pay symbol
	Low3 = 7

	// Low4 is the lowest-pay symbol and the only symbol that pays from 2 reels
	Low4 = 8

	// WildID is the wild symbol, it substitutes for symbols 1-8 and appears on reels 2-4 only
	WildID = 9
)

// Game mechanics constants
const (
	// MaxRowIndex is the index of the last row of the board (the board has 3 rows).
	MaxRowIndex = 2

	// Step_TYPE_BASE_SPIN is the type of a base spin step.
	Step_TYPE_BASE_SPIN = 1
	// Step_TYPE_FREE_SPIN is the type of a free spin step (not used in this game).
	Step_TYPE_FREE_SPIN = 2
	// Step_TYPE_RE_SPIN is the type of a respin step (not used in this game).
	Step_TYPE_RE_SPIN = 3
)
