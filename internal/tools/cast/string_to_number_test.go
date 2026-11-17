//nolint:all // it's ok for cast package
package cast

import (
	"errors"
	"math"
	"testing"
)

// TestToNumber tests the ToNumber function for various numeric types and input strings.
func TestToNumber(t *testing.T) {
	tests := []struct {
		name      string
		convert   func() (interface{}, error)
		wantValue interface{}
		wantError error
		isFloat   bool
	}{
		{
			name:      "string to int8",
			convert:   func() (interface{}, error) { return ToNumber[int8]("100") },
			wantValue: int8(100),
			wantError: nil,
		},
		{
			name:      "string to int16",
			convert:   func() (interface{}, error) { return ToNumber[int16]("1000") },
			wantValue: int16(1000),
			wantError: nil,
		},
		{
			name:      "string to int32",
			convert:   func() (interface{}, error) { return ToNumber[int32]("100000") },
			wantValue: int32(100000),
			wantError: nil,
		},
		{
			name:      "string to int64",
			convert:   func() (interface{}, error) { return ToNumber[int64]("1000000000") },
			wantValue: int64(1000000000),
			wantError: nil,
		},
		{
			name:      "string to uint8",
			convert:   func() (interface{}, error) { return ToNumber[uint8]("255") },
			wantValue: uint8(255),
			wantError: nil,
		},
		{
			name:      "string to uint16",
			convert:   func() (interface{}, error) { return ToNumber[uint16]("65535") },
			wantValue: uint16(65535),
			wantError: nil,
		},
		{
			name:      "string to uint32",
			convert:   func() (interface{}, error) { return ToNumber[uint32]("4294967295") },
			wantValue: uint32(4294967295),
			wantError: nil,
		},
		{
			name:      "string to uint64",
			convert:   func() (interface{}, error) { return ToNumber[uint64]("9007199254740992") },
			wantValue: uint64(9007199254740992),
			wantError: nil,
		},
		{
			name:      "string to float32",
			convert:   func() (interface{}, error) { return ToNumber[float32]("3.14") },
			wantValue: float32(3.14),
			wantError: nil,
			isFloat:   true,
		},
		{
			name:      "string to float64",
			convert:   func() (interface{}, error) { return ToNumber[float64]("3.141592653589793") },
			wantValue: float64(3.141592653589793),
			wantError: nil,
			isFloat:   true,
		},
		{
			name:      "string with decimal to float64",
			convert:   func() (interface{}, error) { return ToNumber[float64]("123.456") },
			wantValue: float64(123.456),
			wantError: nil,
			isFloat:   true,
		},
		{
			name:      "negative string to int8",
			convert:   func() (interface{}, error) { return ToNumber[int8]("-42") },
			wantValue: int8(-42),
			wantError: nil,
		},
		{
			name:      "negative string to int32",
			convert:   func() (interface{}, error) { return ToNumber[int32]("-1000000") },
			wantValue: int32(-1000000),
			wantError: nil,
		},
		{
			name:      "negative string to float64",
			convert:   func() (interface{}, error) { return ToNumber[float64]("-123.456") },
			wantValue: float64(-123.456),
			wantError: nil,
			isFloat:   true,
		},
		// Нулевые значения
		{
			name:      "zero string to int",
			convert:   func() (interface{}, error) { return ToNumber[int]("0") },
			wantValue: int(0),
			wantError: nil,
		},
		{
			name:      "zero string to uint",
			convert:   func() (interface{}, error) { return ToNumber[uint]("0") },
			wantValue: uint(0),
			wantError: nil,
		},
		{
			name:      "zero string to float64",
			convert:   func() (interface{}, error) { return ToNumber[float64]("0.0") },
			wantValue: float64(0.0),
			wantError: nil,
		},
		{
			name:      "int8 max",
			convert:   func() (interface{}, error) { return ToNumber[int8]("127") },
			wantValue: int8(127),
			wantError: nil,
		},
		{
			name:      "int8 min",
			convert:   func() (interface{}, error) { return ToNumber[int8]("-128") },
			wantValue: int8(-128),
			wantError: nil,
		},
		{
			name:      "uint8 max",
			convert:   func() (interface{}, error) { return ToNumber[uint8]("255") },
			wantValue: uint8(255),
			wantError: nil,
		},
		{
			name:      "uint8 min",
			convert:   func() (interface{}, error) { return ToNumber[uint8]("0") },
			wantValue: uint8(0),
			wantError: nil,
		},
		{
			name:      "int16 max",
			convert:   func() (interface{}, error) { return ToNumber[int16]("32767") },
			wantValue: int16(32767),
			wantError: nil,
		},
		{
			name:      "int16 min",
			convert:   func() (interface{}, error) { return ToNumber[int16]("-32768") },
			wantValue: int16(-32768),
			wantError: nil,
		},
		{
			name:      "int32 max",
			convert:   func() (interface{}, error) { return ToNumber[int32]("2147483647") },
			wantValue: int32(2147483647),
			wantError: nil,
		},
		{
			name:      "int32 min",
			convert:   func() (interface{}, error) { return ToNumber[int32]("-2147483648") },
			wantValue: int32(-2147483648),
			wantError: nil,
		},
		{
			name:      "integer with .0 to int",
			convert:   func() (interface{}, error) { return ToNumber[int]("100.0") },
			wantValue: int(100),
			wantError: nil,
		},
		{
			name:      "integer with .00 to int32",
			convert:   func() (interface{}, error) { return ToNumber[int32]("42.00") },
			wantValue: int32(42),
			wantError: nil,
		},
		{
			name:      "exponential notation to float64",
			convert:   func() (interface{}, error) { return ToNumber[float64]("1.5e2") },
			wantValue: float64(150.0),
			wantError: nil,
			isFloat:   true,
		},
		{
			name:      "exponential notation to int",
			convert:   func() (interface{}, error) { return ToNumber[int]("1e2") },
			wantValue: int(100),
			wantError: nil,
		},
		{
			name:      "invalid format - letters",
			convert:   func() (interface{}, error) { return ToNumber[int]("abc") },
			wantValue: int(0),
			wantError: ErrInvalidFormat,
		},
		{
			name:      "invalid format - mixed",
			convert:   func() (interface{}, error) { return ToNumber[int]("12abc") },
			wantValue: int(0),
			wantError: ErrInvalidFormat,
		},
		{
			name:      "invalid format - empty string",
			convert:   func() (interface{}, error) { return ToNumber[int]("") },
			wantValue: int(0),
			wantError: ErrInvalidFormat,
		},
		{
			name:      "invalid format - spaces",
			convert:   func() (interface{}, error) { return ToNumber[int]("  ") },
			wantValue: int(0),
			wantError: ErrInvalidFormat,
		},
		{
			name:      "invalid format - multiple dots",
			convert:   func() (interface{}, error) { return ToNumber[float64]("1.2.3") },
			wantValue: float64(0),
			wantError: ErrInvalidFormat,
		},
		{
			name:      "NaN string to int",
			convert:   func() (interface{}, error) { return ToNumber[int]("NaN") },
			wantValue: int(0),
			wantError: ErrNaN,
		},
		{
			name:      "NaN string to float64",
			convert:   func() (interface{}, error) { return ToNumber[float64]("NaN") },
			wantValue: float64(0),
			wantError: ErrNaN,
		},
		{
			name:      "positive infinity to int",
			convert:   func() (interface{}, error) { return ToNumber[int]("+Inf") },
			wantValue: int(0),
			wantError: ErrInfinity,
		},
		{
			name:      "negative infinity to int",
			convert:   func() (interface{}, error) { return ToNumber[int]("-Inf") },
			wantValue: int(0),
			wantError: ErrInfinity,
		},
		{
			name:      "infinity to float64",
			convert:   func() (interface{}, error) { return ToNumber[float64]("Inf") },
			wantValue: float64(0),
			wantError: ErrInfinity,
		},
		{
			name:      "negative infinity to float32",
			convert:   func() (interface{}, error) { return ToNumber[float32]("-Inf") },
			wantValue: float32(0),
			wantError: ErrInfinity,
		},
		{
			name:      "overflow int8 positive",
			convert:   func() (interface{}, error) { return ToNumber[int8]("128") },
			wantValue: int8(0),
			wantError: ErrOverflow,
		},
		{
			name:      "overflow int8 large",
			convert:   func() (interface{}, error) { return ToNumber[int8]("1000") },
			wantValue: int8(0),
			wantError: ErrOverflow,
		},
		{
			name:      "overflow uint8",
			convert:   func() (interface{}, error) { return ToNumber[uint8]("256") },
			wantValue: uint8(0),
			wantError: ErrOverflow,
		},
		{
			name:      "overflow int16",
			convert:   func() (interface{}, error) { return ToNumber[int16]("32768") },
			wantValue: int16(0),
			wantError: ErrOverflow,
		},
		{
			name:      "overflow uint16",
			convert:   func() (interface{}, error) { return ToNumber[uint16]("65536") },
			wantValue: uint16(0),
			wantError: ErrOverflow,
		},
		{
			name:      "overflow int32",
			convert:   func() (interface{}, error) { return ToNumber[int32]("2147483648") },
			wantValue: int32(0),
			wantError: ErrOverflow,
		},
		{
			name:      "overflow uint32",
			convert:   func() (interface{}, error) { return ToNumber[uint32]("4294967296") },
			wantValue: uint32(0),
			wantError: ErrOverflow,
		},
		{
			name:      "fractional part for int",
			convert:   func() (interface{}, error) { return ToNumber[int]("42.5") },
			wantValue: int(0),
			wantError: ErrOverflow,
		},
		{
			name:      "fractional part for int8",
			convert:   func() (interface{}, error) { return ToNumber[int8]("10.1") },
			wantValue: int8(0),
			wantError: ErrOverflow,
		},
		{
			name:      "fractional part for uint32",
			convert:   func() (interface{}, error) { return ToNumber[uint32]("100.999") },
			wantValue: uint32(0),
			wantError: ErrOverflow,
		},
		{
			name:      "small fractional for int64",
			convert:   func() (interface{}, error) { return ToNumber[int64]("123.0001") },
			wantValue: int64(0),
			wantError: ErrOverflow,
		},
		{
			name:      "underflow int8",
			convert:   func() (interface{}, error) { return ToNumber[int8]("-129") },
			wantValue: int8(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "underflow int8 large negative",
			convert:   func() (interface{}, error) { return ToNumber[int8]("-1000") },
			wantValue: int8(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "underflow uint - negative",
			convert:   func() (interface{}, error) { return ToNumber[uint]("-1") },
			wantValue: uint(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "underflow uint8 negative",
			convert:   func() (interface{}, error) { return ToNumber[uint8]("-1") },
			wantValue: uint8(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "underflow int16",
			convert:   func() (interface{}, error) { return ToNumber[int16]("-32769") },
			wantValue: int16(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "underflow uint16 negative",
			convert:   func() (interface{}, error) { return ToNumber[uint16]("-100") },
			wantValue: uint16(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "underflow int32",
			convert:   func() (interface{}, error) { return ToNumber[int32]("-2147483649") },
			wantValue: int32(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "underflow uint32 negative",
			convert:   func() (interface{}, error) { return ToNumber[uint32]("-1") },
			wantValue: uint32(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "underflow uint64 negative",
			convert:   func() (interface{}, error) { return ToNumber[uint64]("-1000") },
			wantValue: uint64(0),
			wantError: ErrUnderflow,
		},
		{
			name:      "negative float to uint",
			convert:   func() (interface{}, error) { return ToNumber[uint]("-0.1") },
			wantValue: uint(0),
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
					var gotFloat, wantFloat float64
					switch v := got.(type) {
					case float32:
						gotFloat = float64(v)
					case float64:
						gotFloat = v
					}
					switch v := tt.wantValue.(type) {
					case float32:
						wantFloat = float64(v)
					case float64:
						wantFloat = v
					}
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
