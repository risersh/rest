package deployment

import (
	"fmt"
	"github.com/labstack/echo/v4"
	riser "github.com/risersh/controller"
	"net/http"
)

func Trigger(c echo.Context) error {
	return c.JSON(http.StatusOK, fmt.Sprintf("Push a trigger %v", c.Param("id")))
}
