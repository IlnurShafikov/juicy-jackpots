package cast

import (
	"math"
)

// SafeCast converts a numeric value from one type to another with validation.
// It checks for NaN, Infinity, overflow, and underflow. When converting from
// floating-point to integer types, it also validates that no fractional part exists.
//
// Parameters:
//   - value: The source numeric value to convert.
//
// Returns:
//   - To: The converted value of the target type.
//   - Error: ErrNaN, ErrInfinity, ErrOverflow, or ErrUnderflow if validation fails.
//
//nolint:ireturn
func SafeCast[To Number, From Number](value From) (To, error) {
	var zero To

	// Convert to float64 for uniform comparison across all numeric types.
	floatValue := float64(value)

	// Check for NaN values which cannot be represented in integer types.
	if math.IsNaN(floatValue) {
		return zero, ErrNaN
	}

	// Check for Infinity values which cannot be represented in finite numeric types.
	if math.IsInf(floatValue, 0) {
		return zero, ErrInfinity
	}

	minValue, maxValue := getMinMax[To]()

	// Check for underflow: value is smaller than the target type can represent.
	if floatValue < minValue {
		return zero, ErrUnderflow
	}

	// Check for overflow: value is larger than the target type can represent.
	if floatValue > maxValue {
		return zero, ErrOverflow
	}

	// Check for precision loss when converting from float to integer.
	// Ensure the float value has no fractional part before casting.
	if isInteger[To]() && isFloat[From]() {
		if floatValue != math.Trunc(floatValue) {
			return zero, ErrOverflow
		}
	}

	// Conversion is safe after all validation checks.
	return To(value), nil
}
