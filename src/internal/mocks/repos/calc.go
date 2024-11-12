package reposmock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mortgage-calculator/src/internal/domain/dto"
	"mortgage-calculator/src/internal/domain/dto/requests"
)

type MockCacheGetSaver struct {
	mock.Mock
}

func (m *MockCacheGetSaver) Get(ctx context.Context, in *requests.CalculateRequest) (*dto.CalcAggregates, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*dto.CalcAggregates), args.Error(1)
}

func (m *MockCacheGetSaver) Set(ctx context.Context, in *requests.CalculateRequest, aggregates *dto.CalcAggregates) error {
	args := m.Called(ctx, in, aggregates)
	return args.Error(0)
}
