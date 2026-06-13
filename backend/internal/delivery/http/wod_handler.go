package http

import (
	"errors"
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

func (h *WODHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	var req CreateWODRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	cmd, err := toCreateCommand(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	result, err := h.service.Update(c.Request().Context(), id, cmd)
	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusOK, toDetailResponse(result))
}

func toCreateCommand(req CreateWODRequest) (appwod.CreateWODCommand, error) {
	if len(req.Stages) == 0 {
		return appwod.CreateWODCommand{}, domainwod.ErrStageRequired
	}

	stages := make([]appwod.StageInput, 0, len(req.Stages))
	for _, stage := range req.Stages {
		movements := make([]appwod.MovementInput, 0, len(stage.Movements))
		for _, movement := range stage.Movements {
			movements = append(movements, appwod.MovementInput{
				Position:     movement.Position,
				Label:        movement.Label,
				Name:         movement.Name,
				Prescription: movement.Prescription,
				Sets:         movement.Sets,
				Reps:         movement.Reps,
				LoadValue:    movement.LoadValue,
				LoadUnit:     movement.LoadUnit,
				Notes:        movement.Notes,
			})
		}

		stages = append(stages, appwod.StageInput{
			Kind:         domainwod.StageKind(stage.Kind),
			Instructions: stage.Instructions,
			Config: appwod.StageConfigInput{
				Type:            domainwod.WODType(stage.Type),
				TimeCapSeconds:  stage.Config.TimeCapSeconds,
				Rounds:          stage.Config.Rounds,
				WorkSeconds:     stage.Config.WorkSeconds,
				RestSeconds:     stage.Config.RestSeconds,
				Cycles:          stage.Config.Cycles,
				IntervalSeconds: stage.Config.IntervalSeconds,
			},
			Movements: movements,
		})
	}

	return appwod.CreateWODCommand{
		Name:        req.Name,
		Description: req.Description,
		Stages:      stages,
	}, nil
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
		errors.Is(err, domainwod.ErrInvalidMovementLabel),
		errors.Is(err, domainwod.ErrInvalidLoadUnit),
		errors.Is(err, domainwod.ErrInvalidSets),
		errors.Is(err, domainwod.ErrInvalidReps),
		errors.Is(err, domainwod.ErrInvalidPosition),
		errors.Is(err, domainwod.ErrUnknownWODType),
		errors.Is(err, domainwod.ErrStageRequired),
		errors.Is(err, domainwod.ErrInvalidStageKind),
		errors.Is(err, domainwod.ErrInvalidStagePosition),
		errors.Is(err, domainwod.ErrNilConfig),
		errors.Is(err, appwod.ErrMissingConfigField):
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, postgres.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "wod not found"})
	default:
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
