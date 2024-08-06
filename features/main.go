package features

import (
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo) {
	// g := e.Group("/features", middleware.SessionMiddleware())
	g := e.Group("/features")

	g.GET("", Search)
	g.PUT("", Enable)
}
