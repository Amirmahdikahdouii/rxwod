package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	appauth "github.com/rxwod/backend/internal/application/auth"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/gym"
)

const gymIDHeader = "X-Gym-ID"

func AuthMiddleware(auth *appauth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get(echo.HeaderAuthorization)
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || strings.TrimSpace(token) == "" {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "authentication is required"})
			}
			userID, err := auth.AuthenticateAccessToken(strings.TrimSpace(token))
			if err != nil {
				return mapError(c, err)
			}
			ctx := appauthz.WithPrincipal(c.Request().Context(), appauthz.Principal{UserID: userID})
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func GymContextMiddleware(authorizer *appauthz.Authorizer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			principal, ok := appauthz.PrincipalFromContext(c.Request().Context())
			if !ok || principal.UserID == "" {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "authentication is required"})
			}
			headerGymID := strings.TrimSpace(c.Request().Header.Get(gymIDHeader))
			if headerGymID == "" {
				return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "x-gym-id header is required"})
			}
			resolved, err := authorizer.ResolveGym(c.Request().Context(), principal.UserID, gym.GymID(headerGymID))
			if err != nil {
				return mapError(c, err)
			}
			ctx := appauthz.WithPrincipal(c.Request().Context(), resolved)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func RequirePermission(permission domainauthz.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, err := appauthz.Require(c.Request().Context(), permission); err != nil {
				return mapError(c, err)
			}
			return next(c)
		}
	}
}
