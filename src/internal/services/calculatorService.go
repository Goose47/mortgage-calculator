// Package services provides application business logic.
package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"mortgage-calculator/src/internal/domain/dto"
	"time"
)

// CalculatorService provides api for calculating aggregates.
type CalculatorService struct {
	log *slog.Logger
}

// NewCalculatorService is a constructor for CalculatorService.
func NewCalculatorService(log *slog.Logger) *CalculatorService {
	return &CalculatorService{
		log: log,
	}
}

// ErrInsufficientInitialPayment represents error when the initial payment to object cost ratio is too small.
var ErrInsufficientInitialPayment = errors.New("the initial payment should be more")

const minInitialPaymentRatio = 0.2

// Calculate calculates aggregates based on params and program.
func (s *CalculatorService) Calculate(
	_ context.Context,
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

	months := float64(params.Months)
	objectCost := float64(params.ObjectCost)
	initialPayment := float64(params.InitialPayment)

	annualRate := getAnnualRate(program)
	G := annualRate / 12             // monthly rate
	S := objectCost - initialPayment // mortgage debt (loan sum)

	lastPaymentDate := time.Now().AddDate(0, params.Months, 0)
	T := months // interest periods count

	totalRate := math.Pow(1+G, T)
	PM := math.Ceil(S * G * totalRate / (totalRate - 1)) // monthly payment

	overpayment := PM*months - S

	log.Info(
		"aggregates calculated",
		slog.Any("lastPaymentDate", lastPaymentDate),
		slog.Float64("annualRate", annualRate),
		slog.Float64("G", G),
		slog.Float64("S", S),
		slog.Float64("T", T),
		slog.Float64("PM", PM),
		slog.Float64("overpayment", overpayment),
	)

	return &dto.CalcAggregates{
		LastPaymentDate: lastPaymentDate.Format("2006-01-02"),
		Rate:            int(annualRate * 100),
		LoanSum:         int(S),
		MonthlyPayment:  int(PM),
		Overpayment:     int(overpayment),
	}, nil
}

const salaryRate = 0.08
const militaryRate = 0.09
const baseRate = 0.1

func getAnnualRate(program dto.CalcProgram) float64 {
	switch {
	case program.Salary:
		return salaryRate
	case program.Military:
		return militaryRate
	case program.Base:
		return baseRate
	default:
		return 0
	}
}
