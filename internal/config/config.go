package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Host  string `env:"HOST" envDefault:"localhost"`
	Port  int    `env:"PORT" envDefault:"8080"`
	MySQL MySQLConfig
}

type MySQLConfig struct {
	Host     string `env:"MYSQL_HOST"`
	Login    string `env:"MYSQL_USER"`
	Password string `env:"MYSQL_PASSWORD"`
	Database string `env:"MYSQL_DATABASE"`
}

func Configure() (*Config, error) {
	config := &Config{MySQL: MySQLConfig{}}
	if err := env.Parse(config); err != nil {
		return nil, errors.Wrap(err, "error during parsing env variables")
	}
	return config, nil
}
