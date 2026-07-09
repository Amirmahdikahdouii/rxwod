package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	appauth "github.com/rxwod/backend/internal/application/auth"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	appgym "github.com/rxwod/backend/internal/application/gym"
	appwod "github.com/rxwod/backend/internal/application/wod"
	domaingym "github.com/rxwod/backend/internal/domain/gym"
	domainuser "github.com/rxwod/backend/internal/domain/user"
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
	var statusFilter *domainwod.WODStatus
	if status := c.QueryParam("status"); status != "" {
		value := domainwod.WODStatus(status)
		statusFilter = &value
	}

	items, err := h.service.List(c.Request().Context(), statusFilter)
	if err != nil {
		return mapError(c, err)
	}

	responses := make([]WODSummaryResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toSummaryResponse(item))
	}
	return c.JSON(http.StatusOK, responses)
}

func (h *WODHandler) Calendar(c echo.Context) error {
	from, err := parseRequiredDate(c.QueryParam("from"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid from date"})
	}
	to, err := parseRequiredDate(c.QueryParam("to"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid to date"})
	}

	days, err := h.service.Calendar(c.Request().Context(), from, to)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toCalendarResponses(days))
}

func (h *WODHandler) Publish(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	result, err := h.service.Publish(c.Request().Context(), id)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toDetailResponse(result))
}

func (h *WODHandler) Archive(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	result, err := h.service.Archive(c.Request().Context(), id)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toDetailResponse(result))
}

func (h *WODHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
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

	scheduledDate, err := parseOptionalDate(req.ScheduledDate)
	if err != nil {
		return appwod.CreateWODCommand{}, err
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
		Name:          req.Name,
		Description:   req.Description,
		ScheduledDate: scheduledDate,
		Stages:        stages,
	}, nil
}

func parseOptionalDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := parseRequiredDate(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseRequiredDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
}

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, appauthz.ErrUnauthenticated):
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "authentication is required"})
	case errors.Is(err, appauthz.ErrGymRequired):
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "x-gym-id header is required"})
	case errors.Is(err, appauthz.ErrActiveMembershipMissing),
		errors.Is(err, appauthz.ErrForbidden):
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied"})
	case errors.Is(err, appauthz.ErrGymMismatch):
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "resource not found"})
	case errors.Is(err, appgym.ErrInvitationEmailMismatch):
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied"})
	case errors.Is(err, appauth.ErrInvalidCredentials),
		errors.Is(err, appauth.ErrRefreshTokenInvalid):
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	case errors.Is(err, domainwod.ErrInvalidName),
		errors.Is(err, appauth.ErrPasswordTooShort),
		errors.Is(err, domainuser.ErrInvalidEmail),
		errors.Is(err, domainuser.ErrInvalidDisplayName),
		errors.Is(err, domainuser.ErrPasswordHashEmpty),
		errors.Is(err, domaingym.ErrInvalidGymName),
		errors.Is(err, domaingym.ErrInvalidMembershipStatus),
		errors.Is(err, domaingym.ErrInvalidInvitationStatus),
		errors.Is(err, domaingym.ErrInvalidInvitationRole),
		errors.Is(err, domaingym.ErrInvitationNotPending),
		errors.Is(err, domaingym.ErrInvitationExpired),
		errors.Is(err, appgym.ErrRoleNotAssignable),
		errors.Is(err, appgym.ErrOwnerMembershipProtected),
		errors.Is(err, domainwod.ErrInvalidTimeCap),
		errors.Is(err, domainwod.ErrInvalidRounds),
		errors.Is(err, domainwod.ErrInvalidWorkSeconds),
		errors.Is(err, domainwod.ErrInvalidRestSeconds),
		errors.Is(err, domainwod.ErrInvalidInterval),
		errors.Is(err, domainwod.ErrInvalidCycles),
		errors.Is(err, domainwod.ErrMovementRequired),
		errors.Is(err, domainwod.ErrMovementTextRequired),
		errors.Is(err, domainwod.ErrLoadValueRequiresUnit),
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
		errors.Is(err, domainwod.ErrInvalidStatusTransition),
		errors.Is(err, domainwod.ErrCannotDeletePublished),
		errors.Is(err, domainwod.ErrScheduledDateRequired),
		errors.Is(err, appwod.ErrMissingConfigField):
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, postgres.ErrNotFound),
		errors.Is(err, appgym.ErrMemberNotFound),
		errors.Is(err, appgym.ErrInvitationNotFound):
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "resource not found"})
	default:
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
