package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

type Config struct {
	Host  string `env:"HOST" envDefault:"localhost"`
	Port  int    `env:"PORT" envDefault:"8080"`
	MySQL struct {
		Host     string `env:"MYSQL_HOST"`
		Login    string `env:"MYSQL_LOGIN"`
		Password string `env:"MYSQL_PASSWORD"`
		Database string `env:"MYSQL_DATABASE"`
	}
}

func Configure() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, errors.Wrap(err, "error during parsing env variables")
	}
	return config, nil
}
