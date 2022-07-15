package api

import (
	"context"

	"github.com/gorilla/mux"

	"packs/internal/app"
)

type Application interface {
	GetUserBalance(ctx context.Context, userId uint, currency string) (*app.Wallet, error)
	ChangeBalance(ctx context.Context, change app.ChangeBalance) (*app.Wallet, error)
	Transfer(ctx context.Context, transfer app.Transaction) (*app.TransactionsLists, error)
	GetUserTransactions(ctx context.Context, params app.UserTransactionsParam) ([]app.TransactionsLists, int, error)
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
	r.HandleFunc("/user/wallet/upbalance", a.ChangeUserBalance).Methods("POST")
	r.HandleFunc("/user/wallet/transfer", a.TransferBetweenWallets).Methods("POST")
	r.HandleFunc("/user/wallet/balance", a.GetBalance).Methods("GET")
	r.HandleFunc("/user/wallet/transactions", a.GetTransactions).Methods("GET")

	return r

}
