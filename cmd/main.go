package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

	"packs/internal/adapters/repo"
	"packs/internal/adapters/rest_api"
	"packs/internal/api"
	"packs/internal/app"
	"packs/internal/config"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("can`t loading env")
	}

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

	db, err := sqlx.Open(c.Db.Driver, fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s", c.Db.User,
		os.Getenv("PASSWORD_DB"), c.Db.Dbname, c.Db.Mode, c.Db.HostDb, c.Db.PortDb))
	if err != nil {
		return fmt.Errorf("sqlx.Open: %w", err)
	}
	defer db.Close()

	repo := repo.New(db)
	defer repo.StopConnect()

	restAPiClient := rest_api.New(os.Getenv("REST_API_API_KEY"), c.Client.APILayerBasePath)

	a := app.New(repo, restAPiClient)

	server := &http.Server{
		Addr:    c.Server.Host + ":" + c.Server.Port.Http,
		Handler: api.NewAPI(a),
	}

	go func() {
		err := server.ListenAndServe()
		if err == http.ErrServerClosed {
		} else {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), c.Server.Timeout.ServerTimeout*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
