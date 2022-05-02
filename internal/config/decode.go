package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

func TakeConfigFromYaml(s *string) (*Config, error) {

	file, err := os.Open(*s)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config = &Config{}

	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil

}
