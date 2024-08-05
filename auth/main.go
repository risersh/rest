package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/mateothegreat/go-rest-starter/middleware"
)

func Router(e *echo.Echo) *echo.Group {
	e.POST("/auth/login", Login)

	g := e.Group("/auth", middleware.SessionMiddleware())
	g.GET("/logout", Logout)

	return g
}
