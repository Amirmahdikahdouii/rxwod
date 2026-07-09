package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	appathletedefaultsession "github.com/rxwod/backend/internal/application/athletedefaultsession"
)

type AthleteDefaultSessionHandler struct {
	service *appathletedefaultsession.Service
}

func NewAthleteDefaultSessionHandler(service *appathletedefaultsession.Service) *AthleteDefaultSessionHandler {
	return &AthleteDefaultSessionHandler{service: service}
}

func (h *AthleteDefaultSessionHandler) SetDefaultSession(c echo.Context) error {
	var req SetDefaultSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	result, err := h.service.SetDefaultSession(c.Request().Context(), appathletedefaultsession.SetDefaultSessionCommand{
		DayOfWeek: req.DayOfWeek,
		TimeSlot:  req.TimeSlot,
	})
	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusCreated, toDefaultSessionResponse(result))
}
