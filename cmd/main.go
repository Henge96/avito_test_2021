package main

import (
	_ "github.com/lib/pq"
	"log"
	"packs/internal/config"
)

func main() {

	configServer, db := config.TakeConfigFromYaml()
	defer db.Close()

	log.Fatal(configServer.ListenAndServe())

}
