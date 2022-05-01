package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"packs/internal/api"
	"packs/internal/app"
	"packs/internal/repo"
)

var (
	conn   = "user=postgres password=postgres dbname=postgres sslmode=disable"
	decode = "Nekorrektnii zapros(decode)"
	encode = "Owibka pri formatirovanii "
)

func main() {

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	nachinka := repo.NewNachinka(db)
	a := app.NewApplication(nachinka)
	r := api.NewAPI(a)

	log.Fatal(http.ListenAndServe(":8080", r))

}

/*

func PrintBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")

	var balance Wallet

	err := json.NewDecoder(r.Body).Decode(&balance)
	if err != nil {
		http.Error(w, "Nekorrektnii zapros(decode)", 400)
		log.Println(err)
		return
	}

	row := CheckBalance(balance.UserId)
	err = row.Scan(&balance.Balance)
	if err != nil {
		http.Error(w, "Nekorrektnii zapros", 400) // уточнить описание и статус код ошибки, мб 500 ?
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(balance.Balance)
	if err != nil {
		http.Error(w, "Nekorrektnii zapros(encode)", 400)
		log.Println(err) // уточнить описание и статус код ошибки, мб 500 ?
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
/*
}
/*
func Transfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")

	var wallet Wallet

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		http.Error(w, "Nekorrektnii zapros", 400)
		log.Println(err)
		return
	}

	// уточнить нужна ли двойная проверка + логика чисел у "money"
	if wallet.Money < 0 || wallet.ReceiverId == wallet.UserId {
		http.Error(w, "Owibka v predostavlenix dannix", 400)
		return
	}

	err = CheckWallet(wallet.UserId)
	if err != nil {
		http.Error(w, "Y otpravitelya net kowelka", 400)
		return
	}

	err = CheckWallet(wallet.ReceiverId)
	if err != nil {
		http.Error(w, "Y polychatelya net kowelka", 400)
		return
	}

	err = TransferWithWallet(wallet.UserId, wallet.ReceiverId, wallet.Money)
	if err != nil {
		http.Error(w, "Problema s osywestvleniem tranzakcii", 400)
		return
	}

	row := CheckBalance(wallet.UserId)
	err = row.Scan(&wallet.Balance)
	if err != nil {
		http.Error(w, "Nekorrektnii zapros", 400)
		return
	}

	err = json.NewEncoder(w).Encode(wallet.Balance)
	if err != nil {
		log.Println(err)
		return
	}
}

func PrintTransactions(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application-json")

	var wallet Wallet

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		http.Error(w, "Nekorrektnii zapros", 400)
		log.Println(err)
	}

	rows, err := database.Query("select * from transaction where wallet_id = (select id from wallet where user_id = 100) order by money, date")
	if err != nil {
		http.Error(w, "Nekorrektnii zapros", 400)
		return
	}
	defer rows.Close()

	result := []Transaction{}

	for rows.Next() {
		r := Transaction{}
		err := rows.Scan(&r.TransactionId, &r.WalletId, &r.ReceiverWalletId, &r.Money, &r.Date)
		if err != nil {
			fmt.Println(err)
			continue // корректный код?
		}
		result = append(result, r)
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
		return
	}

}

func TransferWithWallet(userId int, receiverId int, money float64) error {
	tx, err := database.Begin()

	err = CreateTransaction(tx, userId, receiverId, money)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = UpdateBalance(tx, money, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = UpdateBalance(tx, money, receiverId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func DepositOrWithdrow(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application-json")
	var wallet Wallet
	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		http.Error(w, "Nekorrektnii zapros", 400)
		log.Println(err)
		return
	}

	ok := CheckWallet(wallet.UserId)
	if ok {
		http.Error(w, "Y otpravitelya net kowelka", 400)
		return
	}
	wallet.ReceiverId = wallet.UserId
	if wallet.Money == 0 {
		http.Error(w, "Nekorrektnii zapros", 400)
		log.Println(err)
		return
	} else if wallet.Money > 0 {
		tx, err := database.Begin()

		err = UpdateBalance(tx, wallet.Money*-1, wallet.UserId)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Problema v sql1", 500)
			return
		}

		err = CreateTransaction(tx, wallet.UserId, wallet.ReceiverId, wallet.Money)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Problema v sql", 500)
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return
		}

	} else if wallet.Money < 0 {
		tx, err := database.Begin()
		err = UpdateBalance(tx, wallet.Money*-1, wallet.UserId)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Problema v sql", 500)
			return
		}

		err = CreateTransaction(tx, wallet.UserId, wallet.ReceiverId, wallet.Money*-1)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Problema v sq", 500)
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return
		}
	}

	row := CheckBalance(wallet.UserId)

	err = row.Scan(&wallet.Balance)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.NewEncoder(w).Encode(wallet.Balance)
	if err != nil {
		log.Println(err)
		return
	}

}
func CheckWallet(i int) error {

	row := database.QueryRow("select user_id from wallet where User_Id = $1", i)
	err := row.Scan(&i)
	// уточнить про разделение проверок на ошибки в бд(когда такой таблицы нет например) и ошибку отсутствия кошелька,
	// когда мы используем полученное бул для его создания
	if err != nil {
		return err
	}
	return nil
}

func UpdateBalance(t *sql.Tx, f float64, i int) error {
	_, err := t.Exec("UPDATE wallet SET balance = balance - $1 WHERE User_Id = $2", f, i)
	return err
}

func CheckBalance(i int) *sql.Row {
	row := database.QueryRow("SELECT balance FROM wallet where user_id = $1", i)
	return row
}

func CreateTransaction(t *sql.Tx, i int, r int, f float64) error {
	_, err := t.Exec("insert into transaction (wallet_id, receiver_wallet_id, money, date) values ((select id from wallet where user_id = $1), (select id from wallet where user_id = $2), $3, $4)",
		i, r, f, time.Now())
	return err
}
*/
