package app

import (
	"time"

	"github.com/shopspring/decimal"
)

type RequestBalance struct {
	UserID   int    `json:"user_id"`
	Currency string `json:"currency"`
}

type Wallet struct {
	ID      int
	UserID  int
	Balance decimal.Decimal
}

type TransferBetweenUsers struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
	Amount     int `json:"amount"`
}

type TransactionsLists struct {
	ID         int             `json:"id"`
	FromWallet int             `json:"from_wallet"`
	ToWallet   int             `json:"to_wallet"`
	Amount     decimal.Decimal `json:"amount"`
	CreatedAt  time.Time       `json:"created_at"`
}

type UserTransactionsParam struct {
	UserID int
	Limit  uint
	Offset uint
}
