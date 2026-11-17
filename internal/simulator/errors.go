// Package simulator provides tools for running game simulations, including
// simple batch mode, interactive step-by-step mode, and player balance simulation.
// It is intended for internal use during game development and RTP verification.
package simulator

import (
	"errors"
	"fmt"
)

var (
	// ErrGameNotFound is returned when a requested game mod or RTP version
	// does not exist in the loaded configuration.
	ErrGameNotFound = errors.New("game not found")

	// ErrInvalidRngType is returned when an unrecognised RNG type string is provided.
	// Supported values are "mt" (Mersenne Twister) and "fortuna".
	ErrInvalidRngType = errors.New("invalid RNG type")
)

// SimulationError wraps an underlying error with the name of the failed operation.
// It provides contextual information about where in the simulation pipeline
// the failure occurred.
type SimulationError struct {
	// Cause is the underlying error that caused the simulation to fail.
	Cause error
	// Operation is a human-readable name for the step that failed.
	Operation string
}

// Error implements the error interface and returns a formatted error message.
func (e SimulationError) Error() string {
	return fmt.Sprintf("simulation %s failed: %v", e.Operation, e.Cause)
}

// NewSimulationError creates a new SimulationError for the given operation and cause.
//
// Parameters:
//   - operation: human-readable name of the step that failed (e.g. "execute spin loop")
//   - cause: the underlying error
//
// Returns:
//   - error wrapping the cause with the operation context
func NewSimulationError(operation string, cause error) error {
	return SimulationError{Operation: operation, Cause: cause}
}
