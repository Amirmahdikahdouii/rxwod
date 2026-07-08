package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	appgym "github.com/rxwod/backend/internal/application/gym"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
)

type GymHandler struct {
	service *appgym.Service
}

func NewGymHandler(service *appgym.Service) *GymHandler {
	return &GymHandler{service: service}
}

func (h *GymHandler) Create(c echo.Context) error {
	var req CreateGymRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	result, err := h.service.Create(c.Request().Context(), appgym.CreateGymCommand{Name: req.Name})
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, toGymResponse(result))
}

func (h *GymHandler) List(c echo.Context) error {
	result, err := h.service.ListForCurrentUser(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toWorkspaceResponses(result))
}

func (h *GymHandler) Get(c echo.Context) error {
	result, err := h.service.Get(c.Request().Context(), c.Param("gymId"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toGymResponse(result))
}

func (h *GymHandler) Members(c echo.Context) error {
	result, err := h.service.ListMembers(c.Request().Context(), c.Param("gymId"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toMemberResponses(result))
}

func (h *GymHandler) InviteCoach(c echo.Context) error {
	var req InviteMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	result, err := h.service.InviteCoach(c.Request().Context(), c.Param("gymId"), req.Email)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, toInvitationResponse(result))
}

func (h *GymHandler) InviteAthlete(c echo.Context) error {
	var req InviteMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	result, err := h.service.InviteAthlete(c.Request().Context(), c.Param("gymId"), req.Email)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, toInvitationResponse(result))
}

func (h *GymHandler) AcceptInvitation(c echo.Context) error {
	var req AcceptInvitationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	result, err := h.service.AcceptInvitation(c.Request().Context(), c.Param("gymId"), req.Token)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toMemberResponse(result))
}

func (h *GymHandler) UpdateMemberRole(c echo.Context) error {
	var req UpdateMemberRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	result, err := h.service.UpdateMemberRole(
		c.Request().Context(),
		c.Param("gymId"),
		c.Param("userId"),
		domainauthz.Role(req.Role),
	)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toMemberResponse(result))
}

func (h *GymHandler) RemoveMember(c echo.Context) error {
	if err := h.service.RemoveMember(c.Request().Context(), c.Param("gymId"), c.Param("userId")); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}
