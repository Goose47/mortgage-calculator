package dto

// CalcParams represent parameters required for calculation.
type CalcParams struct {
	ObjectCost     int `json:"object_cost" binding:"required"`
	InitialPayment int `json:"initial_payment" binding:"required"`
	Months         int `json:"months" binding:"required"`
}
