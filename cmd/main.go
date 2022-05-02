package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"packs/internal/api"
	"packs/internal/app"
	"packs/internal/config"
	"packs/internal/repo"
)

func main() {

	var configPath string

	flag.StringVar(&configPath, "config", "/config.yml", "path to config file")

	flag.Parse()

	configStruct, err := config.TakeConfigFromYaml(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = Run(configStruct)
	if err != nil {
		log.Fatal(err)
	}

}

func Run(c *config.Config) error {

	db, err := sql.Open(c.Db.Driver, fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", c.Db.User,
		c.Db.Password, c.Db.Dbname, c.Db.Mode))
	defer db.Close()
	if err != nil {
		panic(err)
	}

	nachinka := repo.NewNachinka(db)
	a := app.NewApplication(nachinka)

	server := &http.Server{
		Addr:    c.Server.Host + ":" + c.Server.Port.Http,
		Handler: api.NewAPI(a),
	}

	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
	} else if err != nil {
		return err
	}

	return nil

}
