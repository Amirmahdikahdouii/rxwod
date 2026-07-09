package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	appauth "github.com/rxwod/backend/internal/application/auth"
	appgym "github.com/rxwod/backend/internal/application/gym"
)

type AuthHandler struct {
	auth *appauth.Service
	gyms *appgym.Service
}

func NewAuthHandler(auth *appauth.Service, gyms *appgym.Service) *AuthHandler {
	return &AuthHandler{auth: auth, gyms: gyms}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	result, err := h.auth.Register(c.Request().Context(), appauth.RegisterCommand{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, toTokenResponse(result))
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	result, err := h.auth.Login(c.Request().Context(), appauth.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toTokenResponse(result))
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	result, err := h.auth.Refresh(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, toAccessTokenResponse(result))
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	if err := h.auth.RequestPasswordReset(c.Request().Context(), req.Email); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	if err := h.auth.ResetPassword(c.Request().Context(), req.Token, req.Password); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	var req VerifyEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}
	if err := h.auth.VerifyEmail(c.Request().Context(), req.Token); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) ResendVerification(c echo.Context) error {
	if err := h.auth.ResendVerificationEmail(c.Request().Context()); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) Me(c echo.Context) error {
	currentUser, err := h.auth.CurrentUser(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	gyms, err := h.gyms.ListForCurrentUser(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, MeResponse{
		User: toUserResponse(currentUser),
		Gyms: toWorkspaceResponses(gyms),
	})
}
