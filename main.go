package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/risersh/rest/auth"
	"github.com/risersh/rest/conf"
	"github.com/risersh/rest/deployments"
	"github.com/risersh/rest/features"
	"github.com/risersh/rest/invitations"
	"github.com/risersh/rest/monitoring"
	"github.com/risersh/rest/registration"
	"github.com/risersh/rest/sessions"
	"github.com/risersh/rest/util"
	"github.com/risersh/rest/util/database"
)

func main() {
	conf.Init()

	go monitoring.Setup()
	go util.ConnectRabbitMQ()
	go database.Connect(conf.Config.Database.URI)
	go sessions.InitKeys()

	shutdown, err := monitoring.InitTracer()
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())

	e := echo.New()
	// e.Use(monitoring.OtelMiddleware)
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: conf.Config.Server.Cors.Origins,
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
	features.Router(e)
	deployments.Router(e)
	registration.Router(e)
	invitations.Router(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Config.Server.Port)))
}
