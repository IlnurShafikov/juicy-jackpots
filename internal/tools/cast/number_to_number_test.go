//nolint:all // it's ok for cast package
package cast

import (
	"errors"
	"math"
	"testing"
)

func TestSafeCast(t *testing.T) {
	tests := []struct {
		name      string
		convert   func() (interface{}, error)
		wantValue interface{}
		wantError error
		isFloat   bool
	}{
		{
			name:      "int8 to int16",
			convert:   func() (interface{}, error) { return SafeCast[int16](int8(100)) },
			wantValue: int16(100),
			wantError: nil,
		},
		{
			name:      "int32 to int64",
			convert:   func() (interface{}, error) { return SafeCast[int64](int32(1000000)) },
			wantValue: int64(1000000),
			wantError: nil,
		},
		{
			name:      "uint8 to uint32",
			convert:   func() (interface{}, error) { return SafeCast[uint32](uint8(255)) },
			wantValue: uint32(255),
			wantError: nil,
		},
		{
			name:      "float32 to float64",
			convert:   func() (interface{}, error) { return SafeCast[float64](float32(3.14)) },
			wantValue: float64(3.14),
			wantError: nil,
			isFloat:   true,
		},
		{
			name:      "float64 to int no fraction",
			convert:   func() (interface{}, error) { return SafeCast[int](float64(100.0)) },
			wantValue: 100,
			wantError: nil,
		},
		{
			name:      "negative int to int64",
			convert:   func() (interface{}, error) { return SafeCast[int64](int(-42)) },
			wantValue: int64(-42),
			wantError: nil,
		},
		{
			name:      "zero int to uint",
			convert:   func() (interface{}, error) { return SafeCast[uint](0) },
			wantValue: uint(0),
			wantError: nil,
		},
		{
			name:      "int8 max",
			convert:   func() (interface{}, error) { return SafeCast[int8](int8(math.MaxInt8)) },
			wantValue: int8(127),
			wantError: nil,
		},
		{
			name:      "int8 min",
			convert:   func() (interface{}, error) { return SafeCast[int8](int8(math.MinInt8)) },
			wantValue: int8(-128),
			wantError: nil,
		},
		{
			name:      "uint8 max",
			convert:   func() (interface{}, error) { return SafeCast[uint8](uint8(math.MaxUint8)) },
			wantValue: uint8(255),
			wantError: nil,
		},
		{
			name:      "uint8 min",
			convert:   func() (interface{}, error) { return SafeCast[uint8](uint8(0)) },
			wantValue: uint8(0),
			wantError: nil,
		},
		{
			name:      "int16 127 to int8",
			convert:   func() (interface{}, error) { return SafeCast[int8](int16(127)) },
			wantValue: int8(127),
			wantError: nil,
		},
		{
			name:      "int16 -128 to int8",
			convert:   func() (interface{}, error) { return SafeCast[int8](int16(-128)) },
			wantValue: int8(-128),
			wantError: nil,
		},
		{
			name:      "int32 max",
			convert:   func() (interface{}, error) { return SafeCast[int32](int32(math.MaxInt32)) },
			wantValue: int32(math.MaxInt32),
			wantError: nil,
		},
		{
			name:      "int32 min",
			convert:   func() (interface{}, error) { return SafeCast[int32](int32(math.MinInt32)) },
			wantValue: int32(math.MinInt32),
			wantError: nil,
		},
		{
			name:      "uint32 max",
			convert:   func() (interface{}, error) { return SafeCast[uint32](uint32(math.MaxUint32)) },
			wantValue: uint32(math.MaxUint32),
			wantError: nil,
		},
		{
			name:      "float64 NaN to int",
			convert:   func() (interface{}, error) { return SafeCast[int](math.NaN()) },
			wantValue: 0,
			wantError: ErrNaN,
		},
		{
			name:      "float32 NaN to int32",
			convert:   func() (interface{}, error) { return SafeCast[int32](float32(math.NaN())) },
			wantValue: int32(0),
			wantError: ErrNaN,
		},
		{
			name:      "float64 NaN to float32",
			convert:   func() (interface{}, error) { return SafeCast[float32](math.NaN()) },
			wantValue: float32(0),
			wantError: ErrNaN,
		},
		{
			name:      "positive infinity to int",
			convert:   func() (interface{}, error) { return SafeCast[int](math.Inf(1)) },
			wantValue: 0,
			wantError: ErrInfinity,
		},
		{
			name:      "negative infinity to int",
			convert:   func() (interface{}, error) { return SafeCast[int](math.Inf(-1)) },
			wantValue: 0,
			wantError: ErrInfinity,
		},
		{
			name:      "positive infinity to uint",
			convert:   func() (interface{}, error) { return SafeCast[uint](math.Inf(1)) },
			wantValue: uint(0),
			wantError: ErrInfinity,
		},
		{
			name:      "negative infinity to float32",
			convert:   func() (interface{}, error) { return SafeCast[float32](math.Inf(-1)) },
			wantValue: float32(0),
			wantError: ErrInfinity,
		},
		{
			name:      "int64 max to int8",
			convert:   func() (interface{}, error) { return SafeCast[int8](int64(math.MaxInt64)) },
			wantValue: int8(0),
			wantError: ErrOverflow,
		},
		{
			name:      "int32 128 to int8",
			convert:   func() (interface{}, error) { return SafeCast[int8](int32(128)) },
			wantValue: int8(0),
			wantError: ErrOverflow,
		},
		{
			name:      "uint32 256 to uint8",
			convert:   func() (interface{}, error) { return SafeCast[uint8](uint32(256)) },
			wantValue: uint8(0),
			wantError: ErrOverflow,
		},
		{
			name:      "float64 with fraction to int",
			convert:   func() (interface{}, error) { return SafeCast[int](float64(42.5)) },
			wantValue: 0,
			wantError: ErrOverflow,
		},
		{
			name:      "float32 with fraction to int64",
			convert:   func() (interface{}, error) { return SafeCast[int64](float32(123.456)) },
			wantValue: int64(0),
			wantError: ErrOverflow,
		},
		{
			name:      "large float to int8",
			convert:   func() (interface{}, error) { return SafeCast[int8](float64(1000.0)) },
			wantValue: int8(0),
			wantError: ErrOverflow,
		},
		{
			name:      "int64 max to int32",
			convert:   func() (interface{}, error) { return SafeCast[int32](int64(math.MaxInt64)) },
			wantValue: int32(0),
			wantError: ErrOverflow,
		},
		{
			name:      "uint64 max to uint32",
			convert:   func() (interface{}, error) { return SafeCast[uint32](uint64(math.MaxUint64)) },
			wantValue: uint32(0),
			wantError: ErrOverflow,
		},
		{
			name:      "int64 min to int8",
			convert:   func() (interface{}, error) { return SafeCast[int8](int64(math.MinInt64)) },
			wantValue: int8(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "negative int32 to uint",
			convert:   func() (interface{}, error) { return SafeCast[uint](int32(-1)) },
			wantValue: uint(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "int32 -129 to int8",
			convert:   func() (interface{}, error) { return SafeCast[int8](int32(-129)) },
			wantValue: int8(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "negative float to uint",
			convert:   func() (interface{}, error) { return SafeCast[uint](float64(-1.0)) },
			wantValue: uint(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "int64 min to uint64",
			convert:   func() (interface{}, error) { return SafeCast[uint64](int64(math.MinInt64)) },
			wantValue: uint64(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "int32 -32769 to int16",
			convert:   func() (interface{}, error) { return SafeCast[int16](int32(-32769)) },
			wantValue: int16(0),
			wantError: ErrUnderflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.convert()

			if !errors.Is(err, tt.wantError) {
				t.Errorf("error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError == nil {
				if tt.isFloat {
					// Для float используем приблизительное сравнение
					gotFloat, _ := got.(float64)
					wantFloat, _ := tt.wantValue.(float64)
					if math.Abs(gotFloat-wantFloat) > 0.01 {
						t.Errorf("value = %v, wantValue %v", got, tt.wantValue)
					}
				} else if got != tt.wantValue {
					t.Errorf("value = %v, wantValue %v", got, tt.wantValue)
				}
			}
		})
	}
}

func TestGetMinMax(t *testing.T) {
	tests := []struct {
		name    string
		getFunc func() (float64, float64)
		wantMin float64
		wantMax float64
	}{
		{
			name:    "int8",
			getFunc: func() (float64, float64) { return getMinMax[int8]() },
			wantMin: float64(math.MinInt8),
			wantMax: float64(math.MaxInt8),
		},
		{
			name:    "int16",
			getFunc: func() (float64, float64) { return getMinMax[int16]() },
			wantMin: float64(math.MinInt16),
			wantMax: float64(math.MaxInt16),
		},
		{
			name:    "int32",
			getFunc: func() (float64, float64) { return getMinMax[int32]() },
			wantMin: float64(math.MinInt32),
			wantMax: float64(math.MaxInt32),
		},
		{
			name:    "int64",
			getFunc: func() (float64, float64) { return getMinMax[int64]() },
			wantMin: float64(math.MinInt64),
			wantMax: float64(math.MaxInt64),
		},
		{
			name:    "int",
			getFunc: func() (float64, float64) { return getMinMax[int]() },
			wantMin: float64(math.MinInt),
			wantMax: float64(math.MaxInt),
		},
		{
			name:    "uint8",
			getFunc: func() (float64, float64) { return getMinMax[uint8]() },
			wantMin: 0,
			wantMax: float64(math.MaxUint8),
		},
		{
			name:    "uint16",
			getFunc: func() (float64, float64) { return getMinMax[uint16]() },
			wantMin: 0,
			wantMax: float64(math.MaxUint16),
		},
		{
			name:    "uint32",
			getFunc: func() (float64, float64) { return getMinMax[uint32]() },
			wantMin: 0,
			wantMax: float64(math.MaxUint32),
		},
		{
			name:    "uint64",
			getFunc: func() (float64, float64) { return getMinMax[uint64]() },
			wantMin: 0,
			wantMax: float64(math.MaxUint64),
		},
		{
			name:    "uint",
			getFunc: func() (float64, float64) { return getMinMax[uint]() },
			wantMin: 0,
			wantMax: float64(math.MaxUint),
		},
		{
			name:    "float32",
			getFunc: func() (float64, float64) { return getMinMax[float32]() },
			wantMin: -math.MaxFloat32,
			wantMax: math.MaxFloat32,
		},
		{
			name:    "float64",
			getFunc: func() (float64, float64) { return getMinMax[float64]() },
			wantMin: -math.MaxFloat64,
			wantMax: math.MaxFloat64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minVal, maxVal := tt.getFunc()

			if minVal != tt.wantMin {
				t.Errorf("minVal = %v, want %v", minVal, tt.wantMin)
			}

			if maxVal != tt.wantMax {
				t.Errorf("maxVal = %v, want %v", maxVal, tt.wantMax)
			}
		})
	}
}

func TestIsInteger(t *testing.T) {
	tests := []struct {
		name string
		fn   func() bool
		want bool
	}{
		{
			name: "int",
			fn:   func() bool { return isInteger[int]() },
			want: true,
		},
		{
			name: "int8",
			fn:   func() bool { return isInteger[int8]() },
			want: true,
		},
		{
			name: "int16",
			fn:   func() bool { return isInteger[int16]() },
			want: true,
		},
		{
			name: "int32",
			fn:   func() bool { return isInteger[int32]() },
			want: true,
		},
		{
			name: "int64",
			fn:   func() bool { return isInteger[int64]() },
			want: true,
		},
		{
			name: "uint",
			fn:   func() bool { return isInteger[uint]() },
			want: true,
		},
		{
			name: "uint8",
			fn:   func() bool { return isInteger[uint8]() },
			want: true,
		},
		{
			name: "uint16",
			fn:   func() bool { return isInteger[uint16]() },
			want: true,
		},
		{
			name: "uint32",
			fn:   func() bool { return isInteger[uint32]() },
			want: true,
		},
		{
			name: "uint64",
			fn:   func() bool { return isInteger[uint64]() },
			want: true,
		},
		{
			name: "float32",
			fn:   func() bool { return isInteger[float32]() },
			want: false,
		},
		{
			name: "float64",
			fn:   func() bool { return isInteger[float64]() },
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(); got != tt.want {
				t.Errorf("isInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFloat(t *testing.T) {
	tests := []struct {
		name string
		fn   func() bool
		want bool
	}{
		{
			name: "float32",
			fn:   func() bool { return isFloat[float32]() },
			want: true,
		},
		{
			name: "float64",
			fn:   func() bool { return isFloat[float64]() },
			want: true,
		},
		{
			name: "int",
			fn:   func() bool { return isFloat[int]() },
			want: false,
		},
		{
			name: "int8",
			fn:   func() bool { return isFloat[int8]() },
			want: false,
		},
		{
			name: "int16",
			fn:   func() bool { return isFloat[int16]() },
			want: false,
		},
		{
			name: "int32",
			fn:   func() bool { return isFloat[int32]() },
			want: false,
		},
		{
			name: "int64",
			fn:   func() bool { return isFloat[int64]() },
			want: false,
		},
		{
			name: "uint",
			fn:   func() bool { return isFloat[uint]() },
			want: false,
		},
		{
			name: "uint8",
			fn:   func() bool { return isFloat[uint8]() },
			want: false,
		},
		{
			name: "uint16",
			fn:   func() bool { return isFloat[uint16]() },
			want: false,
		},
		{
			name: "uint32",
			fn:   func() bool { return isFloat[uint32]() },
			want: false,
		},
		{
			name: "uint64",
			fn:   func() bool { return isFloat[uint64]() },
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(); got != tt.want {
				t.Errorf("isFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
