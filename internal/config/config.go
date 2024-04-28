package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Port          string `env:"PORT" envDefault:"9000"`
	DriverName    string `env:"DRIVER_NAME" envDefault:"postgres"`
	DbConnect     string `env:"DB_CONNECT" envDefault:"host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"`
	SecretKey     string `env:"SECRET_KEY" envDefault:"hash"`
	TokenLifetime int    `env:"TOKEN_LIFETIME" envDefault:"2"`
	HashCost      int    `env:"HASH_COST" envDefault:"6"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
