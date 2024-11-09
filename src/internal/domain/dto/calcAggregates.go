// Package dto provides dtos used across whole domain.
package dto

import "time"

// CalcAggregates represents calculation result.
type CalcAggregates struct {
	Rate            int       `json:"rate"`
	LoanSum         int       `json:"loan_sum"`
	MonthlyPayment  int       `json:"monthly_payment"`
	Overpayment     int       `json:"overpayment"`
	LastPaymentDate time.Time `json:"last_payment_date"`
}