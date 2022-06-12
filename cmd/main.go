package main

import (
	"context"
	_ "database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

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

	configStruct, err := config.TakeConfigFromYaml(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = Run(configStruct)
	if err != nil {
		log.Fatal(err)
	}

}

//go:embed migrate/*.sql
var embedMigrations embed.FS

func Run(c *config.Config) error {

	db, err := sqlx.Open(c.Db.Driver, fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s", c.Db.User,
		c.Db.Password, c.Db.Dbname, c.Db.Mode, c.Db.HostDb, c.Db.PortDb))
	defer db.Close()
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrate"); err != nil {
		panic(err)
	}

	nachinka := repo.NewNachinka(db)
	defer nachinka.StopConnect()

	restAPiClient := rest_api.New(c.Client.APILayerAPIKey, c.Client.APILayerBasePath)

	a := app.NewApplication(nachinka, restAPiClient)

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
