package deployments

import "github.com/labstack/echo/v4"

func Router(e *echo.Echo) {
	e.POST("/deployments/:id/trigger", Trigger)

	e.GET("/deployments/:id/logs", Logs)
}
