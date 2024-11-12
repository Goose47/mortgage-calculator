// Package servicesmock provides mock for services.
package servicesmock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mortgage-calculator/src/internal/domain/dto"
)

// MockCalculator mocks service layer for calculations.
type MockCalculator struct {
	mock.Mock
}

// Calculate mocks calculations.
func (m *MockCalculator) Calculate(ctx context.Context, params dto.CalcParams, program dto.CalcProgram) (*dto.CalcAggregates, error) {
	args := m.Called(ctx, params, program)
	return args.Get(0).(*dto.CalcAggregates), args.Error(1) //nolint:wrapcheck,errcheck // already returns wrapped errors
}
