package api

import (
	"time"

	"github.com/shopspring/decimal"

	"packs/internal/app"
)

type ResponseBalance struct {
	Balance string `json:"balance"`
}

type ResponseTransactions struct {
	Response []app.TransactionsLists `json:"response"`
	Total    int                     `json:"total"`
}

type Response struct {
	ID         uint            `json:"id"`
	FromWallet uint            `json:"from_wallet"`
	ToWallet   uint            `json:"to_wallet"`
	Amount     decimal.Decimal `json:"amount"`
	CreatedAt  time.Time       `json:"created_at"`
	Status     string          `json:"status"`
}

type ResponseTransactionID struct {
	ID int `json:"id"`
}
