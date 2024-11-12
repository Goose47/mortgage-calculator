// Package reposmock provides mock for repos.
package reposmock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mortgage-calculator/src/internal/domain/dto"
	"mortgage-calculator/src/internal/domain/dto/requests"
)

// MockCacheGetSaver represent mock cache repository for calculations.
type MockCacheGetSaver struct {
	mock.Mock
}

// Get mocks getting calc result.
func (m *MockCacheGetSaver) Get(ctx context.Context, in *requests.CalculateRequest) (*dto.CalcAggregates, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*dto.CalcAggregates), args.Error(1) //nolint:wrapcheck,errcheck // already returns wrapped errors
}

// Set mocks setting calc result.
func (m *MockCacheGetSaver) Set(ctx context.Context, in *requests.CalculateRequest, aggregates *dto.CalcAggregates) error {
	args := m.Called(ctx, in, aggregates)
	return args.Error(0) //nolint:wrapcheck // already returns wrapped errors
}

// List mocks listing calc results.
func (m *MockCacheGetSaver) List(ctx context.Context) ([]*dto.CacheEntry, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*dto.CacheEntry), args.Error(1) //nolint:wrapcheck,errcheck // already returns wrapped errors
}
