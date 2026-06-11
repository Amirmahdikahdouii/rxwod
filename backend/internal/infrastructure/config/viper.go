package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Load() (Config, error) {
	return load(defaultConfigPaths()...)
}

func load(configPaths ...string) (Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	mustBindEnv(v, "app.env", "APP_ENV")
	mustBindEnv(v, "http.port", "HTTP_PORT")
	mustBindEnv(v, "database.url", "DATABASE_URL")

	v.SetDefault("app.env", "development")
	v.SetDefault("http.port", 8080)
	v.SetDefault("database.url", "postgres://rxwod:rxwod@localhost:5432/rxwod?sslmode=disable")

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return Config{}, fmt.Errorf("read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("bind config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func defaultConfigPaths() []string {
	return []string{
		".",
		"./backend",
		"../..",
		"../../..",
	}
}

func mustBindEnv(v *viper.Viper, key string, envVars ...string) {
	keys := append([]string{key}, envVars...)
	if err := v.BindEnv(keys...); err != nil {
		panic(fmt.Sprintf("bind env %q: %v", key, err))
	}
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Database.URL) == "" {
		return fmt.Errorf("database.url is required")
	}

	if c.HTTP.Port < 1 || c.HTTP.Port > 65535 {
		return fmt.Errorf("http.port must be between 1 and 65535")
	}

	return nil
}

func (c Config) AppEnv() string {
	return c.App.Env
}

func (c Config) HTTPPort() int {
	return c.HTTP.Port
}

func (c Config) DatabaseURL() string {
	return c.Database.URL
}
