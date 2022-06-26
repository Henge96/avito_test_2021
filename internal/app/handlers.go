package app

import (
	"context"
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

func (a *Core) ChangeBalance(ctx context.Context, change ChangeBalance) (*Wallet, error) {
	wallet, err := a.repo.GetWallet(ctx, change.UserID)
	if err != nil && err != ErrNotFound {
		return nil, fmt.Errorf("a.repo.GetWallet: %w", err)
	}

	if err == ErrNotFound && change.Amount.GreaterThanOrEqual(decimal.NewFromInt(0)) {
		newWallet, err := a.repo.CreateWallet(ctx, change.UserID)
		if err != nil {
			return nil, fmt.Errorf("a.repo.CreateWallet: %w", err)
		}
		wallet.ID = newWallet.ID
	}

	wallet, err = a.repo.Change(ctx, wallet.ID, change.Amount, stReplenishment)
	if err != nil {
		return nil, fmt.Errorf("a.repo.Change: %w", err)
	}

	return wallet, nil
}

func (a *Core) Transfer(ctx context.Context, transfer TransferBetweenUsers) (*TransactionsLists, error) {

	senderWallet, err := a.repo.GetWallet(ctx, transfer.SenderID)
	if err != nil {
		return nil, fmt.Errorf("a.repo.GetWalletByUser: %w", err)
	}

	receiverWallet, err := a.repo.GetWallet(ctx, transfer.ReceiverID)
	if err != nil && err != ErrNotFound {
		return nil, fmt.Errorf("a.repo.GetWalletByUser: %w", err)
	}

	if err == ErrNotFound {
		receiverWallet, err = a.repo.CreateWallet(ctx, transfer.ReceiverID)
		if err != nil {
			return nil, fmt.Errorf("a.repo.CreateWallet: %w", err)
		}
	}

	tr := TransferBetweenUsers{
		SenderID:   senderWallet.UserID,
		Amount:     transfer.Amount,
		ReceiverID: receiverWallet.UserID,
	}

	transaction, err := a.repo.TransactionBetweenUsers(ctx, tr, stTransfer)
	if err != nil {
		return nil, fmt.Errorf("a.repo.TransactionBetweenUsers: %w", err)
	}

	return transaction, nil
}

func (a *Core) GetUserTransactions(ctx context.Context, params UserTransactionsParam) ([]TransactionsLists, int, error) {
	return a.repo.GetUserTransactionsByParams(ctx, params)
}
