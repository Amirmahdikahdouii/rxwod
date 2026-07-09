package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	appauth "github.com/rxwod/backend/internal/application/auth"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	appgym "github.com/rxwod/backend/internal/application/gym"
	appwod "github.com/rxwod/backend/internal/application/wod"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
)

func NewRouter(authService *appauth.Service, gymService *appgym.Service, wodService *appwod.Service, authorizer *appauthz.Authorizer) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	authHandler := NewAuthHandler(authService, gymService)
	gymHandler := NewGymHandler(gymService)
	wodHandler := NewWODHandler(wodService)

	v1 := e.Group("/api/v1")
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)
	v1.POST("/auth/refresh", authHandler.Refresh)
	v1.POST("/auth/forgot-password", authHandler.ForgotPassword)
	v1.POST("/auth/reset-password", authHandler.ResetPassword)
	v1.POST("/auth/verify-email", authHandler.VerifyEmail)
	v1.GET("/invitations/:token", gymHandler.GetInvitationPreview)

	authenticated := v1.Group("", AuthMiddleware(authService))
	authenticated.GET("/me", authHandler.Me)
	authenticated.POST("/auth/resend-verification", authHandler.ResendVerification)
	authenticated.POST("/gyms", gymHandler.Create)
	authenticated.GET("/gyms", gymHandler.List)
	authenticated.POST("/gyms/:gymId/members/accept", gymHandler.AcceptInvitation)

	gymContext := authenticated.Group("", GymContextMiddleware(authorizer))
	gymContext.GET("/gyms/:gymId", gymHandler.Get, RequirePermission(domainauthz.PermissionGymRead))
	gymContext.GET("/gyms/:gymId/members", gymHandler.Members, RequirePermission(domainauthz.PermissionMemberList))
	gymContext.PATCH("/gyms/:gymId/members/:userId", gymHandler.UpdateMemberRole, RequirePermission(domainauthz.PermissionMemberUpdateRole))
	gymContext.DELETE("/gyms/:gymId/members/:userId", gymHandler.RemoveMember, RequirePermission(domainauthz.PermissionMemberRemove))
	gymContext.POST("/gyms/:gymId/coaches", gymHandler.InviteCoach, RequirePermission(domainauthz.PermissionMemberInviteCoach))
	gymContext.POST("/gyms/:gymId/athletes", gymHandler.InviteAthlete, RequirePermission(domainauthz.PermissionMemberInviteAthlete))
	gymContext.POST("/wods", wodHandler.Create, RequirePermission(domainauthz.PermissionWODCreate))
	gymContext.GET("/wods/calendar", wodHandler.Calendar, RequirePermission(domainauthz.PermissionWODRead))
	gymContext.GET("/wods", wodHandler.List, RequirePermission(domainauthz.PermissionWODRead))
	gymContext.GET("/wods/:id", wodHandler.GetByID, RequirePermission(domainauthz.PermissionWODRead))
	gymContext.PUT("/wods/:id", wodHandler.Update, RequirePermission(domainauthz.PermissionWODUpdate))
	gymContext.POST("/wods/:id/publish", wodHandler.Publish, RequirePermission(domainauthz.PermissionWODPublish))
	gymContext.DELETE("/wods/:id", wodHandler.Delete, RequirePermission(domainauthz.PermissionWODDelete))
	gymContext.POST("/wods/:id/archive", wodHandler.Archive, RequirePermission(domainauthz.PermissionWODDelete))

	return e
}
