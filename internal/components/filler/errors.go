package filler

import "errors"

var (
	// ErrInvalidPairsCoinCount is returned when the drawn number of coins has no
	// matching entry in the jackpot pairs configuration.
	ErrInvalidPairsCoinCount = errors.New("invalid pairs coin count")

	// ErrInvalidType is returned when a value does not match any known type.
	ErrInvalidType = errors.New("invalid spin type feature")
)
