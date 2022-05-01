package api

import (
	"context"
	"github.com/gorilla/mux"
	"packs/internal/app"
)

type Application interface {
	CheckBalance(ctx context.Context, userId int) (float64, error)
	UpdateBalance(ctx context.Context, money float64, userId int) error
	CreateTransaction(ctx context.Context, userId int, receiverId int, money float64) error
	CheckWallet(ctx context.Context, userId int) error
	TransferWithWallet(ctx context.Context, userId int, receiverId int, money float64) error
	GetUserTransactions(ctx context.Context, wallet app.Wallet) ([]app.Transaction, error)
}

type Api struct {
	app Application
}

func NewAPI(apl Application) *mux.Router {

	a := &Api{
		apl,
	}

	r := mux.NewRouter()
	r.Use(MwHandler1, MwHandler2)
	r.HandleFunc("/user/wallet/upbalance", a.DepositOrWithdrow).Methods("POST")
	r.HandleFunc("/user/wallet/transfer", a.Transfer).Methods("POST")
	r.HandleFunc("/user/wallet/balance", a.PrintBalance).Methods("GET")
	r.HandleFunc("/user/wallet/transactions", a.PrintTransactions).Methods("GET")

	return r

}
