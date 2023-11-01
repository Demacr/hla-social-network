package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Host       string `env:"HOST" envDefault:"localhost"`
	Port       int    `env:"PORT" envDefault:"8080"`
	JWTSecret  string `env:"JWT_SECRET,required"`
	MySQL      MySQLConfig
	PostgreSQL PostgreSQLConfig
	Redis      RedisConfig
}

type MySQLConfig struct {
	Host       string `env:"MYSQL_HOST"`
	Login      string `env:"MYSQL_USER"`
	Password   string `env:"MYSQL_PASSWORD"`
	Database   string `env:"MYSQL_DATABASE"`
	SlaveHosts string `env:"MYSQL_SLAVE_HOSTS"`
}

type PostgreSQLConfig struct {
	Host       string `env:"POSTGRES_HOST"`
	Login      string `env:"POSTGRES_USER"`
	Password   string `env:"POSTGRES_PASSWORD"`
	Database   string `env:"POSTGRES_DB"`
	SlaveHosts string `env:"POSTGRES_SLAVE_HOSTS"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"redis:6379"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	Database int    `env:"REDIS_DATABASE" envDefault:"0"`
}

func Configure() (*Config, error) {
	config := &Config{
		MySQL:      MySQLConfig{},
		PostgreSQL: PostgreSQLConfig{},
		Redis:      RedisConfig{},
	}
	if err := env.Parse(config); err != nil {
		return nil, errors.Wrap(err, "error during parsing env variables")
	}
	return config, nil
}
