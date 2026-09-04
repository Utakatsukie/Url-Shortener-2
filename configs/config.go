package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("Failed to load .env file: %w", err)
	}
	cfg := &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("TOKEN"),
		},
	}

	if cfg.Db.Dsn == "" {
		return nil, errors.New("REQUIRED_ENV_MISSING: DSN is empty")
	}
	if cfg.Auth.Secret == "" {
		return nil, errors.New("REQUIRED_ENV_MISSING: TOKEN is empty")
	}

	return cfg, nil
}
