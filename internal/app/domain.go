package app

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	stReplenishment = "replenishment"
	stTransfer      = "transfer"
)

type (
	RequestBalance struct {
		UserID   uint   `json:"user_id" validate:"required"`
		Currency string `json:"currency" validate:"required,len=3,alpha"`
	}

	ChangeBalance struct {
		UserID uint            `json:"user_id" validate:"required"`
		Amount decimal.Decimal `json:"amount" validate:"required"`
	}

	Wallet struct {
		ID      uint
		UserID  int
		Balance decimal.Decimal
	}

	TransferBetweenUsers struct {
		SenderID   uint            `json:"sender_id" validate:"required"`
		ReceiverID uint            `json:"receiver_id" validate:"required"`
		Amount     decimal.Decimal `json:"amount" validate:"required,gt=0"`
	}

	TransactionsLists struct {
		ID         uint            `json:"id"`
		FromWallet uint            `json:"from_wallet"`
		ToWallet   uint            `json:"to_wallet"`
		Amount     decimal.Decimal `json:"amount"`
		CreatedAt  time.Time       `json:"created_at"`
		Status     string          `json:"status"`
	}

	UserTransactionsParam struct {
		UserID uint `json:"user_id" validate:"required"`
		Limit  uint `json:"limit" validate:"required"`
		Offset uint `json:"offset"`
	}
)
