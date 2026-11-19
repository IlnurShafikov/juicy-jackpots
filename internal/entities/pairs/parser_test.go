package pairs

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePairs(t *testing.T) {
	errIs := func(expErr error) func(t *testing.T, gotErr error) {
		return func(t *testing.T, gotErr error) {
			t.Helper()
			require.ErrorIs(t, gotErr, expErr)
		}
	}

	tests := []struct {
		name     string
		data     [][]string
		exp      Pairs
		checkErr func(t *testing.T, gotErr error)
	}{
		{
			name:     "not enough value",
			data:     [][]string{{"1"}},
			exp:      nil,
			checkErr: errIs(ErrInvalidData),
		},
		{
			name: "invalid key",
			data: [][]string{{"abc", "100"}},
			exp:  nil,
			checkErr: func(t *testing.T, gotErr error) {
				t.Helper()
				require.Error(t, gotErr)
				require.Contains(t, gotErr.Error(), "invalid ID")

				var expErr *strconv.NumError
				require.ErrorAs(t, gotErr, &expErr, "invalid ID")
			},
		},
		{
			name: "invalid value",
			data: [][]string{{"1", "xyz"}},
			exp:  nil,
			checkErr: func(t *testing.T, gotErr error) {
				t.Helper()
				require.Error(t, gotErr)
				require.Contains(t, gotErr.Error(), "invalid value")

				var expErr *strconv.NumError
				require.ErrorAs(t, gotErr, &expErr, "invalid value")
			},
		},
		{
			name: "parse data successful",
			data: [][]string{
				{"1", "100"},
				{"2", "200"},
			},
			exp: Pairs{
				{Key: 1, Value: 100},
				{Key: 2, Value: 200},
			},
			checkErr: errIs(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Parse(tt.data)
			tt.checkErr(t, err)
			require.Equal(t, tt.exp, resp)
		})
	}
}
