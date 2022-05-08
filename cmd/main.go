package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"packs/internal/api"
	"packs/internal/app"
	"packs/internal/config"
	"packs/internal/repo"
	"time"
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

	db, err := sql.Open(c.Db.Driver, fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s", c.Db.User,
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

	// need err?
	defer nachinka.StopConnect()

	a := app.NewApplication(nachinka)

	server := &http.Server{
		Addr:    c.Server.Host + ":" + c.Server.Port.Http,
		Handler: api.NewAPI(a),
	}

	go func() {
		err = server.ListenAndServe()
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
