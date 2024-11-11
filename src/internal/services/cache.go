package services

import (
	"context"
	"fmt"
	"log/slog"
	"mortgage-calculator/src/internal/domain/dto"
)

type cache interface {
	List(ctx context.Context) ([]*dto.CacheEntry, error)
}

// CacheService provides api for cache entries analysis.
type CacheService struct {
	log   *slog.Logger
	cache cache
}

// NewCacheService is a constructor for CacheService.
func NewCacheService(
	log *slog.Logger,
	cache cache,
) *CacheService {
	return &CacheService{
		log:   log,
		cache: cache,
	}
}

// List lists all active cache entries.
func (s *CacheService) List(
	ctx context.Context,
) ([]*dto.CacheEntry, error) {
	const op = "cacheService.List"
	log := s.log.With(slog.String("op", op))

	log.Info("retrieving cache entries")

	entries, err := s.cache.List(ctx)
	if err != nil {
		log.Error("failed to retrieve cache entries")

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("cache entries retrieved")

	return entries, nil
}
