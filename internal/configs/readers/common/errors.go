// Package common provides shared utilities for reading and parsing configuration files.
// It contains common functions used across different configuration readers.
package common

import (
	"errors"
)

var (
	// ErrEmptyData is returned when attempting to parse empty or invalid CSV data.
	ErrEmptyData = errors.New("empty data")

	// ErrInvalidData is returned when attempting to parse invalid CSV data.
	ErrInvalidData = errors.New("invalid data")

	// ErrMismatchedValuesAndWeights is returned when the number of
	// values and weights don't match in the configuration.
	ErrMismatchedValuesAndWeights = errors.New("mismatched values and weights count")
)
