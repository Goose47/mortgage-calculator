package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"mortgage-calculator/src/internal/domain/dto"
	reposmock "mortgage-calculator/src/internal/mocks/repos"
	"testing"
)

func setup() (*CacheService, *reposmock.MockCacheGetSaver) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	cache := new(reposmock.MockCacheGetSaver)
	s := NewCacheService(log, cache)
	return s, cache
}

func TestNewCacheService(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	cache := new(reposmock.MockCacheGetSaver)
	s := NewCacheService(log, cache)

	require.NotEmpty(t, s)
}

func TestCacheService_List(t *testing.T) {
	service, c := setup()

	entries := []*dto.CacheEntry{
		{},
	}

	c.On("List", mock.Anything).Return(entries, nil)

	res, err := service.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, entries, res)
}

func TestCacheService_List_Error(t *testing.T) {
	service, c := setup()

	var entries []*dto.CacheEntry

	c.On("List", mock.Anything).Return(entries, errors.New("failed to retrieve cache entries"))

	res, err := service.List(context.Background())
	require.Error(t, err)
	require.Empty(t, res)
	require.Contains(t, err.Error(), "failed to retrieve cache entries")
}
