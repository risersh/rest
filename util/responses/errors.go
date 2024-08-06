package responses

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/mateothegreat/go-multilog/multilog"
)

type ResponseError struct {
	Message string `json:"message"`
}

func readRequestBody(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return buf.String()
}

func HandleException(c echo.Context, code int, message string, err error) error {
	multilog.Error("exception", "error", map[string]interface{}{
		"code":    code,
		"message": message,
		"error":   err,
		"request": Request{
			Method: c.Request().Method,
			Path:   c.Request().URL.Path,
			Query:  c.Request().URL.Query(),
			Body:   readRequestBody(c.Request().Body),
			Header: c.Request().Header,
		},
	})
	return c.JSON(code, ResponseError{
		Message: message,
	})
}

func ReturnError(c echo.Context, code int, message string) error {
	multilog.Error("error", "error", map[string]interface{}{
		"code":    code,
		"message": message,
		"request": Request{
			Method: c.Request().Method,
			Path:   c.Request().URL.Path,
			Query:  c.Request().URL.Query(),
			Body:   readRequestBody(c.Request().Body),
			Header: c.Request().Header,
		},
	})
	return c.JSON(code, ResponseError{
		Message: message,
	})
}
