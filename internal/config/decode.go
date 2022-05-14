package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

func TakeConfigFromYaml(s string) (*Config, error) {
	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config = &Config{}

	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}

	if config.Server.Port.Http == "" {
		config.Server.Port.Http = "8080"
	}

	if config.Server.Timeout.ServerTimeout == 0 {
		config.Server.Timeout.ServerTimeout = 30
	}

	if config.Db.Driver == "" {
		config.Db.Driver = "postgres"
	}

	if config.Db.User == "" {
		config.Db.User = "postgres"
	}

	if config.Db.Password == "" {
		config.Db.Password = "postgres"
	}

	if config.Db.HostDb == "" {
		config.Db.HostDb = "localhost"
	}

	if config.Db.PortDb == "" {
		config.Db.PortDb = "5432"
	}

	if config.Db.Dbname == "" {
		config.Db.Dbname = "postgres"
	}

	if config.Db.Mode == "" {
		config.Db.Mode = "disable"
	}

	return config, nil
}
