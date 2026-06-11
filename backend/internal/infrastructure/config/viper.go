package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv      string
	HTTPPort    int
	DatabaseURL string
}

func Load() (Config, error) {
	v := viper.New()
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("HTTP_PORT", 8080)
	v.SetDefault("DATABASE_URL", "postgres://rxwod:rxwod@localhost:5432/rxwod?sslmode=disable")

	v.AutomaticEnv()

	cfg := Config{
		AppEnv:      v.GetString("APP_ENV"),
		HTTPPort:    v.GetInt("HTTP_PORT"),
		DatabaseURL: v.GetString("DATABASE_URL"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}
