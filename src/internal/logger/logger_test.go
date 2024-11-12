package logger

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		in string
	}{
		{"local"},
		{"dev"},
		{"prod"},
	}

	for _, tt := range cases {
		res := New(tt.in)
		require.NotEmpty(t, res)
	}
}
