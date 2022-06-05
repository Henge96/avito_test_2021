package app

import (
	"context"

	"github.com/shopspring/decimal"
)

type (
	Repo interface {
		GetWalletByUserId(ctx context.Context, i int) (*Wallet, error)
		UpdateBalanceByUserId(ctx context.Context, money decimal.Decimal, userId int) error
		CreateTransactionByUsers(ctx context.Context, userId int, receiverId int, money decimal.Decimal) error
		GetUserTransactionsByUserId(ctx context.Context, userId int) ([]Transaction, error)
		CreateWalletByUserId(ctx context.Context, UserId int) error
	}
	ExchangeClient interface {
		ExchangeCurrency(ctx context.Context, sum decimal.Decimal, ticker string) (decimal.Decimal, error)
	}
)

type Core struct {
	repo     Repo
	exchange ExchangeClient
}

func NewApplication(repo Repo, exchange ExchangeClient) *Core {
	return &Core{repo: repo, exchange: exchange}
}
