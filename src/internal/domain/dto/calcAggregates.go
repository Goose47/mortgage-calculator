// Package dto provides dtos used across whole domain.
package dto

// CalcAggregates represents calculation result.
type CalcAggregates struct {
	LastPaymentDate string `json:"last_payment_date"`
	Rate            int    `json:"rate"`
	LoanSum         int    `json:"loan_sum"`
	MonthlyPayment  int    `json:"monthly_payment"`
	Overpayment     int    `json:"overpayment"`
}
