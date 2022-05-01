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
