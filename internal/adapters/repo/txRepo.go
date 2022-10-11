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

func (t txRepo) Transaction(ctx context.Context, tr app.Transaction) (int, error) {
	var id int

	const query = `insert into transaction (sender_id, receiver_id, amount, status) values ($1, $2, $3, $4) returning id`
	err := t.tx.GetContext(ctx, &id, query, tr.SenderID, tr.ReceiverID, tr.Amount, tr.Status)
	if err != nil {
		return 0, fmt.Errorf("tx.QueryRowContext: %w", convertErr(err))
	}

	return id, nil
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
	var wallet app.Wallet

	query := fmt.Sprintf("update %s set balance = balance + $1 where id = $2 returning id, user_id, balance", tableWallet)
	err := t.tx.QueryRowContext(ctx, query, amount, walletID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, convertErr(err)
	}

	return &wallet, nil
}

func (t txRepo) GetUserTransactionsByParams(ctx context.Context, params app.UserTransactionsParam) ([]app.TransactionsLists, int, error) {
	query := fmt.Sprintf("select * from %s where sender_id = (select id from %s where user_id = $1) OR receiver_id = (select id from %s where user_id = $1) order by created_at desc, amount limit $2 offset $3 for update", tableTransaction, tableWallet, tableWallet)

	rows, err := t.tx.QueryContext(ctx, query, params.UserID, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("r.db.QueryContext: %w", convertErr(err))
	}
	defer rows.Close()

	var lists []app.TransactionsLists
	for rows.Next() {
		list := app.TransactionsLists{}
		err := rows.Scan(&list.ID, &list.FromWallet, &list.ToWallet, &list.Amount, &list.CreatedAt, &list.Status)
		if err != nil {
			return nil, 0, convertErr(err)
		}
		lists = append(lists, list)
	}

	getTotal := fmt.Sprintf("select count(*) from %s where sender_id = (select id from %s where user_id = $1) OR receiver_id = (select id from %s where user_id = $1) for update", tableTransaction, tableWallet, tableWallet)

	var total int
	err = t.tx.Get(&total, getTotal, params.UserID)
	if err != nil {
		return nil, 0, fmt.Errorf("r.db.Get: %w", convertErr(err))
	}

	return lists, total, nil
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

func (t txRepo) Tx(_ context.Context, _ func(repo app.Repo) error) error {
	panic("couldn`t start transaction into transaction")
}

var _ app.Repo = &txRepo{}
