package cacherepos

import (
	"context"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	cachepkg "mortgage-calculator/src/internal/cache"
	"mortgage-calculator/src/internal/domain/dto"
	"mortgage-calculator/src/internal/domain/dto/requests"
	cachemock "mortgage-calculator/src/internal/mocks/cache"
	"testing"
)

func setup() (context.Context, *CalcRepository, *cachemock.MockCache) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	cache := new(cachemock.MockCache)

	return context.Background(), NewCalcRepository(log, cache), cache
}

func TestNewCalcRepository(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	cache := new(cachemock.MockCache)
	repo := NewCalcRepository(log, cache)

	require.NotEmpty(t, repo)
}

func TestGenerateKey(t *testing.T) {
	cases := []struct {
		in   *requests.CalculateRequest
		want string
	}{
		{
			&requests.CalculateRequest{
				CalcParams: dto.CalcParams{
					ObjectCost:     1,
					InitialPayment: 2,
					Months:         3,
				},
				Program: dto.CalcProgram{
					Base: true,
				},
			},
			"{\"object_cost\":1,\"initial_payment\":2,\"months\":3,\"program\":{\"base\":true}}",
		},
		{
			&requests.CalculateRequest{
				CalcParams: dto.CalcParams{},
				Program:    dto.CalcProgram{},
			},
			"{\"object_cost\":0,\"initial_payment\":0,\"months\":0,\"program\":{}}",
		},
	}

	for _, tt := range cases {
		res, err := generateKey(tt.in)
		require.NoError(t, err)
		require.Equal(t, tt.want, res)
	}
}

func TestCalcRepository_Get(t *testing.T) {
	ctx, repo, cache := setup()

	in := &requests.CalculateRequest{
		CalcParams: dto.CalcParams{},
		Program:    dto.CalcProgram{},
	}

	marshalled, err := generateKey(in)
	require.NoError(t, err)

	agg := &dto.CalcAggregates{
		LastPaymentDate: "123",
		Rate:            456,
		LoanSum:         789,
		MonthlyPayment:  10,
		Overpayment:     11,
	}
	aggMarshalled := "{\"last_payment_date\":\"123\",\"rate\":456,\"loan_sum\":789,\"monthly_payment\":10,\"overpayment\":11}"
	cache.On("Get", ctx, marshalled).Return([]byte(aggMarshalled), nil)

	res, err := repo.Get(ctx, in)

	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, *agg, *res)
}

func TestCalcRepository_Get_NonexistentKey(t *testing.T) {
	ctx, repo, cache := setup()

	in := &requests.CalculateRequest{
		CalcParams: dto.CalcParams{},
		Program:    dto.CalcProgram{},
	}

	marshalled, err := generateKey(in)
	require.NoError(t, err)

	cache.On("Get", ctx, marshalled).Return(make([]byte, 0), cachepkg.ErrKeyNotExists)

	res, err := repo.Get(ctx, in)

	require.Error(t, err)
	require.Empty(t, res)
	require.ErrorIs(t, err, cachepkg.ErrKeyNotExists)
}

func TestCalcRepository_Set(t *testing.T) {
	ctx, repo, cache := setup()

	in := &requests.CalculateRequest{
		CalcParams: dto.CalcParams{},
		Program:    dto.CalcProgram{},
	}
	inMarshalled := "{\"object_cost\":0,\"initial_payment\":0,\"months\":0,\"program\":{}}"

	agg := &dto.CalcAggregates{
		LastPaymentDate: "123",
		Rate:            456,
		LoanSum:         789,
		MonthlyPayment:  10,
		Overpayment:     11,
	}
	aggMarshalled := "{\"last_payment_date\":\"123\",\"rate\":456,\"loan_sum\":789,\"monthly_payment\":10,\"overpayment\":11}"

	cache.On("Set", ctx, inMarshalled, []byte(aggMarshalled)).Return(nil)

	err := repo.Set(ctx, in, agg)
	require.NoError(t, err)

	cache.On("Get", ctx, inMarshalled).Return([]byte(aggMarshalled), nil)

	res, err := repo.Get(ctx, in)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, *agg, *res)
}

func TestCalcRepository_List(t *testing.T) {
	ctx, repo, cache := setup()

	list := []*cachepkg.Entry{
		{
			Key: "{}",
			Val: []byte("{}"),
			ID:  0,
		},
	}

	cache.On("List", ctx).Return(list, nil)

	res, err := repo.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, res)
}

func TestCalcRepository_Clear(t *testing.T) {
	ctx, repo, cache := setup()
	cache.On("Clear", ctx).Return(nil)
	repo.Clear(ctx)
}
