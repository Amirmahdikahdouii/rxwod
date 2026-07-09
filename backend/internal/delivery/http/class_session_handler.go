package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	appclasssession "github.com/rxwod/backend/internal/application/classsession"
)

type ClassSessionHandler struct {
	service *appclasssession.Service
}

func NewClassSessionHandler(service *appclasssession.Service) *ClassSessionHandler {
	return &ClassSessionHandler{service: service}
}

func (h *ClassSessionHandler) Create(c echo.Context) error {
	var req CreateClassSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	startTime, err := parseRFC3339(req.StartTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "startTime must be RFC3339"})
	}
	endTime, err := parseRFC3339(req.EndTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "endTime must be RFC3339"})
	}

	result, err := h.service.Create(c.Request().Context(), appclasssession.CreateClassSessionCommand{
		WodID:     req.WodID,
		StartTime: startTime,
		EndTime:   endTime,
		Capacity:  req.Capacity,
	})
	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusCreated, toCreateClassSessionResponse(result))
}

func (h *ClassSessionHandler) List(c echo.Context) error {
	fromRaw := c.QueryParam("from")
	toRaw := c.QueryParam("to")
	if fromRaw == "" || toRaw == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "from and to query params are required"})
	}

	from, err := parseRFC3339(fromRaw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "from must be RFC3339"})
	}
	to, err := parseRFC3339(toRaw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "to must be RFC3339"})
	}

	sessions, err := h.service.List(c.Request().Context(), appclasssession.ListClassSessionsCommand{
		From: from,
		To:   to,
	})
	if err != nil {
		return mapError(c, err)
	}

	responses := make([]ClassSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, toClassSessionResponse(session))
	}
	return c.JSON(http.StatusOK, responses)
}

func parseRFC3339(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}
