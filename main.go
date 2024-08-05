package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/risersh/rest/auth"
	"github.com/risersh/rest/conf"
	"github.com/risersh/rest/deployment"
	"github.com/risersh/rest/monitoring"
	"github.com/risersh/rest/util"
)

func main() {
	conf.Init()

	go util.ConnectRabbitMQ()

	shutdown, err := monitoring.InitTracer()
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())

	e := echo.New()
	// e.Use(monitoring.OtelMiddleware)
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:*",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderXRequestedWith,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	auth.Router(e)
	deployment.Router(e)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Config.Port)))
}
