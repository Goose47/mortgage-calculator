// Package cachemock provides mock for cache.
package cachemock

import (
	"context"
	"github.com/stretchr/testify/mock"
	cachepkg "mortgage-calculator/src/internal/cache"
)

// MockCache is a mock realization of cache.
type MockCache struct {
	mock.Mock
}

// Get returns value by key.
func (m *MockCache) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	return args.Get(0).([]byte), args.Error(1) //nolint:wrapcheck,errcheck // already returns wrapped errors
}

// Set saves given value by given key.
func (m *MockCache) Set(ctx context.Context, key string, value []byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0) //nolint:wrapcheck // already returns wrapped errors
}

// Clear checks whether entries are expired and deletes them if true.
func (m *MockCache) Clear(_ context.Context) {}

// List returns all active cache entries.
func (m *MockCache) List(ctx context.Context) ([]*cachepkg.Entry, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*cachepkg.Entry), args.Error(1) //nolint:wrapcheck,errcheck // already returns wrapped errors
}
