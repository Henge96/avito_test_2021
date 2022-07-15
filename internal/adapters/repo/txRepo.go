package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"

	"packs/internal/app"
)

type txRepo struct {
	tx *sqlx.Tx
}

func (t txRepo) Transaction(ctx context.Context, tr app.Transaction) (*app.TransactionsLists, error) {
	//TODO implement me
	panic("implement me")
}

func (t txRepo) GetWallet(ctx context.Context, userID uint) (*app.Wallet, error) {
	var wallet app.Wallet

	query := fmt.Sprintf("SELECT * FROM %s where user_id = $1 for update", tableWallet)
	err := t.tx.QueryRowContext(ctx, query, userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, convertErr(err)
	}

	return &wallet, nil
}

func (t txRepo) Change(ctx context.Context, walletID uint, amount decimal.Decimal) (*app.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (t txRepo) GetUserTransactionsByParams(ctx context.Context, params app.UserTransactionsParam) ([]app.TransactionsLists, int, error) {
	//TODO implement me
	panic("implement me")
}

func (t txRepo) CreateWallet(ctx context.Context, userID uint) (*app.Wallet, error) {
	var wallet app.Wallet

	query := fmt.Sprintf("insert into %s (user_id) values ($1) returning id, user_id, balance", tableWallet)
	err := t.tx.QueryRowContext(ctx, query, userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, fmt.Errorf("r.db.QueryRowContext: %w", err)
	}

	return &wallet, nil
}

func (t txRepo) Tx(ctx context.Context, f func(repo app.Repo) error) error {
	panic("implement me")
}

var _ app.Repo = &txRepo{}
