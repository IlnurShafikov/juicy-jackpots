//nolint:all // it's ok for cast package
package cast

import (
	"errors"
	"math"
	"testing"
)

func TestToString(t *testing.T) {
	tests := []struct {
		name      string
		convert   func() (string, error)
		wantValue string
		wantError error
	}{
		{
			name:      "int positive",
			convert:   func() (string, error) { return ToString(int(42)) },
			wantValue: "42",
			wantError: nil,
		},
		{
			name:      "int negative",
			convert:   func() (string, error) { return ToString(int(-42)) },
			wantValue: "-42",
			wantError: nil,
		},
		{
			name:      "int zero",
			convert:   func() (string, error) { return ToString(int(0)) },
			wantValue: "0",
			wantError: nil,
		},
		{
			name:      "int max",
			convert:   func() (string, error) { return ToString(int(math.MaxInt64)) },
			wantValue: "9223372036854775807",
			wantError: nil,
		},
		{
			name:      "int min",
			convert:   func() (string, error) { return ToString(int(math.MinInt64)) },
			wantValue: "-9223372036854775808",
			wantError: nil,
		},

		// int8
		{
			name:      "int8 positive",
			convert:   func() (string, error) { return ToString(int8(100)) },
			wantValue: "100",
			wantError: nil,
		},
		{
			name:      "int8 max",
			convert:   func() (string, error) { return ToString(int8(math.MaxInt8)) },
			wantValue: "127",
			wantError: nil,
		},
		{
			name:      "int8 min",
			convert:   func() (string, error) { return ToString(int8(math.MinInt8)) },
			wantValue: "-128",
			wantError: nil,
		},
		{
			name:      "int16 positive",
			convert:   func() (string, error) { return ToString(int16(1000)) },
			wantValue: "1000",
			wantError: nil,
		},
		{
			name:      "int16 max",
			convert:   func() (string, error) { return ToString(int16(math.MaxInt16)) },
			wantValue: "32767",
			wantError: nil,
		},
		{
			name:      "int16 min",
			convert:   func() (string, error) { return ToString(int16(math.MinInt16)) },
			wantValue: "-32768",
			wantError: nil,
		},
		{
			name:      "int32 positive",
			convert:   func() (string, error) { return ToString(int32(100000)) },
			wantValue: "100000",
			wantError: nil,
		},
		{
			name:      "int32 max",
			convert:   func() (string, error) { return ToString(int32(math.MaxInt32)) },
			wantValue: "2147483647",
			wantError: nil,
		},
		{
			name:      "int32 min",
			convert:   func() (string, error) { return ToString(int32(math.MinInt32)) },
			wantValue: "-2147483648",
			wantError: nil,
		},
		{
			name:      "int64 positive",
			convert:   func() (string, error) { return ToString(int64(1000000000)) },
			wantValue: "1000000000",
			wantError: nil,
		},
		{
			name:      "int64 max",
			convert:   func() (string, error) { return ToString(int64(math.MaxInt64)) },
			wantValue: "9223372036854775807",
			wantError: nil,
		},
		{
			name:      "int64 min",
			convert:   func() (string, error) { return ToString(int64(math.MinInt64)) },
			wantValue: "-9223372036854775808",
			wantError: nil,
		},
		{
			name:      "uint positive",
			convert:   func() (string, error) { return ToString(uint(42)) },
			wantValue: "42",
			wantError: nil,
		},
		{
			name:      "uint zero",
			convert:   func() (string, error) { return ToString(uint(0)) },
			wantValue: "0",
			wantError: nil,
		},
		{
			name:      "uint8 max",
			convert:   func() (string, error) { return ToString(uint8(math.MaxUint8)) },
			wantValue: "255",
			wantError: nil,
		},
		{
			name:      "uint16 max",
			convert:   func() (string, error) { return ToString(uint16(math.MaxUint16)) },
			wantValue: "65535",
			wantError: nil,
		},
		{
			name:      "uint32 max",
			convert:   func() (string, error) { return ToString(uint32(math.MaxUint32)) },
			wantValue: "4294967295",
			wantError: nil,
		},
		{
			name:      "uint64 positive",
			convert:   func() (string, error) { return ToString(uint64(1000000000)) },
			wantValue: "1000000000",
			wantError: nil,
		},

		// float32
		{
			name:      "float32 simple",
			convert:   func() (string, error) { return ToString(float32(42.5)) },
			wantValue: "42.50000000",
			wantError: nil,
		},
		{
			name:      "float32 zero",
			convert:   func() (string, error) { return ToString(float32(0.0)) },
			wantValue: "0.00000000",
			wantError: nil,
		},
		{
			name:      "float32 negative",
			convert:   func() (string, error) { return ToString(float32(-42.5)) },
			wantValue: "-42.50000000",
			wantError: nil,
		},
		{
			name:      "float32 precision",
			convert:   func() (string, error) { return ToString(float32(3.14159265)) },
			wantValue: "3.14159274",
			wantError: nil,
		},
		{
			name:      "float64 simple",
			convert:   func() (string, error) { return ToString(float64(42.5)) },
			wantValue: "42.50000000",
			wantError: nil,
		},
		{
			name:      "float64 zero",
			convert:   func() (string, error) { return ToString(float64(0.0)) },
			wantValue: "0.00000000",
			wantError: nil,
		},
		{
			name:      "float64 negative",
			convert:   func() (string, error) { return ToString(float64(-42.5)) },
			wantValue: "-42.50000000",
			wantError: nil,
		},
		{
			name:      "float64 small",
			convert:   func() (string, error) { return ToString(float64(0.00000001)) },
			wantValue: "0.00000001",
			wantError: nil,
		},
		{
			name:      "float64 very small",
			convert:   func() (string, error) { return ToString(float64(0.000000001)) },
			wantValue: "0.00000000",
			wantError: nil,
		},
		{
			name:      "float64 large",
			convert:   func() (string, error) { return ToString(float64(123456789.123456789)) },
			wantValue: "123456789.12345679",
			wantError: nil,
		},
		{
			name:      "float64 precision",
			convert:   func() (string, error) { return ToString(float64(3.14159265)) },
			wantValue: "3.14159265",
			wantError: nil,
		},
		{
			name:      "float64 NaN",
			convert:   func() (string, error) { return ToString(math.NaN()) },
			wantValue: "",
			wantError: ErrNaN,
		},
		{
			name:      "float32 NaN",
			convert:   func() (string, error) { return ToString(float32(math.NaN())) },
			wantValue: "",
			wantError: ErrNaN,
		},
		{
			name:      "float64 positive infinity",
			convert:   func() (string, error) { return ToString(math.Inf(1)) },
			wantValue: "",
			wantError: ErrInfinity,
		},
		{
			name:      "float64 negative infinity",
			convert:   func() (string, error) { return ToString(math.Inf(-1)) },
			wantValue: "",
			wantError: ErrInfinity,
		},
		{
			name:      "float32 positive infinity",
			convert:   func() (string, error) { return ToString(float32(math.Inf(1))) },
			wantValue: "",
			wantError: ErrInfinity,
		},
		{
			name:      "float32 negative infinity",
			convert:   func() (string, error) { return ToString(float32(math.Inf(-1))) },
			wantValue: "",
			wantError: ErrInfinity,
		},
		{
			name:      "float64 max",
			convert:   func() (string, error) { return ToString(math.MaxFloat64) },
			wantValue: "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.00000000",
			wantError: nil,
		},
		{
			name:      "float64 smallest nonzero",
			convert:   func() (string, error) { return ToString(math.SmallestNonzeroFloat64) },
			wantValue: "0.00000000",
			wantError: nil,
		},
		{
			name:      "float32 max",
			convert:   func() (string, error) { return ToString(float32(math.MaxFloat32)) },
			wantValue: "340282346638528859811704183484516925440.00000000",
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.convert()

			if !errors.Is(err, tt.wantError) {
				t.Errorf("error = %v, wantError %v", err, tt.wantError)
				return
			}

			if got != tt.wantValue {
				t.Errorf("value = %q, wantValue %q", got, tt.wantValue)
			}
		})
	}
}
