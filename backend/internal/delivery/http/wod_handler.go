package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	appwod "github.com/rxwod/backend/internal/application/wod"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
	"github.com/rxwod/backend/internal/infrastructure/postgres"
)

type WODHandler struct {
	service *appwod.Service
}

func NewWODHandler(service *appwod.Service) *WODHandler {
	return &WODHandler{service: service}
}

func (h *WODHandler) Create(c echo.Context) error {
	var req CreateWODRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	cmd, err := toCreateCommand(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	result, err := h.service.Create(c.Request().Context(), cmd)
	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusCreated, toCreateResponse(result))
}

func (h *WODHandler) List(c echo.Context) error {
	items, err := h.service.List(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}

	responses := make([]WODSummaryResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toSummaryResponse(item))
	}
	return c.JSON(http.StatusOK, responses)
}

func (h *WODHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	item, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toDetailResponse(item))
}

func toCreateCommand(req CreateWODRequest) (appwod.CreateCommand, error) {
	movements := make([]appwod.MovementInput, 0, len(req.Movements))
	for _, movement := range req.Movements {
		movements = append(movements, appwod.MovementInput{
			Position:  movement.Position,
			Name:      movement.Name,
			Reps:      movement.Reps,
			LoadValue: movement.LoadValue,
			LoadUnit:  movement.LoadUnit,
			Notes:     movement.Notes,
		})
	}

	switch domainwod.WODType(req.Type) {
	case domainwod.WODTypeAMRAP:
		if req.Config.TimeCapSeconds == nil {
			return appwod.CreateCommand{}, fmt.Errorf("timeCapSeconds is required for AMRAP")
		}
		return appwod.CreateCommand{
			Type: domainwod.WODTypeAMRAP,
			AMRAP: &appwod.CreateAMRAPCommand{
				Name:        req.Name,
				Description: req.Description,
				TimeCap:     *req.Config.TimeCapSeconds,
				Movements:   movements,
			},
		}, nil
	case domainwod.WODTypeForTime:
		if req.Config.Rounds == nil {
			return appwod.CreateCommand{}, fmt.Errorf("rounds is required for FORTIME")
		}
		return appwod.CreateCommand{
			Type: domainwod.WODTypeForTime,
			ForTime: &appwod.CreateForTimeCommand{
				Name:        req.Name,
				Description: req.Description,
				Rounds:      *req.Config.Rounds,
				TimeCap:     req.Config.TimeCapSeconds,
				Movements:   movements,
			},
		}, nil
	case domainwod.WODTypeTabata:
		if req.Config.WorkSeconds == nil || req.Config.RestSeconds == nil || req.Config.Rounds == nil || req.Config.Cycles == nil {
			return appwod.CreateCommand{}, fmt.Errorf("workSeconds, restSeconds, rounds, and cycles are required for TABATA")
		}
		return appwod.CreateCommand{
			Type: domainwod.WODTypeTabata,
			Tabata: &appwod.CreateTabataCommand{
				Name:        req.Name,
				Description: req.Description,
				WorkSeconds: *req.Config.WorkSeconds,
				RestSeconds: *req.Config.RestSeconds,
				Rounds:      *req.Config.Rounds,
				Cycles:      *req.Config.Cycles,
				Movements:   movements,
			},
		}, nil
	case domainwod.WODTypeEMOM:
		if req.Config.IntervalSeconds == nil || req.Config.Rounds == nil {
			return appwod.CreateCommand{}, fmt.Errorf("intervalSeconds and rounds are required for EMOM")
		}
		return appwod.CreateCommand{
			Type: domainwod.WODTypeEMOM,
			EMOM: &appwod.CreateEMOMCommand{
				Name:            req.Name,
				Description:     req.Description,
				IntervalSeconds: *req.Config.IntervalSeconds,
				Rounds:          *req.Config.Rounds,
				Movements:       movements,
			},
		}, nil
	default:
		return appwod.CreateCommand{}, domainwod.ErrUnknownWODType
	}
}

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domainwod.ErrInvalidName),
		errors.Is(err, domainwod.ErrInvalidTimeCap),
		errors.Is(err, domainwod.ErrInvalidRounds),
		errors.Is(err, domainwod.ErrInvalidWorkSeconds),
		errors.Is(err, domainwod.ErrInvalidRestSeconds),
		errors.Is(err, domainwod.ErrInvalidInterval),
		errors.Is(err, domainwod.ErrInvalidCycles),
		errors.Is(err, domainwod.ErrMovementRequired),
		errors.Is(err, domainwod.ErrInvalidMovement),
		errors.Is(err, domainwod.ErrInvalidLoadUnit),
		errors.Is(err, domainwod.ErrInvalidReps),
		errors.Is(err, domainwod.ErrInvalidPosition),
		errors.Is(err, domainwod.ErrUnknownWODType):
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, postgres.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "wod not found"})
	default:
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
