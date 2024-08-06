package registration

import (
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo) *echo.Group {
	g := e.Group("/register")

	g.POST("", Register)
	g.POST("/confirm", ConfirmRegistration)

	return g
}
