package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"

	"packs/internal/app"
)

// MwHandler1 for an example of use
func MwHandler1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")
		log.Println(r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// MwHandler2 for an example of use
func MwHandler2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")
		log.Println(r.Method)
		next.ServeHTTP(w, r)
	})
}

// ErrHandler response.
func ErrHandler(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	log.Println(err)
}

// GetBalance get user balance.
func (a *Api) GetBalance(w http.ResponseWriter, r *http.Request) {
	var wallet app.RequestBalance

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		ErrHandler(w, err, http.StatusBadRequest)
		return
	}

	if wallet.Currency == "" {
		wallet.Currency = "RUB"
	}

	err = validateStruct(wallet)
	if err != nil {
		ErrHandler(w, app.ErrInvalidArgument, http.StatusBadRequest)
		return
	}

	returnBalance, err := a.app.GetUserBalance(r.Context(), wallet.UserID, strings.ToUpper(wallet.Currency))
	if err != nil {
		switch {
		case errors.Is(err, app.ErrNotFound):
			ErrHandler(w, err, http.StatusNotFound)
			return
		default:
			ErrHandler(w, err, http.StatusInternalServerError)
			return
		}
	}

	err = json.NewEncoder(w).Encode(app.Wallet{
		ID:      returnBalance.ID,
		UserID:  returnBalance.UserID,
		Balance: returnBalance.Balance,
	})
	if err != nil {
		ErrHandler(w, err, http.StatusInternalServerError)
		return
	}
}

// TransferBetweenWallets between wallets.
func (a *Api) TransferBetweenWallets(w http.ResponseWriter, r *http.Request) {

	var transfer app.Transaction

	err := json.NewDecoder(r.Body).Decode(&transfer)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	err = validateStruct(transfer)
	if err != nil {
		ErrHandler(w, app.ErrInvalidArgument, http.StatusBadRequest)
		return
	}

	if !transfer.Amount.GreaterThan(decimal.NewFromInt(0)) {
		ErrHandler(w, app.ErrInvalidArgument, http.StatusBadRequest)
		return
	}

	//TODO: add logic whem we got err not found
	id, err := a.app.Transfer(r.Context(), transfer)
	if err != nil {
		ErrHandler(w, err, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(ResponseTransactionID{ID: id})
	if err != nil {
		ErrHandler(w, err, http.StatusInternalServerError)
		return
	}
}

// GetTransactions getting transactions by params.
func (a *Api) GetTransactions(w http.ResponseWriter, r *http.Request) {

	var params app.UserTransactionsParam

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	err = validateStruct(params)
	if err != nil {
		ErrHandler(w, app.ErrInvalidArgument, 400)
		return
	}

	res, total, err := a.app.GetUserTransactions(r.Context(), params)
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}

	err = json.NewEncoder(w).Encode(ResponseTransactions{Response: res, Total: total})
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}
}

// ChangeUserBalance balance.
func (a *Api) ChangeUserBalance(w http.ResponseWriter, r *http.Request) {

	var wallet app.ChangeBalance

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	err = validateStruct(wallet)
	if err != nil {
		ErrHandler(w, err, 400)
		return
	}

	transaction, err := a.app.ChangeBalance(r.Context(), wallet)
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}

	err = json.NewEncoder(w).Encode(ResponseBalance{Balance: transaction.Balance.StringFixedBank(2)})
	if err != nil {
		ErrHandler(w, err, 500)
		return
	}
}

//TODO: add another validate logic.
func validateStruct(object interface{}) error {
	validate := validator.New()
	err := validate.Struct(object)
	if err != nil {
		return fmt.Errorf("validate.Struct: %w", err)
	}

	return nil
}
