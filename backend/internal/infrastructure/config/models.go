package config

import "time"

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

type AppConfig struct {
	Env string `mapstructure:"env"`
}

type HTTPConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

type AuthConfig struct {
	JWTSecret       string `mapstructure:"jwtSecret"`
	AccessTokenTTL  string `mapstructure:"accessTokenTTL"`
	RefreshTokenTTL string `mapstructure:"refreshTokenTTL"`
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0
	}
	return duration
}
