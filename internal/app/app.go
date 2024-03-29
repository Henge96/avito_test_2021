package app

import (
	"context"

	"github.com/shopspring/decimal"
)

type (
	// Repo for current service repository.
	Repo interface {
		GetWallet(ctx context.Context, i uint) (*Wallet, error)
		Change(ctx context.Context, walletID uint, amount decimal.Decimal) (*Wallet, error)
		GetUserTransactionsByParams(ctx context.Context, params UserTransactionsParam) ([]TransactionsLists, int, error)
		CreateWallet(ctx context.Context, userID uint) (*Wallet, error)
		Tx(ctx context.Context, cb func(repo Repo) error) (err error)
		Transaction(ctx context.Context, tr Transaction) (int, error)
	}
	// ExchangeClient for exchange currency at another api.
	ExchangeClient interface {
		ExchangeCurrency(ctx context.Context, sum decimal.Decimal, ticker string) (decimal.Decimal, error)
	}
)

type Core struct {
	repo     Repo
	exchange ExchangeClient
}

func New(repo Repo, exchange ExchangeClient) *Core {
	return &Core{repo: repo, exchange: exchange}
}
