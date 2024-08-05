package deployment

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Logs(c echo.Context) error {
	return c.JSON(http.StatusOK, fmt.Sprintf("Logs for %v", c.Param("id")))
}
