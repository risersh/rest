package responses

import (
	"github.com/labstack/echo/v4"
	"github.com/mateothegreat/go-multilog/multilog"
)

type ResponseError struct {
	Message string `json:"message"`
}

func HandleException(c echo.Context, code int, message string, err error) error {
	multilog.Error("exception", "error", map[string]interface{}{
		"code":    code,
		"message": message,
		"error":   err,
	})
	return c.JSON(code, ResponseError{
		Message: message,
	})
}
