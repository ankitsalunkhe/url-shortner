package config

import (
	"fmt"
	"log/slog"

	"github.com/ankitsalunkhe/url-shortner/internal/db"
	"github.com/ankitsalunkhe/url-shortner/internal/retriever"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        int    `envconfig:"PORT" required:"true"`
	BasePath    string `envconfig:"BASE_PATH" required:"true"`
	Environment string `envconfig:"ENVIRONMENT" default:"local"`
	RtConfig    retriever.Config
	DBConfig    db.Config
}

func New() (Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Debug("Failed to Load Environment Variables", "error", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("unable to parse config variables: %w", err)
	}

	return cfg, nil
}
