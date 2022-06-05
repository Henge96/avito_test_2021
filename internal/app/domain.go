package app

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	Id         int             `json:"Id"`
	UserId     int             `json:"UserId"`
	Balance    decimal.Decimal `json:"Balance"`
	Currency   string          `json:"Currency"`
	ReceiverId int             `json:"ReceiverId"`
	Transaction
}

type Transaction struct {
	TransactionId    int             `json:"TransactionId"`
	WalletId         int             `json:"wallet_id"`
	ReceiverWalletId int             `json:"ReceiverWalletId"`
	Money            decimal.Decimal `json:"Money"`
	Date             time.Time       `json:"Date"`
}
