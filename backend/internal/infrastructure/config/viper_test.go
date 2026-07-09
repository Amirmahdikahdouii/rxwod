package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadUsesDefaultsWhenConfigFileIsMissing(t *testing.T) {
	clearConfigEnv(t)

	cfg, err := load(t.TempDir())
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if got, want := cfg.AppEnv(), "development"; got != want {
		t.Fatalf("AppEnv() = %q, want %q", got, want)
	}

	if got, want := cfg.HTTPPort(), 8080; got != want {
		t.Fatalf("HTTPPort() = %d, want %d", got, want)
	}

	if got, want := cfg.DatabaseURL(), "postgres://rxwod:rxwod@localhost:5432/rxwod?sslmode=disable"; got != want {
		t.Fatalf("DatabaseURL() = %q, want %q", got, want)
	}

	if got, want := cfg.LogLevel(), "info"; got != want {
		t.Fatalf("LogLevel() = %q, want %q", got, want)
	}

	if got, want := cfg.AllowedOrigins(), []string{"http://localhost:5173"}; len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("AllowedOrigins() = %v, want %v", got, want)
	}

	if got, want := cfg.AuthPerIPBurst(), 10; got != want {
		t.Fatalf("AuthPerIPBurst() = %d, want %d", got, want)
	}

	if got, want := cfg.AuthPerIdentifierBurst(), 3; got != want {
		t.Fatalf("AuthPerIdentifierBurst() = %d, want %d", got, want)
	}
}

