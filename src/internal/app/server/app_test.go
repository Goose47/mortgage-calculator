package serverapp

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"testing"
)

func TestNew(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	app := New(log, 1000, gin.New())

	require.NotEmpty(t, app)
}
