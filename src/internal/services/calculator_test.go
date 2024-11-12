package services

import (
	"context"
	"io"
	"log/slog"
	"mortgage-calculator/src/internal/domain/dto"
	"testing"
	"time"
)

func Test_NewCalculatorService(t *testing.T) {
	log := slog.Logger{}
	service := NewCalculatorService(&log)

	if service == nil {
		t.Fatalf("calculator service is nil")
	}
}

func TestCalculatorService_Calculate_HappyPath(t *testing.T) {
	now := time.Now()

	cases := []struct {
		params  dto.CalcParams
		program dto.CalcProgram
		want    *dto.CalcAggregates
	}{
		{
			dto.CalcParams{
				ObjectCost:     5000000,
				InitialPayment: 1000000,
				Months:         240,
			},
			dto.CalcProgram{Salary: true},
			&dto.CalcAggregates{
				LastPaymentDate: now.AddDate(0, 240, 0).Format("2006-01-02"),
				Rate:            8,
				LoanSum:         4000000,
				MonthlyPayment:  33458,
				Overpayment:     4029920,
			},
		},
		{
			dto.CalcParams{
				ObjectCost:     100,
				InitialPayment: 20,
				Months:         12,
			},
			dto.CalcProgram{Military: true},
			&dto.CalcAggregates{
				LastPaymentDate: now.AddDate(0, 12, 0).Format("2006-01-02"),
				Rate:            9,
				LoanSum:         80,
				MonthlyPayment:  7,
				Overpayment:     4,
			},
		},
		{
			dto.CalcParams{
				ObjectCost:     100000000000,
				InitialPayment: 20000000000,
				Months:         1200,
			},
			dto.CalcProgram{Base: true},
			&dto.CalcAggregates{
				LastPaymentDate: now.AddDate(0, 1200, 0).Format("2006-01-02"),
				Rate:            10,
				LoanSum:         80000000000,
				MonthlyPayment:  666698216,
				Overpayment:     720037859200,
			},
		},
	}

	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	service := NewCalculatorService(log)

	for _, tt := range cases {
		res, err := service.Calculate(ctx, tt.params, tt.program)
		if err != nil {
			t.Fatalf("CalculatorService.Calculate error: %s", err.Error())
		}
		if *res != *tt.want {
			t.Fatalf("want %v, got %v", *tt.want, *res)
		}
	}
}

func TestCalculatorService_Calculate_InsufficientInitialPayment(t *testing.T) {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	service := NewCalculatorService(log)

	res, err := service.Calculate(
		ctx,
		dto.CalcParams{
			ObjectCost:     5000000,
			InitialPayment: 999999,
			Months:         240,
		},
		dto.CalcProgram{Base: true},
	)

	if res != nil {
		t.Fatalf("want error, got %v", *res)
	}
	if err == nil {
		t.Fatalf("want error, got nil")
	}
}

func TestGetAnnualRate(t *testing.T) {
	cases := []struct {
		in   dto.CalcProgram
		want float64
	}{
		{dto.CalcProgram{Salary: true}, 0.08},
		{dto.CalcProgram{Military: true}, 0.09},
		{dto.CalcProgram{Base: true}, 0.10},
		{dto.CalcProgram{}, 0},
	}

	for _, tt := range cases {
		if res := getAnnualRate(tt.in); res != tt.want {
			t.Fatalf("want %f, got %f", tt.want, res)
		}
	}
}
