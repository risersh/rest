package invitations

import (
	"github.com/labstack/echo/v4"
	"github.com/risersh/rest/middleware"
)

func Router(e *echo.Echo) *echo.Group {
	e.POST("/invitations/confirm", Confirm)

	g := e.Group("/invitations", middleware.SessionMiddleware())

	g.GET("", Search)
	g.POST("", Create)
	g.DELETE("/:id", Delete)

	return g
}
