package servicesmock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mortgage-calculator/src/internal/domain/dto"
)

type MockCalculator struct {
	mock.Mock
}

func (m *MockCalculator) Calculate(ctx context.Context, params dto.CalcParams, program dto.CalcProgram) (*dto.CalcAggregates, error) {
	args := m.Called(ctx, params, program)
	return args.Get(0).(*dto.CalcAggregates), args.Error(1)
}
