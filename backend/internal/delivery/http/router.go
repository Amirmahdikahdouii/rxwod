package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	appwod "github.com/rxwod/backend/internal/application/wod"
)

func NewRouter(service *appwod.Service) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	handler := NewWODHandler(service)
	v1 := e.Group("/api/v1")
	v1.POST("/wods", handler.Create)
	v1.GET("/wods", handler.List)
	v1.GET("/wods/:id", handler.GetByID)
	v1.PUT("/wods/:id", handler.Update)

	return e
}