func TestLoadReadsYAMLConfigFile(t *testing.T) {
	clearConfigEnv(t)
	dir := t.TempDir()
	writeConfig(t, dir, `
app:
  env: test
http:
  port: 9090
database:
  url: "postgres://yaml:yaml@localhost:5432/yaml?sslmode=disable"
`)

	cfg, err := load(dir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if got, want := cfg.AppEnv(), "test"; got != want {
		t.Fatalf("AppEnv() = %q, want %q", got, want)
	}

	if got, want := cfg.HTTPPort(), 9090; got != want {
		t.Fatalf("HTTPPort() = %d, want %d", got, want)
	}

	if got, want := cfg.DatabaseURL(), "postgres://yaml:yaml@localhost:5432/yaml?sslmode=disable"; got != want {
		t.Fatalf("DatabaseURL() = %q, want %q", got, want)
	}
}

func TestLoadLetsEnvironmentOverrideYAMLConfig(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, `
app:
  env: yaml
http:
  port: 9090
database:
  url: "postgres://yaml:yaml@localhost:5432/yaml?sslmode=disable"
`)
	t.Setenv("APP_ENV", "production")
	t.Setenv("HTTP_PORT", "7070")
	t.Setenv("DATABASE_URL", "postgres://env:env@localhost:5432/env?sslmode=disable")
	t.Setenv("LOG_LEVEL", "error")

	cfg, err := load(dir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if got, want := cfg.AppEnv(), "production"; got != want {
		t.Fatalf("AppEnv() = %q, want %q", got, want)
	}

	if got, want := cfg.HTTPPort(), 7070; got != want {
		t.Fatalf("HTTPPort() = %d, want %d", got, want)
	}

	if got, want := cfg.DatabaseURL(), "postgres://env:env@localhost:5432/env?sslmode=disable"; got != want {
		t.Fatalf("DatabaseURL() = %q, want %q", got, want)
	}

	if got, want := cfg.LogLevel(), "error"; got != want {
		t.Fatalf("LogLevel() = %q, want %q", got, want)
	}
}

func TestLoadRejectsEmptyDatabaseURL(t *testing.T) {
	clearConfigEnv(t)
	dir := t.TempDir()
	writeConfig(t, dir, `
app:
  env: test
http:
  port: 9090
database:
  url: ""
`)

	_, err := load(dir)
	if err == nil {
		t.Fatal("load config succeeded, want error")
	}

	if !strings.Contains(err.Error(), "database.url") {
		t.Fatalf("error = %q, want database.url validation", err.Error())
	}
}

func TestLoadRejectsInvalidHTTPPort(t *testing.T) {
	clearConfigEnv(t)
	dir := t.TempDir()
	writeConfig(t, dir, `
app:
  env: test
http:
  port: 70000
database:
  url: "postgres://yaml:yaml@localhost:5432/yaml?sslmode=disable"
`)

	_, err := load(dir)
	if err == nil {
		t.Fatal("load config succeeded, want error")
	}

	if !strings.Contains(err.Error(), "http.port") {
		t.Fatalf("error = %q, want http.port validation", err.Error())
	}
}

func TestLoadRejectsInvalidLogLevel(t *testing.T) {
	clearConfigEnv(t)
	dir := t.TempDir()
	writeConfig(t, dir, `
app:
  env: test
http:
  port: 9090
database:
  url: "postgres://yaml:yaml@localhost:5432/yaml?sslmode=disable"
logging:
  level: verbose
`)

	_, err := load(dir)
	if err == nil {
		t.Fatal("load config succeeded, want error")
	}

	if !strings.Contains(err.Error(), "logging.level") {
		t.Fatalf("error = %q, want logging.level validation", err.Error())
	}
}

func TestLoadLetsEnvironmentOverrideRateLimitConfig(t *testing.T) {
	clearConfigEnv(t)
	dir := t.TempDir()
	writeConfig(t, dir, `
app:
  env: test
http:
  port: 9090
database:
  url: "postgres://yaml:yaml@localhost:5432/yaml?sslmode=disable"
ratelimit:
  authPerIPRequestsPerMinute: 20
  authPerIPBurst: 10
  authPerIdentifierRequestsPerMinute: 5
  authPerIdentifierBurst: 3
`)
	t.Setenv("RATELIMIT_AUTH_PER_IP_RPM", "30")
	t.Setenv("RATELIMIT_AUTH_PER_IP_BURST", "15")
	t.Setenv("RATELIMIT_AUTH_PER_IDENTIFIER_RPM", "8")
	t.Setenv("RATELIMIT_AUTH_PER_IDENTIFIER_BURST", "4")

	cfg, err := load(dir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if got, want := cfg.RateLimit.AuthPerIPRequestsPerMinute, 30.0; got != want {
		t.Fatalf("AuthPerIPRequestsPerMinute = %v, want %v", got, want)
	}

	if got, want := cfg.AuthPerIPBurst(), 15; got != want {
		t.Fatalf("AuthPerIPBurst() = %d, want %d", got, want)
	}

	if got, want := cfg.RateLimit.AuthPerIdentifierRequestsPerMinute, 8.0; got != want {
		t.Fatalf("AuthPerIdentifierRequestsPerMinute = %v, want %v", got, want)
	}

	if got, want := cfg.AuthPerIdentifierBurst(), 4; got != want {
		t.Fatalf("AuthPerIdentifierBurst() = %d, want %d", got, want)
	}
}

func TestLoadRejectsInvalidRateLimitConfig(t *testing.T) {
	clearConfigEnv(t)
	dir := t.TempDir()
	writeConfig(t, dir, `
app:
  env: test
http:
  port: 9090
database:
  url: "postgres://yaml:yaml@localhost:5432/yaml?sslmode=disable"
ratelimit:
  authPerIPRequestsPerMinute: 0
  authPerIPBurst: 10
  authPerIdentifierRequestsPerMinute: 5
  authPerIdentifierBurst: 3
`)

	_, err := load(dir)
	if err == nil {
		t.Fatal("load config succeeded, want error")
	}

	if !strings.Contains(err.Error(), "ratelimit.authPerIPRequestsPerMinute") {
		t.Fatalf("error = %q, want ratelimit.authPerIPRequestsPerMinute validation", err.Error())
	}
}

func clearConfigEnv(t *testing.T) {
	t.Helper()
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_ALLOWED_ORIGINS", "")
	t.Setenv("HTTP_PORT", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("AUTH_JWT_SECRET", "")
	t.Setenv("AUTH_ACCESS_TOKEN_TTL", "")
	t.Setenv("AUTH_REFRESH_TOKEN_TTL", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("RATELIMIT_AUTH_PER_IP_RPM", "")
	t.Setenv("RATELIMIT_AUTH_PER_IP_BURST", "")
	t.Setenv("RATELIMIT_AUTH_PER_IDENTIFIER_RPM", "")
	t.Setenv("RATELIMIT_AUTH_PER_IDENTIFIER_BURST", "")
}

func writeConfig(t *testing.T, dir string, contents string) {
	t.Helper()

	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(strings.TrimSpace(contents)+"\n"), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}
}
