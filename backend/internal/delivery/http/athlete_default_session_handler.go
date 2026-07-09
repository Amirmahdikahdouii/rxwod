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

func (h *AthleteDefaultSessionHandler) ListMine(c echo.Context) error {
	sessions, err := h.service.ListMyDefaultSessions(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}

	responses := make([]DefaultSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, toDefaultSessionResponse(session))
	}
	return c.JSON(http.StatusOK, responses)
}

func (h *AthleteDefaultSessionHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	if err := h.service.RemoveDefaultSession(c.Request().Context(), id); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}
