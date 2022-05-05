package app

import "context"

type Repo interface {
	GetWalletByUserId(ctx context.Context, i int) (*Wallet, error)
	UpdateBalanceByUserId(ctx context.Context, money float64, userId int) error
	CreateTransactionByUsers(ctx context.Context, userId int, receiverId int, money float64) error
	GetUserTransactionsByUserId(ctx context.Context, userId int) ([]Transaction, error)
	CreateWalletByUserId(ctx context.Context, UserId int) error
}

type Core struct {
	repo Repo
}

func NewApplication(repo Repo) *Core {
	return &Core{repo: repo}
}
