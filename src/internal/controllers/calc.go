// Package controllers provides transport layer logic for application endpoints.
package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"mortgage-calculator/src/internal/domain/dto"
	"mortgage-calculator/src/internal/domain/dto/requests"
	"mortgage-calculator/src/internal/services"
	"net/http"
)

// Calculator calculates result based on given parameters.
type Calculator interface {
	Calculate(ctx context.Context, params dto.CalcParams, program dto.CalcProgram) (*dto.CalcAggregates, error)
}

// CacheGetSaver interacts with cache.
type CacheGetSaver interface {
	Get(ctx context.Context, in *requests.CalculateRequest) (*dto.CalcAggregates, error)
	Set(ctx context.Context, in *requests.CalculateRequest, aggregates *dto.CalcAggregates) error
}

// CalcController deals with calculation endpoints.
type CalcController struct {
	log        *slog.Logger
	calculator Calculator
	cache      CacheGetSaver
}

// NewCalcController is a constructor for CalcController.
func NewCalcController(
	log *slog.Logger,
	calculator Calculator,
	cache CacheGetSaver,
) *CalcController {
	return &CalcController{
		log:        log,
		calculator: calculator,
		cache:      cache,
	}
}

type calculateResponse struct {
	Aggregates dto.CalcAggregates `json:"aggregates"`
	Params     dto.CalcParams     `json:"params"`
	Program    dto.CalcProgram    `json:"program"`
}

// Calculate validates request params, calculates params and composes result message.
func (con *CalcController) Calculate(c *gin.Context) {
	ctx := c.Request.Context()

	// validate request
	in, err := validateRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	params := dto.CalcParams{
		ObjectCost:     in.ObjectCost,
		InitialPayment: in.InitialPayment,
		Months:         in.Months,
	}

	// retrieve result from cache
	res, err := con.cache.Get(ctx, in)
	if err != nil {
		// calculate result
		res, err = con.calculator.Calculate(ctx, params, in.Program)

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

		// save calculated result
		if err := con.cache.Set(ctx, in, res); err != nil {
			con.log.Warn("failed to cache result", slog.Any("error", err))
		}
	}

	// compose response
	out := calculateResponse{
		Params:     params,
		Program:    in.Program,
		Aggregates: *res,
	}

	c.JSON(200, out)
}

var errNoPayload = errors.New("no json payload")
var errValidation = errors.New("validation error")
var errNoProgram = errors.New("choose program")
var errTooManyPrograms = errors.New("choose only 1 program")

func validateRequest(c *gin.Context) (*requests.CalculateRequest, error) {
	// validate request
	var in requests.CalculateRequest
	err := c.ShouldBindJSON(&in)

	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errNoPayload
		}
		return nil, fmt.Errorf("%w: %s", errValidation, err.Error())
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
		return nil, errNoProgram
	}

	// only one program must be selected
	if programCount > 1 {
		return nil, errTooManyPrograms
	}

	return &in, nil
}
