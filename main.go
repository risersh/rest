package main

import (
	"context"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mateothegreat/go-rest-starter/auth"
	"github.com/mateothegreat/go-rest-starter/monitoring"
)

func main() {
	godotenv.Load("../.env.local")

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

	e.Logger.Fatal(e.Start(":8080"))
}
