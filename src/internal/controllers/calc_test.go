package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	cachepkg "mortgage-calculator/src/internal/cache"
	"mortgage-calculator/src/internal/domain/dto"
	"mortgage-calculator/src/internal/domain/dto/requests"
	reposmock "mortgage-calculator/src/internal/mocks/repos"
	servicesmock "mortgage-calculator/src/internal/mocks/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() (*CalcController, *servicesmock.MockCalculator, *reposmock.MockCacheGetSaver) {
	service := new(servicesmock.MockCalculator)
	repo := new(reposmock.MockCacheGetSaver)
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	con := NewCalcController(log, service, repo)

	return con, service, repo
}

func TestNewCalcController(t *testing.T) {
	service := new(servicesmock.MockCalculator)
	repo := new(reposmock.MockCacheGetSaver)
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	con := NewCalcController(log, service, repo)

	require.NotEmpty(t, con)
}

func TestCalcController_Calculate(t *testing.T) {
	con, s, r := setup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, err := json.Marshal(requests.CalculateRequest{
		CalcParams: dto.CalcParams{
			ObjectCost:     100,
			InitialPayment: 20,
			Months:         12,
		},
		Program: dto.CalcProgram{
			Salary: true,
		},
	})
	require.NoError(t, err)
	req, _ := http.NewRequest("POST", "/calculate", bytes.NewBuffer(body))
	c.Request = req

	r.On("Get", mock.Anything, mock.Anything).Return(&dto.CalcAggregates{}, cachepkg.ErrKeyNotExists)
	r.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.On("Calculate", mock.Anything, mock.Anything, mock.Anything).Return(&dto.CalcAggregates{
		LastPaymentDate: "1",
		Rate:            2,
		LoanSum:         3,
		MonthlyPayment:  4,
		Overpayment:     5,
	}, nil)

	con.Calculate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	require.Contains(
		t,
		w.Body.String(),
		"{\"aggregates\":{\"last_payment_date\":\"1\",\"rate\":2,\"loan_sum\":3,\"monthly_payment\":4,\"overpayment\":5},\"params\":{\"object_cost\":100,\"initial_payment\":20,\"months\":12},\"program\":{\"salary\":true}}",
	)
}

func TestValidateRequest_PassCases(t *testing.T) {
	cases := []struct {
		in requests.CalculateRequest
	}{
		{
			requests.CalculateRequest{
				CalcParams: dto.CalcParams{
					ObjectCost:     100,
					InitialPayment: 20,
					Months:         12,
				},
				Program: dto.CalcProgram{
					Salary: true,
				},
			},
		},
		{
			requests.CalculateRequest{
				CalcParams: dto.CalcParams{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         120,
				},
				Program: dto.CalcProgram{
					Base: true,
				},
			},
		},
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	for _, tt := range cases {
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(tt.in)
		req, _ := http.NewRequest("POST", "/calculate", bytes.NewBuffer(body))
		c.Request = req
		data, err := validateRequest(c)

		require.NoError(t, err)
		require.NotEmpty(t, data)
		require.Equal(t, tt.in, *data)
	}
}

func TestValidateRequest_FailCases(t *testing.T) {
	cases := []struct {
		in  requests.CalculateRequest
		err error
	}{
		{
			requests.CalculateRequest{
				CalcParams: dto.CalcParams{
					ObjectCost:     100,
					InitialPayment: 20,
					Months:         120,
				},
				Program: dto.CalcProgram{
					Salary: true,
					Base:   true,
				},
			},
			errTooManyPrograms,
		},
		{
			requests.CalculateRequest{
				CalcParams: dto.CalcParams{
					ObjectCost:     100,
					InitialPayment: 20,
					Months:         120,
				},
				Program: dto.CalcProgram{},
			},
			errNoProgram,
		},
		{
			requests.CalculateRequest{
				CalcParams: dto.CalcParams{},
				Program:    dto.CalcProgram{},
			},
			errValidation,
		},
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	for _, tt := range cases {
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(tt.in)
		req, _ := http.NewRequest("POST", "/calculate", bytes.NewBuffer(body))
		c.Request = req
		data, err := validateRequest(c)

		require.Error(t, err)
		require.Empty(t, data)
		require.ErrorIs(t, err, tt.err)
	}
}
