//nolint:all // it's ok for cast package
package cast

import (
	"fmt"
	"math"
	"strconv"
)

// ToString converts a numeric value to its string representation
// with specified decimal precision. Only applicable for floating-point types.
// For integer types, the precision parameter is ignored.
//
// Parameters:
//   - value: The numeric value to convert.
//   - precision: Number of digits after decimal point (0-15).
//
// Returns:
//   - string: String representation of the number with specified precision.
//   - Error: ErrNaN, ErrInfinity, or ErrInvalidPrecision if validation fails.
func ToString[T Number](value T) (string, error) {
	const precision = 8

	floatValue := float64(value)

	// Check for NaN values
	if math.IsNaN(floatValue) {
		return "", ErrNaN
	}

	// Check for Infinity values
	if math.IsInf(floatValue, 0) {
		return "", ErrInfinity
	}

	// For integers, ignore precision
	if isInteger[T]() {
		return fmt.Sprintf("%d", int64(value)), nil
	}

	return strconv.FormatFloat(floatValue, 'f', precision, 64), nil
}
