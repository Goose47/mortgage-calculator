// Package requests provides dtos for controller requests.
package requests

import "mortgage-calculator/src/internal/domain/dto"

// CalculateRequest represents payload for Calculate endpoint.
type CalculateRequest struct {
	dto.CalcParams
	Program dto.CalcProgram `json:"program" binding:"required"`
}
