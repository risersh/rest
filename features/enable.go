package features

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/risersh/rest/util/responses"
)

func Enable(c echo.Context) error {
	name := c.QueryParam("name")
	feature := GetFeature(name)
	if feature == nil {
		return responses.ReturnError(c, http.StatusNotFound, "feature not found")
	}

	return c.JSON(http.StatusOK, feature)
}
