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

	switch {

	case config.Server.Host == "":
		config.Server.Host = "0.0.0.0"

	case config.Server.Port.Http == "":
		config.Server.Port.Http = "8080"

	case config.Db.Driver == "":
		config.Db.Driver = "postgres"

	case config.Db.User == "":
		config.Db.User = "postgres"

	case config.Db.Password == "":
		config.Db.Password = "postgres"

	case config.Db.HostDb == "":
		config.Db.HostDb = "localhost"

	case config.Db.PortDb == "":
		config.Db.PortDb = "5432"

	case config.Db.Dbname == "":
		config.Db.Dbname = "postgres"

	case config.Db.Mode == "":
		config.Db.Mode = "disable"

	}

	return config, nil

}
