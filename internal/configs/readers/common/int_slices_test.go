package common

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDataToSymbols(t *testing.T) {
	errIs := func(expErr error) func(t *testing.T, gotErr error) {
		return func(t *testing.T, gotErr error) {
			t.Helper()
			require.ErrorIs(t, gotErr, expErr)
		}
	}

	tests := []struct {
		name     string
		data     [][]string
		exp      []int
		checkErr func(t *testing.T, gotErr error)
	}{
		{
			name:     "got err empty data",
			data:     [][]string{},
			exp:      nil,
			checkErr: errIs(ErrEmptyData),
		},
		{
			name: "invalid value",
			data: [][]string{{"1"}, {"2"}, {"d"}, {"4"}},
			exp:  nil,
			checkErr: func(t *testing.T, gotErr error) {
				t.Helper()
				require.Error(t, gotErr)

				var expErr *strconv.NumError
				require.ErrorAs(t, gotErr, &expErr, "strconv.Atoi")
			},
		},
		{
			name:     "parse data successful",
			data:     [][]string{{"1"}, {"2"}, {"3"}, {"4"}},
			exp:      []int{1, 2, 3, 4},
			checkErr: errIs(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := parseDataToSliceInt(tt.data)
			tt.checkErr(t, err)
			require.Equal(t, tt.exp, resp)
		})
	}
}
