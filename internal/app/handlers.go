package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

func (a *Core) GetUserBalance(ctx context.Context, userID uint, currency string) (*Wallet, error) {
	wallet, err := a.repo.GetWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("a.repo.GetWalletByUserId: %w", err)
	}

	if currency != "RUB" {
		excBalance, err := a.exchange.ExchangeCurrency(ctx, wallet.Balance, currency)
		if err != nil {
			return nil, fmt.Errorf("a.exchange.ExchangeCurrency: %w", err)
		}
		wallet.Balance = excBalance
	}

	return wallet, nil
}

func (a *Core) ChangeBalance(ctx context.Context, change ChangeBalance) (wallet *Wallet, err error) {
	wallet, err = a.repo.GetWallet(ctx, change.UserID)
	switch {
	case errors.Is(err, ErrNotFound):
		wallet, err = a.repo.CreateWallet(ctx, change.UserID)
		if err != nil {
			return nil, fmt.Errorf("a.repo.CreateWallet: %w", err)
		}
	case err != nil:
		return nil, fmt.Errorf("a.repo.GetWallet: %w", err)
	}

	err = a.repo.Tx(ctx, func(repo Repo) error {
		wallet, err = repo.Change(ctx, wallet.ID, change.Amount)
		if err != nil {
			return fmt.Errorf("a.repo.Change: %w", err)
		}

		tr := Transaction{
			SenderID:   wallet.ID,
			ReceiverID: wallet.ID,
			Amount:     change.Amount,
			Status:     stReplenishment,
		}

		_, err = repo.Transaction(ctx, tr)
		if err != nil {
			return fmt.Errorf("repo.Transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (a *Core) Transfer(ctx context.Context, transfer Transaction) (id int, err error) {

	senderWallet, err := a.repo.GetWallet(ctx, transfer.SenderID)
	if err != nil {
		return 0, fmt.Errorf("a.repo.GetWalletByUser: %w", err)
	}

	receiverWallet, err := a.repo.GetWallet(ctx, transfer.ReceiverID)
	switch {
	case errors.Is(err, ErrNotFound) && transfer.Amount.GreaterThan(decimal.NewFromInt(0)):
		receiverWallet, err = a.repo.CreateWallet(ctx, transfer.ReceiverID)
		if err != nil {
			return 0, fmt.Errorf("a.repo.CreateWallet: %w", err)
		}
	case err != nil:
		return 0, fmt.Errorf("a.repo.GetWalletByUser: %w", err)
	}

	err = a.repo.Tx(ctx, func(repo Repo) error {
		senderWallet, err = repo.Change(ctx, senderWallet.ID, transfer.Amount.Mul(decimal.NewFromInt(-1)))
		if err != nil {
			return fmt.Errorf("a.repo.TransactionBetweenUsers: %w", err)
		}

		receiverWallet, err = repo.Change(ctx, receiverWallet.ID, transfer.Amount)
		if err != nil {
			return fmt.Errorf("a.repo.TransactionBetweenUsers: %w", err)
		}

		tr := Transaction{
			SenderID:   senderWallet.ID,
			Amount:     transfer.Amount,
			ReceiverID: receiverWallet.ID,
			Status:     stTransfer,
		}

		id, err = repo.Transaction(ctx, tr)
		if err != nil {
			return fmt.Errorf("repo.Transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *Core) GetUserTransactions(ctx context.Context, params UserTransactionsParam) ([]TransactionsLists, int, error) {
	return a.repo.GetUserTransactionsByParams(ctx, params)
}
