package app_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"

	"packs/internal/app"
)

func TestCore_GetUserBalance(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())

	wallet := app.Wallet{
		ID:      uint(rand.Uint32()),
		UserID:  uint(rand.Uint64()),
		Balance: decimal.NewFromInt(100),
	}

	walletUSD := app.Wallet{
		ID:      uint(rand.Uint32()),
		UserID:  uint(rand.Uint64()),
		Balance: decimal.NewFromInt(500),
	}
	exchangeRes := walletUSD.Balance.Div(decimal.NewFromInt(100))

	testCases := map[string]struct {
		userID      uint
		currency    string
		repoRes     *app.Wallet
		repoErr     error
		exchangeRes decimal.Decimal
		exchangeErr error
		want        *app.Wallet
		wantErr     error
	}{
		"successRUB": {wallet.UserID, "RUB", &wallet, nil, decimal.Decimal{}, nil, &wallet, nil},
		"successUSD": {walletUSD.UserID, "USD", &walletUSD, nil, exchangeRes, nil, &app.Wallet{ID: walletUSD.ID, UserID: walletUSD.UserID, Balance: exchangeRes}, nil},
	}

	for name, tc := range testCases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mock, assert := start(t)

			mock.repo.EXPECT().GetWallet(gomock.Any(), tc.userID).Return(tc.repoRes, tc.repoErr)

			if tc.currency != "RUB" {
				tc.repoRes.Balance = tc.exchangeRes

				mock.exchange.EXPECT().ExchangeCurrency(gomock.Any(), tc.repoRes.Balance, tc.currency).Return(tc.exchangeRes, tc.exchangeErr)
			}

			res, err := c.GetUserBalance(ctx, tc.userID, tc.currency)
			assert.ErrorIs(err, tc.wantErr)
			assert.Equal(tc.want, res)
		})
	}
}
