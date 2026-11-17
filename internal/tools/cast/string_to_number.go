//nolint:all // it's ok for cast package
package cast

import (
	"math"
	"strconv"
)

// ToNumber converts a string to a numeric type with validation.
// It parses the string and checks for overflow, underflow, and invalid formats.
// For integer types, the string must represent a whole number.
//
// Parameters:
//   - value: The string to convert.
//
// Returns:
//   - To: The converted numeric value.
//   - Error: ErrInvalidFormat, ErrOverflow, or ErrUnderflow if validation fails.
func ToNumber[To Number](value string) (To, error) {
	var zero To

	// Parse string to float64 for uniform validation
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return zero, ErrInvalidFormat
	}

	// Check for NaN and Infinity
	if math.IsNaN(floatValue) {
		return zero, ErrNaN
	}

	if math.IsInf(floatValue, 0) {
		return zero, ErrInfinity
	}

	minValue, maxValue := getMinMax[To]()

	// Check bounds
	if floatValue < minValue {
		return zero, ErrUnderflow
	}

	if floatValue > maxValue {
		return zero, ErrOverflow
	}

	// For integer targets, check for the fractional part
	if isInteger[To]() {
		if floatValue != math.Trunc(floatValue) {
			return zero, ErrOverflow
		}
	}

	return To(floatValue), nil
}
