// Package cast provides generic, type-safe conversion utilities for numeric types.
// It supports safe casting between all integer and floating-point types with explicit
// error reporting for NaN, Infinity, overflow, and underflow conditions.
package cast

import (
	"errors"
)

var (
	// ErrNaN is returned when attempting to cast a NaN (Not a Number) value.
	ErrNaN = errors.New("cannot cast NaN value")

	// ErrInfinity is returned when attempting to cast an Infinity value.
	ErrInfinity = errors.New("cannot cast Infinity value")

	// ErrOverflow is returned when a value exceeds the maximum value for the target type.
	ErrOverflow = errors.New("overflow: value exceeds target type maximum")

	// ErrUnderflow is returned when a value is smaller than the target type can represent.
	ErrUnderflow = errors.New("underflow: value is below target type minimum")

	// ErrInvalidFormat indicates the string cannot be parsed as a number.
	ErrInvalidFormat = errors.New("invalid number format")
)
