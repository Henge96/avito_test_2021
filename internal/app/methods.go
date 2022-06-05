package app

import (
	"context"

	"github.com/shopspring/decimal"
)

func (a *Core) CheckBalance(ctx context.Context, userId int, currency string) (decimal.Decimal, error) {
	wallet, err := a.repo.GetWalletByUserId(ctx, userId)
	if err != nil {
		return decimal.Decimal{}, err
	}

	switch {
	case currency != "" && currency != "RUB":
		res, err := a.exchange.ExchangeCurrency(ctx, wallet.Balance, currency)
		if err != nil {
			return decimal.Decimal{}, err
		}

		return res, nil

	default:
		return wallet.Balance, nil
	}
}

func (a *Core) UpdateBalance(ctx context.Context, money decimal.Decimal, userId int) error {

	err := a.repo.UpdateBalanceByUserId(ctx, money, userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *Core) CreateTransaction(ctx context.Context, userId int, receiverId int, money decimal.Decimal) error {
	err := a.repo.CreateTransactionByUsers(ctx, userId, receiverId, money)
	if err != nil {
		return err
	}
	return nil
}

func (a *Core) CheckWallet(ctx context.Context, userId int) error {
	_, err := a.repo.GetWalletByUserId(ctx, userId)

	if err != nil {
		return err
	}
	return nil
}

func (a *Core) TransferWithWallet(ctx context.Context, userId int, receiverId int, money decimal.Decimal) error {

	err := a.repo.CreateTransactionByUsers(ctx, userId, receiverId, money)
	if err != nil {
		return err
	}

	err = a.repo.UpdateBalanceByUserId(ctx, money, userId)
	if err != nil {
		return err
	}

	err = a.repo.UpdateBalanceByUserId(ctx, money.Mul(decimal.NewFromInt(-1)), receiverId)
	if err != nil {
		return err
	}

	return nil

}

func (a *Core) GetUserTransactions(ctx context.Context, wallet Wallet) ([]Transaction, error) {

	return a.repo.GetUserTransactionsByUserId(ctx, wallet.UserId)
}

func (a *Core) CreateWallet(ctx context.Context, userId int) error {
	return a.repo.CreateWalletByUserId(ctx, userId)
}
