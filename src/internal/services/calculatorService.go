// Package services provides application business logic.
package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mortgage-calculator/src/internal/domain/dto"
)

// CalculatorService provides api for calculating aggregates
type CalculatorService struct {
	log *slog.Logger
}

// NewCalculatorService is a constructor for CalculatorService
func NewCalculatorService(log *slog.Logger) *CalculatorService {
	return &CalculatorService{
		log: log,
	}
}

// ErrInsufficientInitialPayment represents error when the initial payment to object cost ratio is too small
var ErrInsufficientInitialPayment = errors.New("the initial payment should be more")

const minInitialPaymentRatio = 0.2

// Calculate calculates aggregates based on params and program
func (s *CalculatorService) Calculate(
	ctx context.Context,
	params dto.CalcParams,
	program dto.CalcProgram,
) (*dto.CalcAggregates, error) {
	const op = "calculatorService.Calculate"

	log := s.log.With(slog.String("op", op))

	if float64(params.InitialPayment)/float64(params.ObjectCost) < minInitialPaymentRatio {
		log.Warn(
			"insufficient initial payment",
			slog.Int("initial_payment", params.InitialPayment),
			slog.Int("object_cost", params.ObjectCost),
		)

		return nil, fmt.Errorf("%s: %w", op, ErrInsufficientInitialPayment)
	}

	log.Info("calculating aggregates")

	log.Info("aggregates calculated")

	return &dto.CalcAggregates{}, nil
}
