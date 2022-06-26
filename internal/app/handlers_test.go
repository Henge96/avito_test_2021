package app_test

import (
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"math/rand"
	"packs/internal/app"
	"testing"
	"time"
)

func TestCore_GetUserBalance(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())

	wallet := app.Wallet{
		ID:      uint(rand.Uint32()),
		UserID:  uint(rand.Uint64()),
		Balance: decimal.NewFromInt(100),
	}

	testCases := map[string]struct {
		userID      uint
		currency    string
		repoRes     app.Wallet
		repoErr     error
		exchangeRes decimal.Decimal
		exchangeErr error
		want        app.Wallet
		wantErr     error
	}{
		"successRUB": {wallet.UserID, "RUB", wallet, nil, decimal.Decimal{}, nil, wallet, nil},
	}

	for name, tc := range testCases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mock, assert := start(t)

			mock.repo.EXPECT().GetWallet(gomock.Any(), tc.userID).Return(tc.repoRes, tc.repoErr)

			res, err := c.GetUserBalance(ctx, tc.userID, tc.currency)
			assert.ErrorIs(err, tc.wantErr)
			assert.Equal(tc.want, res)
		})
	}
}
