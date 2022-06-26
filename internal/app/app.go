package app

import (
	"context"

	"github.com/shopspring/decimal"
)

type (
	Repo interface {
		GetWallet(ctx context.Context, i uint) (*Wallet, error)
		Change(ctx context.Context, walletID uint, amount decimal.Decimal, status string) (*Wallet, error)
		TransactionBetweenUsers(ctx context.Context, transfer TransferBetweenUsers, status string) (transaction *TransactionsLists, err error)
		GetUserTransactionsByParams(ctx context.Context, params UserTransactionsParam) ([]TransactionsLists, int, error)
		CreateWallet(ctx context.Context, userID uint) (*Wallet, error)
	}
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
