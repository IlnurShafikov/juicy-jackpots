// Package game contains the core game logic, including game mode management,
// phases (base, free spins) and execution of game rounds.
package game

import "errors"

var (
	// ErrFeatureNotImplemented is returned when attempting to use an unimplemented game feature.
	// Used when the requested feature ID is not supported by the current game version.
	ErrFeatureNotImplemented = errors.New("feature id not implemented")

	// ErrInvalidGameModId is returned when an invalid game mode identifier is provided.
	// Used when the requested mod ID does not exist in the game configuration.
	ErrInvalidGameModId = errors.New("invalid game mod id")
)
