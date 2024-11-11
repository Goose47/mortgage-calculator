// Package cacherepos provides repositories to interact with cache.
package cacherepos

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	cachepkg "mortgage-calculator/src/internal/cache"
	"mortgage-calculator/src/internal/controllers"
	"mortgage-calculator/src/internal/domain/dto"
)

// Cache represents cache api.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Clear(ctx context.Context)
	List(ctx context.Context) ([]*cachepkg.Entry, error)
}

// CalcRepository is a repo to save and retrieve calculation results.
type CalcRepository struct {
	log   *slog.Logger
	cache Cache
}

// NewCalcRepository is a constructor for CalcRepository.
func NewCalcRepository(
	log *slog.Logger,
	cache Cache,
) *CalcRepository {
	return &CalcRepository{
		log:   log,
		cache: cache,
	}
}

// Get generates key by input and returns cached result.
func (r *CalcRepository) Get(
	ctx context.Context,
	in *controllers.CalculateRequest,
) (*dto.CalcAggregates, error) {
	const op = "cacherepos.calcRepository.Get"
	log := r.log.With(slog.String("op", op))

	log.Info("generating key")

	key, err := generateKey(in)
	if err != nil {
		log.Error("failed to generate key", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log = log.With(slog.String("key", key))
	log.Info("key generated, trying to retrieve key from cache")

	byteArr, err := r.cache.Get(ctx, key)
	if err != nil {
		log.Info("failed to retrieve key from cache", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("cache hit, unmarshalling result")

	var res dto.CalcAggregates
	err = json.Unmarshal(byteArr, &res)
	if err != nil {
		log.Info("failed to unmarshal result", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &res, nil
}

// Set generates key by input, marshals result and caches it.
func (r *CalcRepository) Set(
	ctx context.Context,
	in *controllers.CalculateRequest,
	aggregates *dto.CalcAggregates,
) error {
	const op = "cacherepos.calcRepository.Set"
	log := r.log.With(slog.String("op", op))

	log.Info("generating key")

	key, err := generateKey(in)
	if err != nil {
		log.Error("failed to generate key", slog.Any("error", err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log = log.With(slog.String("key", key))

	log.Info("key generated, marshalling data")

	byteArr, err := json.Marshal(aggregates)
	if err != nil {
		log.Error("failed to marshal data", slog.Any("error", err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("data marshaled, storing data in cache")

	err = r.cache.Set(ctx, key, byteArr)
	if err != nil {
		log.Error("failed to save data", slog.Any("error", err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Clear cleans expired items from cache.
func (r *CalcRepository) Clear(ctx context.Context) {
	const op = "cacherepos.calcRepository.Clear"
	log := r.log.With(slog.String("op", op))

	log.Info("clearing expired cache entries")
	r.cache.Clear(ctx)
	log.Info("deleted expired items from cache")
}

// List lists all active items.
func (r *CalcRepository) List(ctx context.Context) ([]*dto.CacheEntry, error) {
	const op = "cacherepos.calcRepository.List"
	log := r.log.With(slog.String("op", op))

	log.Info("retrieving all cache items")

	items, err := r.cache.List(ctx)
	if err != nil {
		log.Info("failed to retrieve cache items")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("cache items retrieved, unmarshalling")

	res := make([]*dto.CacheEntry, len(items))

	for i, item := range items {
		var aggregates dto.CalcAggregates
		err := json.Unmarshal(item.Val, &aggregates)
		if err != nil {
			log.Info("failed to unmarshal aggregates")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		var params controllers.CalculateRequest
		err = json.Unmarshal([]byte(item.Key), &params)
		if err != nil {
			log.Info("failed to unmarshal params")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		res[i] = &dto.CacheEntry{
			ID:         item.ID,
			Aggregates: &aggregates,
			Params: &dto.CalcParams{
				ObjectCost:     params.ObjectCost,
				InitialPayment: params.InitialPayment,
				Months:         params.Months,
			},
			Program: &params.Program,
		}
	}

	return res, nil
}

func generateKey(in *controllers.CalculateRequest) (string, error) {
	byteArr, err := json.Marshal(in)

	if err != nil {
		return "", fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return string(byteArr), nil
}
