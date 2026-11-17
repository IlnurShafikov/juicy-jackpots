// Package tools provides helper utilities and general-purpose tools
// for working with game logic, including error handling, data conversion, and path operations.
package tools

import (
	"errors"
)

var (
	// ErrIndexOutOfBounds is returned when attempting to access an element
	// outside the valid index range.
	ErrIndexOutOfBounds = errors.New("index out of bounds")

	// ErrEmptyPositions is returned when working with an empty positions slice
	// when at least one position is expected.
	ErrEmptyPositions = errors.New("empty positions slice")
)
