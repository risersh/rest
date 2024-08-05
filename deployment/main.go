package deployment

import "github.com/labstack/echo/v4"

func Router(e *echo.Echo) {
	e.POST("/deployment/:id/trigger", Trigger)

	e.GET("/deployment/:id/logs", Logs)
}
