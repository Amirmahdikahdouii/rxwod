package http

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type dbPinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	db dbPinger
}

func NewHealthHandler(db dbPinger) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, ErrorResponse{Error: "database unavailable"})
	}
	return c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
