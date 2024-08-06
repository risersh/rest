package rbac

import (
	"github.com/labstack/echo/v4"
)

func HasPermission(c echo.Context, permissions ...string) bool {
	return true
}
