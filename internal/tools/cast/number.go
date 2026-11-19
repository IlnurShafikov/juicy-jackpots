package cast

import (
	"math"
)

// Number defines a constraint for numeric types, including all signed and unsigned
// integer types and floating-point types. The tilde (~) allows matching both
// predeclared types and types with the same underlying type.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// getMinMax returns the minimum and maximum representable values for a numeric type.
// Values are returned as float64 to allow uniform comparison in SafeCast.
//
// Parameters:
//   - T: The numeric type for which to retrieve bounds (provided via type parameter).
//
// Returns:
//   - float64: The minimum value the type can represent.
//   - float64: The maximum value the type can represent.
//   - For unsigned types, the minimum is always 0.
//   - Returns (0, 0) if the type is not recognized (should never happen with Number constraint).
func getMinMax[T Number]() (minVal, maxVal float64) {
	var zero T

	switch any(zero).(type) {
	case int, int8, int16, int32, int64:
		return getSignedMinMax[T]()
	case uint, uint8, uint16, uint32, uint64:
		return getUnsignedMinMax[T]()
	case float32, float64:
		return getFloatMinMax[T]()
	}

	return 0, 0
}

// getSignedMinMax returns the bounds of a signed integer type as float64.
// It returns (0, 0) for any other type.
func getSignedMinMax[T Number]() (minVal, maxVal float64) {
	var zero T

	switch any(zero).(type) {
	case int:
		return float64(math.MinInt), float64(math.MaxInt)
	case int8:
		return float64(math.MinInt8), float64(math.MaxInt8)
	case int16:
		return float64(math.MinInt16), float64(math.MaxInt16)
	case int32:
		return float64(math.MinInt32), float64(math.MaxInt32)
	case int64:
		return float64(math.MinInt64), float64(math.MaxInt64)
	}

	return 0, 0
}

// getUnsignedMinMax returns the bounds of an unsigned integer type as float64.
// It returns (0, 0) for any other type.
func getUnsignedMinMax[T Number]() (minVal, maxVal float64) {
	var zero T

	switch any(zero).(type) {
	case uint:
		return 0, float64(math.MaxUint)
	case uint8:
		return 0, float64(math.MaxUint8)
	case uint16:
		return 0, float64(math.MaxUint16)
	case uint32:
		return 0, float64(math.MaxUint32)
	case uint64:
		return 0, float64(math.MaxUint64)
	}

	return 0, 0
}

// getFloatMinMax returns the bounds of a floating-point type as float64.
// It returns (0, 0) for any other type.
func getFloatMinMax[T Number]() (minVal, maxVal float64) {
	var zero T

	switch any(zero).(type) {
	case float32:
		return -math.MaxFloat32, math.MaxFloat32
	case float64:
		return -math.MaxFloat64, math.MaxFloat64
	}

	return 0, 0
}

// isInteger checks if the given type parameter is an integer type.
// This is an internal helper used to determine if precision loss checks are needed.
//
// Parameters:
//   - T: The numeric type to check (provided via type parameter).
//
// Returns:
//   - bool: true if T is any signed or unsigned integer type, false otherwise.
func isInteger[T Number]() bool {
	var zero T

	switch any(zero).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}

	return false
}

// isFloat checks if the given type parameter is a floating-point type.
// This is an internal helper used to determine if precision loss checks are needed.
//
// Parameters:
//   - T: The numeric type to check (provided via type parameter).
//
// Returns:
//   - bool: true if T is float32 or float64, false otherwise.
func isFloat[T Number]() bool {
	var zero T

	switch any(zero).(type) {
	case float32, float64:
		return true
	}

	return false
}
