package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rxwod/backend/internal/infrastructure/config"
	"golang.org/x/time/rate"
)

const rateLimitExceededMessage = "rate limit exceeded"

type emailIdentifierRequest struct {
	Email string `json:"email"`
}

func authIPRateLimiter(cfg config.Config) echo.MiddlewareFunc {
	return newAuthRateLimiter(cfg.AuthPerIPRateLimit(), cfg.AuthPerIPBurst(), func(c echo.Context) (string, error) {
		return c.RealIP(), nil
	})
}

func authIdentifierRateLimiter(cfg config.Config) echo.MiddlewareFunc {
	return newAuthRateLimiter(cfg.AuthPerIdentifierRateLimit(), cfg.AuthPerIdentifierBurst(), extractAuthEmailIdentifier)
}

func newAuthRateLimiter(limit rate.Limit, burst int, identifierExtractor middleware.Extractor) echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate:      limit,
			Burst:     burst,
			ExpiresIn: 5 * time.Minute,
		}),
		IdentifierExtractor: identifierExtractor,
		DenyHandler:         rateLimitDenyHandler,
		ErrorHandler:        rateLimitErrorHandler,
	})
}

func rateLimitDenyHandler(c echo.Context, _ string, _ error) error {
	return c.JSON(http.StatusTooManyRequests, ErrorResponse{Error: rateLimitExceededMessage})
}

func rateLimitErrorHandler(c echo.Context, _ error) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
}

func extractAuthEmailIdentifier(c echo.Context) (string, error) {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.RealIP(), nil
	}
	c.Request().Body = io.NopCloser(bytes.NewReader(body))

	if len(body) == 0 {
		return c.RealIP(), nil
	}

	var req emailIdentifierRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return c.RealIP(), nil
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return c.RealIP(), nil
	}

	return email, nil
}
