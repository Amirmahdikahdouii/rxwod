package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rxwod/backend/internal/infrastructure/config"
)

func testRateLimitConfig() config.Config {
	return config.Config{
		RateLimit: config.RateLimitConfig{
			AuthPerIPRequestsPerMinute:         60,
			AuthPerIPBurst:                     2,
			AuthPerIdentifierRequestsPerMinute: 60,
			AuthPerIdentifierBurst:             2,
		},
	}
}

func TestAuthIPRateLimiterReturns429WhenExceeded(t *testing.T) {
	cfg := testRateLimitConfig()
	e := echo.New()
	e.POST("/auth/login", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}, authIPRateLimiter(cfg))

	for i := 0; i < cfg.AuthPerIPBurst(); i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"a@test.com","password":"secret"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.RemoteAddr = "203.0.113.10:1234"
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"a@test.com","password":"secret"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.RemoteAddr = "203.0.113.10:1234"
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}

	var body ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error != rateLimitExceededMessage {
		t.Fatalf("error = %q, want %q", body.Error, rateLimitExceededMessage)
	}
}

func TestAuthIdentifierRateLimiterLimitsByEmailAcrossIPs(t *testing.T) {
	cfg := testRateLimitConfig()
	e := echo.New()
	e.POST("/auth/login", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}, authIdentifierRateLimiter(cfg))

	payload := `{"email":"victim@example.com","password":"secret"}`
	remoteAddrs := []string{"203.0.113.10:1234", "203.0.113.11:1234", "203.0.113.12:1234"}

	for i := 0; i < cfg.AuthPerIdentifierBurst(); i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.RemoteAddr = remoteAddrs[i]
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.RemoteAddr = remoteAddrs[2]
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
}

func TestAuthIdentifierRateLimiterRestoresRequestBody(t *testing.T) {
	cfg := testRateLimitConfig()
	e := echo.New()
	e.POST("/auth/login", func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		}
		return c.JSON(http.StatusOK, map[string]string{"email": req.Email})
	}, authIdentifierRateLimiter(cfg))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"secret"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.RemoteAddr = "203.0.113.20:1234"
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["email"] != "user@example.com" {
		t.Fatalf("email = %q, want %q", body["email"], "user@example.com")
	}
}
