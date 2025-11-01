package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("map of environments", func(t *testing.T) {
		result, err := ReadDir("./testdata/env/")
		require.NoError(t, err)

		tests := []struct {
			key        string
			wantValue  string
			wantRemove bool
		}{
			{"BAR", "bar", false},
			{"EMPTY", "", false},
			{"FOO", "   foo\nwith new line", false},
			{"HELLO", `"hello"`, false},
			{"UNSET", "", true},
		}

		for _, tt := range tests {
			t.Run(tt.key, func(t *testing.T) {
				val, ok := result[tt.key]
				require.True(t, ok, "missing key %s", tt.key)
				require.Equal(t, tt.wantValue, val.Value)
				require.Equal(t, tt.wantRemove, val.NeedRemove)
			})
		}
	})
}
