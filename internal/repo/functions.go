package repo

import (
	"context"
	"packs/internal/app"
	"time"
)

func (nac Nachinka) GetWalletByUserId(ctx context.Context, userId int) (*app.Wallet, error) {

	var wallet app.Wallet

	row := nac.db.QueryRowContext(ctx, "SELECT * FROM wallet where user_id = $1", userId)
	err := row.Scan(&wallet.Id, &wallet.UserId, &wallet.Balance)

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (nac Nachinka) UpdateBalanceByUserId(ctx context.Context, money float64, userId int) error {
	_, err := nac.db.ExecContext(ctx, "UPDATE wallet SET balance = balance - $1 WHERE User_Id = $2", money, userId)
	if err != nil {
		return err
	}

	return nil
}

func (nac Nachinka) CreateTransactionByUsers(ctx context.Context, userId int, receiverId int, money float64) error {
	_, err := nac.db.ExecContext(ctx, "insert into transaction (wallet_id, receiver_wallet_id, money, date) values ((select id from wallet where user_id = $1), (select id from wallet where user_id = $2), $3, $4)",
		userId, receiverId, money, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (nac Nachinka) GetUserTransactionsByUserId(ctx context.Context, userId int) ([]app.Transaction, error) {
	rows, err := nac.db.QueryContext(ctx, "select * from transaction where wallet_id = (select id from wallet where user_id = $1) order by date, money", userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []app.Transaction{}

	for rows.Next() {
		r := app.Transaction{}
		err := rows.Scan(&r.TransactionId, &r.WalletId, &r.ReceiverWalletId, &r.Money, &r.Date)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil

}

func (nac Nachinka) CreateWalletByUserId(ctx context.Context, UserId int) error {
	_, err := nac.db.ExecContext(ctx, "insert into wallet (user_id, balance) values ($1, 0.0)", UserId)
	if err != nil {
		return err
	}

	return nil

}

func (nac Nachinka) StopConnect() error {
	err := nac.db.Close()
	if err != nil {
		return err
	}
	return nil
}
