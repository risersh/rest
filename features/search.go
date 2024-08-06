package features

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/risersh/rest/util/responses"
)

func Search(c echo.Context) error {
	group := c.QueryParam("group")
	if group != "" {
		features := GetGroup(group)
		if features == nil {
			return responses.ReturnError(c, http.StatusNotFound, "group not found")
		}
		return c.JSON(http.StatusOK, features)
	}

	return c.JSON(http.StatusOK, Groups)
}
