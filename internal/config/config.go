package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	APIConfig    APIConfig
}

type APIConfig struct {
	GetAgeAPI     string `env:"GET_AGE_API"`
	GetGenderAPI  string `env:"GET_GENDER_API"`
	GetCountryAPI string `env:"GET_COUNTRY_API"`
}

type ServerConfig struct {
	ServerMode string `env:"ENVIRONMENT"`
	ServerPort string `env:"HTTP_PORT"`
}

type DBConfig struct {
	PgUser     string `env:"PGUSER"`
	PgPassword string `env:"PGPASSWORD"`
	PgHost     string `env:"PGHOST"`
	PgPort     uint16 `env:"PGPORT"`
	PgDatabase string `env:"PGDATABASE"`
	PgSSLMode  string `env:"PGSSLMODE"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from environment variables: %w", err)
	}

	return cfg, nil
}
