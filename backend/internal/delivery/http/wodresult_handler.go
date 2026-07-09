package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	appwodresult "github.com/rxwod/backend/internal/application/wodresult"
)

type WODResultHandler struct {
	service *appwodresult.Service
}

func NewWODResultHandler(service *appwodresult.Service) *WODResultHandler {
	return &WODResultHandler{service: service}
}

func (h *WODResultHandler) SubmitResult(c echo.Context) error {
	wodID := c.Param("id")
	if wodID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	var req SubmitResultRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	result, err := h.service.SubmitResult(c.Request().Context(), appwodresult.SubmitResultCommand{
		WODID:      wodID,
		ScoreValue: req.ScoreValue,
		IsRx:       req.IsRx,
		Notes:      req.Notes,
	})
	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusCreated, toResultResponse(result))
}

func (h *WODResultHandler) GetLeaderboard(c echo.Context) error {
	wodID := c.Param("id")
	if wodID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	leaderboard, err := h.service.GetLeaderboard(c.Request().Context(), wodID)
	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusOK, toLeaderboardResponse(leaderboard))
}
