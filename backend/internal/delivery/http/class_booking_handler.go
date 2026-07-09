package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	appclassbooking "github.com/rxwod/backend/internal/application/classbooking"
)

type ClassBookingHandler struct {
	service *appclassbooking.Service
}

func NewClassBookingHandler(service *appclassbooking.Service) *ClassBookingHandler {
	return &ClassBookingHandler{service: service}
}

func (h *ClassBookingHandler) Book(c echo.Context) error {
	sessionID := c.Param("id")
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	result, err := h.service.Book(c.Request().Context(), appclassbooking.BookCommand{
		SessionID: sessionID,
	})
	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusCreated, toClassBookingResponse(result))
}

func (h *ClassBookingHandler) Overbook(c echo.Context) error {
	sessionID := c.Param("id")
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	var req OverbookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	if req.AthleteUserID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "athleteUserId is required"})
	}

	result, err := h.service.Overbook(c.Request().Context(), appclassbooking.OverbookCommand{
		SessionID:     sessionID,
		AthleteUserID: req.AthleteUserID,
	})
	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusCreated, toClassBookingResponse(result))
}

func (h *ClassBookingHandler) Cancel(c echo.Context) error {
	sessionID := c.Param("id")
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	var req CancelBookingRequest
	_ = c.Bind(&req)

	result, err := h.service.Cancel(c.Request().Context(), appclassbooking.CancelCommand{
		SessionID:     sessionID,
		AthleteUserID: req.AthleteUserID,
	})
	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusOK, toClassBookingResponse(result))
}
