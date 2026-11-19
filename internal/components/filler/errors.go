package filler

import "errors"

var (
	// ErrInvalidPairsCoinCount is returned when the number of coins in a pair configuration
	// does not match the expected count for the hunt filler.
	ErrInvalidPairsCoinCount = errors.New("invalid pairs coin count")

	// ErrInvalidType is returned when the spin type value does not match any known
	// spin classification (rest, negative, or positive).
	ErrInvalidType = errors.New("invalid spin type feature")
)
