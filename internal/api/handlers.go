package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
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
	// мб http.StatusText(code) ?
	// спросить в чем разница при выводе err или err.Error()
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

	err = a.app.CheckWallet(r.Context(), balance.UserId)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	ReturnBalance, err := a.app.CheckBalance(r.Context(), balance.UserId)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	err = json.NewEncoder(w).Encode(CheckBalanceResp{RetBalance: ReturnBalance})
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}

	/* balance.Balance = balance.Balance / 74.66

	 как то подвязать валюту надо, про анмаршал не очень понял

	client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get("http://api.exchangeratesapi.io/v1/latest?access_key=96ac3090873
	80f45dd00b0af6b9657c5&symbols=USD")
		if err != nil {
		log.Println(err)
		return
	}

	*/

}

func (a *Api) Transfer(w http.ResponseWriter, r *http.Request) {

	var wallet app.Wallet

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	// уточнить нужна ли двойная проверка + логика чисел у "money"
	if wallet.Money < 0 || wallet.ReceiverId == wallet.UserId {
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

	// не преобразованое число из бд
	ReturnBalance, err := a.app.CheckBalance(r.Context(), wallet.UserId)
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}

	err = json.NewEncoder(w).Encode(CheckBalanceResp{RetBalance: ReturnBalance})
	if err != nil {
		ErrHandler(w, err, 500)
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

	err = a.app.CheckWallet(r.Context(), wallet.UserId)
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

	if wallet.Money == 0 {
		ErrHandler(w, err, 400)
		return
	} else if wallet.Money > 0 {

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

		err = a.app.UpdateBalance(r.Context(), wallet.Money*-1, wallet.UserId)
		if err != nil {
			ErrHandler(w, err, 500)
			return
		}

		err = a.app.CreateTransaction(r.Context(), wallet.UserId, wallet.ReceiverId, wallet.Money)
		if err != nil {
			ErrHandler(w, err, 500)
			return
		}

	} else if wallet.Money < 0 {

		err = a.app.CheckWallet(r.Context(), wallet.UserId)
		if err != nil {
			ErrHandler(w, err, 400)
			return
		}

		err = a.app.UpdateBalance(r.Context(), wallet.Money*-1, wallet.UserId)
		if err != nil {
			ErrHandler(w, err, 500)
			return
		}

		err = a.app.CreateTransaction(r.Context(), wallet.UserId, wallet.ReceiverId, wallet.Money*-1)
		if err != nil {
			ErrHandler(w, err, 500)
			return
		}

	}
	// перевод данных из бд
	ReturnBalance, err := a.app.CheckBalance(r.Context(), wallet.UserId)
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}

	err = json.NewEncoder(w).Encode(CheckBalanceResp{RetBalance: ReturnBalance})
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}

}
