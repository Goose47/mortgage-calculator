package app

import (
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"testing"
)

func TestNew(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	app := New(log, "dev", 1000, 1000)

	require.NotEmpty(t, app)
}
