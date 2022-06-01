package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/shopspring/decimal"

	"packs/internal/app"
)

func MwHandler1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")
		log.Println(r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func MwHandler2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")
		log.Println(r.Method)
		next.ServeHTTP(w, r)
	})
}

func ErrHandler(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	log.Println(err)
}

func (a *Api) PrintBalance(w http.ResponseWriter, r *http.Request) {

	var balance app.Wallet

	err := json.NewDecoder(r.Body).Decode(&balance)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	fmt.Println(balance.Money)

	ReturnBalance, err := a.app.CheckBalance(r.Context(), balance.UserId, balance.Currency)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	result, _ := ReturnBalance.Float64()

	err = json.NewEncoder(w).Encode(CheckBalanceResp{RetBalance: result})
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}
}

func (a *Api) Transfer(w http.ResponseWriter, r *http.Request) {

	var wallet app.Wallet

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	if wallet.Money.IsNegative() || wallet.ReceiverId == wallet.UserId {
		ErrHandler(w, err, 400)
		return
	}

	err = a.app.CheckWallet(r.Context(), wallet.UserId)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	err = a.app.CheckWallet(r.Context(), wallet.ReceiverId)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	err = a.app.TransferWithWallet(r.Context(), wallet.UserId, wallet.ReceiverId, wallet.Money)
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}

	ReturnBalance, err := a.app.CheckBalance(r.Context(), wallet.UserId, wallet.Currency)
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}
	result, _ := ReturnBalance.Float64()

	err = json.NewEncoder(w).Encode(CheckBalanceResp{RetBalance: result})
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}
}

func (a *Api) PrintTransactions(w http.ResponseWriter, r *http.Request) {

	var wallet app.Wallet

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	res, err := a.app.GetUserTransactions(r.Context(), wallet)
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

}

func (a *Api) DepositOrWithdrow(w http.ResponseWriter, r *http.Request) {

	var wallet app.Wallet

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	wallet.ReceiverId = wallet.UserId

	if wallet.Money.IsZero() {
		ErrHandler(w, err, 400)
		return
	} else if wallet.Money.IsPositive() {

		err = a.app.CheckWallet(r.Context(), wallet.UserId)
		if err != nil && false == errors.Is(err, sql.ErrNoRows) {
			ErrHandler(w, err, 400)
			return
		} else if true == errors.Is(err, sql.ErrNoRows) {
			err = a.app.CreateWallet(r.Context(), wallet.UserId)
			if err != nil {
				ErrHandler(w, err, 500)
				return
			}

		}

		err = a.app.UpdateBalance(r.Context(), wallet.Money.Mul(decimal.NewFromInt(-1)), wallet.UserId)
		if err != nil {
			ErrHandler(w, err, 500)
			return
		}

		err = a.app.CreateTransaction(r.Context(), wallet.UserId, wallet.ReceiverId, wallet.Money)
		if err != nil {
			ErrHandler(w, err, 500)
			return
		}

	} else if wallet.Money.IsNegative() {

		err = a.app.CheckWallet(r.Context(), wallet.UserId)
		if err != nil {
			ErrHandler(w, err, 400)
			return
		}

		err = a.app.UpdateBalance(r.Context(), wallet.Money.Mul(decimal.NewFromInt(-1)), wallet.UserId)
		if err != nil {
			ErrHandler(w, err, 500)
			return
		}

		err = a.app.CreateTransaction(r.Context(), wallet.UserId, wallet.ReceiverId, wallet.Money.Mul(decimal.NewFromInt(-1)))
		if err != nil {
			ErrHandler(w, err, 500)
			return
		}

	}

	ReturnBalance, err := a.app.CheckBalance(r.Context(), wallet.UserId, wallet.Currency)
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}
	result, _ := ReturnBalance.Float64()

	err = json.NewEncoder(w).Encode(CheckBalanceResp{RetBalance: result})
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

}
