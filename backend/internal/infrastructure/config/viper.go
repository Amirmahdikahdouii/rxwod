package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
	mustBindEnv(v, "app.frontendURL", "APP_FRONTEND_URL")
	mustBindEnv(v, "http.port", "HTTP_PORT")
	mustBindEnv(v, "database.url", "DATABASE_URL")
	mustBindEnv(v, "auth.jwtSecret", "AUTH_JWT_SECRET")
	mustBindEnv(v, "auth.accessTokenTTL", "AUTH_ACCESS_TOKEN_TTL")
	mustBindEnv(v, "auth.refreshTokenTTL", "AUTH_REFRESH_TOKEN_TTL")
	mustBindEnv(v, "auth.passwordResetTTL", "AUTH_PASSWORD_RESET_TTL")

	v.SetDefault("app.env", "development")
	v.SetDefault("app.frontendURL", "http://localhost:5173")
	v.SetDefault("http.port", 8080)
	v.SetDefault("database.url", "postgres://rxwod:rxwod@localhost:5432/rxwod?sslmode=disable")
	v.SetDefault("auth.jwtSecret", "dev-secret-change-me")
	v.SetDefault("auth.accessTokenTTL", "15m")
	v.SetDefault("auth.refreshTokenTTL", "168h")
	v.SetDefault("auth.passwordResetTTL", "1h")

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

	if strings.TrimSpace(c.Auth.JWTSecret) == "" {
		return fmt.Errorf("auth.jwtSecret is required")
	}

	if c.AccessTokenTTL() <= 0 {
		return fmt.Errorf("auth.accessTokenTTL must be a valid positive duration")
	}

	if c.RefreshTokenTTL() <= 0 {
		return fmt.Errorf("auth.refreshTokenTTL must be a valid positive duration")
	}

	if c.PasswordResetTTL() <= 0 {
		return fmt.Errorf("auth.passwordResetTTL must be a valid positive duration")
	}

	return nil
}

func (c Config) AppEnv() string {
	return c.App.Env
}

func (c Config) FrontendURL() string {
	return c.App.FrontendURL
}

func (c Config) HTTPPort() int {
	return c.HTTP.Port
}

func (c Config) DatabaseURL() string {
	return c.Database.URL
}

func (c Config) JWTSecret() string {
	return c.Auth.JWTSecret
}

func (c Config) AccessTokenTTL() time.Duration {
	return parseDuration(c.Auth.AccessTokenTTL)
}

func (c Config) RefreshTokenTTL() time.Duration {
	return parseDuration(c.Auth.RefreshTokenTTL)
}

func (c Config) PasswordResetTTL() time.Duration {
	return parseDuration(c.Auth.PasswordResetTTL)
}
