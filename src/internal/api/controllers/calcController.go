// Package controllers provides transport layer logic for application endpoints.
package controllers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"mortgage-calculator/src/internal/domain/dto"
	"mortgage-calculator/src/internal/services"
	"net/http"
)

// Calculator calculates result based on given parameters.
type Calculator interface {
	Calculate(ctx context.Context, params dto.CalcParams, program dto.CalcProgram) (*dto.CalcAggregates, error)
}

// CalcController deals with calculation endpoints.
type CalcController struct {
	calculator Calculator
}

// NewCalcController is a constructor for CalcController.
func NewCalcController(calculator Calculator) *CalcController {
	return &CalcController{
		calculator: calculator,
	}
}

type calculateRequest struct {
	dto.CalcParams
	Program dto.CalcProgram `json:"program" binding:"required"`
}

type calculateResponse struct {
	Params     dto.CalcParams     `json:"params"`
	Program    dto.CalcProgram    `json:"program"`
	Aggregates dto.CalcAggregates `json:"aggregates"`
}

// Calculate validates request params, calculates params and composes result message.
func (con *CalcController) Calculate(c *gin.Context) {
	// validate request
	var in calculateRequest
	err := c.ShouldBindJSON(&in)

	if err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no json payload",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var programCount int
	if in.Program.Base {
		programCount++
	}
	if in.Program.Military {
		programCount++
	}
	if in.Program.Salary {
		programCount++
	}

	// program must be specified
	if programCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "choose program",
		})
		return
	}

	// only one program must be selected
	if programCount > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "choose only 1 program",
		})
		return
	}

	params := dto.CalcParams{
		ObjectCost:     in.ObjectCost,
		InitialPayment: in.InitialPayment,
		Months:         in.Months,
	}

	// calculate result
	res, err := con.calculator.Calculate(
		c.Request.Context(),
		params,
		in.Program,
	)

	if err != nil {
		if errors.Is(err, services.ErrInsufficientInitialPayment) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": services.ErrInsufficientInitialPayment.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to calculate params",
		})
		return
	}

	// compose response
	out := calculateResponse{
		Params:     params,
		Program:    in.Program,
		Aggregates: *res,
	}

	c.JSON(200, out)
}
