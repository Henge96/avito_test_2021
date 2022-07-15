package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"

	"packs/internal/app"
)

const (
	tableWallet      = "wallet"
	tableTransaction = "transaction"
)

var _ app.Repo = &Repo{}

func (r Repo) CreateWallet(ctx context.Context, userID uint) (*app.Wallet, error) {
	var wallet app.Wallet

	query := fmt.Sprintf("insert into %s (user_id) values ($1) returning id, user_id, balance", tableWallet)
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, fmt.Errorf("r.db.QueryRowContext: %w", err)
	}

	return &wallet, nil
}

func (r Repo) GetWallet(ctx context.Context, userID uint) (*app.Wallet, error) {
	var wallet app.Wallet

	query := fmt.Sprintf("SELECT * FROM %s where user_id = $1", tableWallet)
	row := r.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, convertErr(err)
	}

	return &wallet, nil
}

func (r Repo) Change(ctx context.Context, walletID uint, amount decimal.Decimal) (*app.Wallet, error) {
	var wallet app.Wallet

	query := fmt.Sprintf("update %s set balance = balance + $1 where id = $2 returning id, user_id, balance", tableWallet)
	err := r.db.QueryRowContext(ctx, query, amount, walletID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, convertErr(err)
	}

	return &wallet, nil
}

func (r Repo) GetUserTransactionsByParams(ctx context.Context, params app.UserTransactionsParam) ([]app.TransactionsLists, int, error) {
	query := fmt.Sprintf("select * from %s where sender_id = (select id from %s where user_id = $1) OR receiver_id = (select id from %s where user_id = $1) order by created_at desc, amount limit $2 offset $3", tableTransaction, tableWallet, tableWallet)

	rows, err := r.db.QueryContext(ctx, query, params.UserID, params.Limit, params.Offset)
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

	getTotal := fmt.Sprintf("select count(*) from %s where sender_id = (select id from %s where user_id = $1) OR receiver_id = (select id from %s where user_id = $1)", tableTransaction, tableWallet, tableWallet)

	var total int
	err = r.db.Get(&total, getTotal, params.UserID)
	if err != nil {
		return nil, 0, fmt.Errorf("r.db.Get: %w", convertErr(err))
	}

	return lists, total, nil
}

func (r Repo) StopConnect() error {
	err := r.db.Close()
	if err != nil {
		return fmt.Errorf("r.db.Close: %w", err)
	}
	return nil
}

func (r Repo) Tx(ctx context.Context, f func(app.Repo) error) error {
	opts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	return txHelper(ctx, r.db, opts, func(tx *sqlx.Tx) error {
		return f(&txRepo{tx: tx})
	})
}

func (r Repo) Transaction(ctx context.Context, tr app.Transaction) (*app.TransactionsLists, error) {
	var transaction app.TransactionsLists

	queryTransaction := fmt.Sprintf("insert into %s (sender_id, receiver_id, amount, status) values ((select id from wallet where user_id = $1), (select id from wallet where user_id = $2), $3, $4) returning id, sender_id, receiver_id, amount, created_at, status", tableTransaction)
	err := r.db.GetContext(ctx, &transaction, queryTransaction, tr.SenderID, tr.ReceiverID, tr.Amount, tr.Amount)
	if err != nil {
		return nil, fmt.Errorf("tx.QueryRowContext: %w", convertErr(err))
	}

	return &transaction, nil
}

func txHelper(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions, cb func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("db.BeginTx: %w", err)
	}

	err = cb(tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			err = fmt.Errorf("%w: %s", err, errRollback)
		}
		return err
	}

	return tx.Commit()
}
