package config

import (
	"database/sql"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"packs/internal/api"
	"packs/internal/app"
	"packs/internal/repo"
)

func TakeConfigFromYaml() (*http.Server, *sql.DB) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	flag.Parse()

	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config = &Config{}

	err = yaml.NewDecoder(file).Decode(config)

	db, err := sql.Open(config.Db.Driver, fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", config.Db.User,
		config.Db.Password, config.Db.Dbname, config.Db.Mode))

	if err != nil {
		panic(err)
	}

	nachinka := repo.NewNachinka(db)
	a := app.NewApplication(nachinka)

	server := &http.Server{
		Addr:    config.Server.Host + ":" + config.Server.Ports.Http,
		Handler: api.NewAPI(a),
	}

	return server, db

}
