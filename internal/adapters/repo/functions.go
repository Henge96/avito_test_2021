package repo

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"packs/internal/app"
)

const (
	tableWallet      = "wallet"
	tableTransaction = "transaction"
)

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

func (r Repo) Change(ctx context.Context, walletID uint, amount decimal.Decimal, status string) (*app.Wallet, error) {
	var wallet app.Wallet
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, convertErr(err)
	}

	queryTransaction := fmt.Sprintf("insert into %s (sender_id, receiver_id, amount, status) values ($1, $1, $2, $3)", tableTransaction)
	row, err := tx.ExecContext(ctx, queryTransaction, walletID, amount, status)
	if err != nil {
		tx.Rollback()
		return nil, convertErr(err)
	}

	total, err := row.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, convertErr(err)
	}

	if total != 1 {
		tx.Rollback()
		return nil, fmt.Errorf("total != 1")
	}

	query := fmt.Sprintf("update %s set balance = balance + $1 where id = $2 returning id, user_id, balance", tableWallet)
	err = tx.QueryRowContext(ctx, query, amount, walletID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		tx.Rollback()
		return nil, convertErr(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("tx.Commit: %w", err)
	}

	return &wallet, nil
}

func (r Repo) TransactionBetweenUsers(ctx context.Context, transfer app.TransferBetweenUsers, status string) (*app.TransactionsLists, error) {
	var transaction app.TransactionsLists
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, convertErr(err)
	}

	querySender := fmt.Sprintf("update %s set balance = balance-$1 where user_id = $2", tableWallet)
	rowsUP, err := tx.ExecContext(ctx, querySender, transfer.Amount, transfer.SenderID)
	if err != nil {
		tx.Rollback()
		return nil, convertErr(err)
	}

	rowsSender, err := rowsUP.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, convertErr(err)
	}

	if rowsSender != 1 {
		tx.Rollback()
		return nil, fmt.Errorf("rowsSender != 1")
	}

	queryReceiver := fmt.Sprintf("update %s set balance = balance+$1 where user_id = $2", tableWallet)
	rowsDown, err := tx.ExecContext(ctx, queryReceiver, transfer.Amount, transfer.ReceiverID)
	if err != nil {
		tx.Rollback()
		return nil, convertErr(err)
	}

	rowsReceiver, err := rowsDown.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Receiver.Rows.Affected: %w", err)
	}

	if rowsReceiver != 1 {
		tx.Rollback()
		return nil, fmt.Errorf("rowsReceiver != 1")
	}

	queryTransaction := fmt.Sprintf("insert into %s (sender_id, receiver_id, amount, status) values ((select id from wallet where user_id = $1), (select id from wallet where user_id = $2), $3, $4) returning id, sender_id, receiver_id, amount, created_at, status", tableTransaction)
	row := tx.QueryRowContext(ctx, queryTransaction, transfer.SenderID, transfer.ReceiverID, transfer.Amount, status)
	err = row.Scan(&transaction.ID, &transaction.FromWallet, &transaction.ToWallet, &transaction.Amount, &transaction.CreatedAt, &transaction.Status)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("tx.QueryRowContext: %w", convertErr(err))
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("tx.Commit: %w", err)
	}

	return &transaction, nil
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
