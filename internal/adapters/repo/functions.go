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

func (r Repo) GetWalletByUserID(ctx context.Context, userID int) (wallet *app.Wallet, err error) {
	query := fmt.Sprintf("SELECT * FROM %s where user_id = $1", tableWallet)

	row := r.db.QueryRowContext(ctx, query, userID)
	err = row.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, fmt.Errorf("r.db.QueryRowContext: %w", convertErr(err))
	}

	return wallet, nil
}

func (r Repo) UpdateBalanceByUserId(ctx context.Context, money decimal.Decimal, userId int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE wallet SET balance = balance - $1 WHERE User_Id = $2", money, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r Repo) CreateTransactionByUsers(ctx context.Context, userId int, receiverId int, money decimal.Decimal) error {
	_, err := r.db.ExecContext(ctx, "insert into transaction (wallet_id, receiver_wallet_id, money) values ((select id from wallet where user_id = $1), (select id from wallet where user_id = $2), $3)",
		userId, receiverId, money)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetUserTransactionsByParams(ctx context.Context, params app.UserTransactionsParam) (lists []app.TransactionsLists, total int, err error) {
	query := fmt.Sprintf("select * from %s where wallet_id = (select id from %s where user_id = $1) OR receiver_wallet_id = (select id from %s where user_id = $1) order by date desc, money limit $2 offset $3", tableTransaction, tableWallet, tableWallet)

	rows, err := r.db.QueryContext(ctx, query, params.UserID, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("r.db.QueryContext: %w", convertErr(err))
	}
	defer rows.Close()

	for rows.Next() {
		list := app.TransactionsLists{}
		err := rows.Scan(&list.ID, &list.FromWallet, &list.ToWallet, &list.Amount, &list.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		lists = append(lists, r)
	}

	return result, nil

}

func (r Repo) CreateWalletByUserId(ctx context.Context, UserId int) error {
	_, err := r.db.ExecContext(ctx, "insert into wallet (user_id, balance) values ($1, 0.0)", UserId)
	if err != nil {
		return err
	}

	return nil

}

func (r Repo) StopConnect() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}
