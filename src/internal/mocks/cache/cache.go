package cachemock

import (
	"context"
	"github.com/stretchr/testify/mock"
	cachepkg "mortgage-calculator/src/internal/cache"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCache) Set(ctx context.Context, key string, value []byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockCache) Clear(ctx context.Context) {}

func (m *MockCache) List(ctx context.Context) ([]*cachepkg.Entry, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*cachepkg.Entry), args.Error(1)
}
