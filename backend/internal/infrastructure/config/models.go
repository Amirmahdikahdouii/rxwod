package config

import (
	"time"

	"golang.org/x/time/rate"
)

type Config struct {
	App       AppConfig       `mapstructure:"app"`
	HTTP      HTTPConfig      `mapstructure:"http"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	RateLimit RateLimitConfig `mapstructure:"ratelimit"`
}

type RateLimitConfig struct {
	AuthPerIPRequestsPerMinute         float64 `mapstructure:"authPerIPRequestsPerMinute"`
	AuthPerIPBurst                     int     `mapstructure:"authPerIPBurst"`
	AuthPerIdentifierRequestsPerMinute float64 `mapstructure:"authPerIdentifierRequestsPerMinute"`
	AuthPerIdentifierBurst             int     `mapstructure:"authPerIdentifierBurst"`
}

type AppConfig struct {
	Env            string `mapstructure:"env"`
	FrontendURL    string `mapstructure:"frontendURL"`
	AllowedOrigins []string `mapstructure:"allowedOrigins"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

type HTTPConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

type AuthConfig struct {
	JWTSecret            string `mapstructure:"jwtSecret"`
	AccessTokenTTL       string `mapstructure:"accessTokenTTL"`
	RefreshTokenTTL      string `mapstructure:"refreshTokenTTL"`
	PasswordResetTTL     string `mapstructure:"passwordResetTTL"`
	EmailVerificationTTL string `mapstructure:"emailVerificationTTL"`
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0
	}
	return duration
}

func requestsPerMinuteToLimit(requestsPerMinute float64) rate.Limit {
	return rate.Limit(requestsPerMinute / 60.0)
}

func (c Config) AuthPerIPRateLimit() rate.Limit {
	return requestsPerMinuteToLimit(c.RateLimit.AuthPerIPRequestsPerMinute)
}

func (c Config) AuthPerIPBurst() int {
	return c.RateLimit.AuthPerIPBurst
}

func (c Config) AuthPerIdentifierRateLimit() rate.Limit {
	return requestsPerMinuteToLimit(c.RateLimit.AuthPerIdentifierRequestsPerMinute)
}

func (c Config) AuthPerIdentifierBurst() int {
	return c.RateLimit.AuthPerIdentifierBurst
}
