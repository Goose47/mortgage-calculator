package memory

import (
	"context"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	cachepkg "mortgage-calculator/src/internal/cache"
	"mortgage-calculator/src/internal/lib/random"
	"testing"
	"time"
)

func setup(ttl int64) (*Cache, context.Context) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	c := New(log, ttl)
	ctx := context.Background()

	return c, ctx
}

func TestNew(t *testing.T) {
	c, _ := setup(100)

	if c == nil {
		t.Fatalf("calculator service is nil")
	}
}

func TestCache_Get_NonexistentKey(t *testing.T) {
	c, ctx := setup(100)

	key, _ := random.String(10)
	val, err := c.Get(ctx, key)
	require.Error(t, err)
	require.ErrorIs(t, err, cachepkg.ErrKeyNotExists)
	require.Empty(t, val)
}

func TestCache_Set(t *testing.T) {
	c, ctx := setup(100)

	key, _ := random.String(10)
	val, _ := random.String(10)

	err := c.Set(ctx, key, []byte(val))
	require.Empty(t, err)

	res, err := c.Get(ctx, key)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, string(res), val)
}

func TestCache_Get(t *testing.T) {
	c, ctx := setup(100)

	key, _ := random.String(10)
	val, _ := random.String(10)

	err := c.Set(ctx, key, []byte(val))
	require.NoError(t, err)

	res, err := c.Get(ctx, key)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, string(res), val)
}

func TestCache_List_EmptyCache(t *testing.T) {
	c, ctx := setup(100)

	list, err := c.List(ctx)
	require.NoError(t, err)
	require.Empty(t, list)
}

func TestCache_List(t *testing.T) {
	c, ctx := setup(100)

	key, _ := random.String(10)

	err := c.Set(ctx, key, []byte(key))
	require.Empty(t, err)

	list, err := c.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, list)
	require.Equal(t, len(list), 1)

	require.Equal(t, key, string(list[0].Val))
}

func TestCache_Clear(t *testing.T) {
	c, ctx := setup(0)

	key, _ := random.String(10)

	err := c.Set(ctx, key, []byte(key))
	require.Empty(t, err)

	time.Sleep(1 * time.Second)

	c.Clear(ctx)

	require.Empty(t, c.data)
}
