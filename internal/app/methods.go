package app

import (
	"context"
)

func (a *Core) CheckBalance(ctx context.Context, userId int) (float64, error) {
	wallet, err := a.repo.GetWalletByUserId(ctx, userId)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (a *Core) UpdateBalance(ctx context.Context, money float64, userId int) error {

	err := a.repo.UpdateBalanceByUserId(ctx, money, userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *Core) CreateTransaction(ctx context.Context, userId int, receiverId int, money float64) error {
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

func (a *Core) TransferWithWallet(ctx context.Context, userId int, receiverId int, money float64) error {

	err := a.repo.CreateTransactionByUsers(ctx, userId, receiverId, money)
	if err != nil {
		return err
	}

	err = a.repo.UpdateBalanceByUserId(ctx, money, userId)
	if err != nil {
		return err
	}

	err = a.repo.UpdateBalanceByUserId(ctx, money*-1, receiverId)
	if err != nil {
		return err
	}

	return nil

}

func (a *Core) GetUserTransactions(ctx context.Context, wallet Wallet) ([]Transaction, error) {

	return a.repo.GetUserTransactionsByUserId(ctx, wallet.UserId)
}
