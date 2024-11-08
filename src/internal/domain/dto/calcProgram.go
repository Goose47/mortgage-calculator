package dto

// CalcProgram represents available programs for calculation.
type CalcProgram struct {
	Salary   bool `json:"salary,omitempty"`
	Military bool `json:"military,omitempty"`
	Base     bool `json:"base,omitempty"`
}
