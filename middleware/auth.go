package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/risersh/rest-api/monitoring"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			parent := trace.SpanFromContext(c.Request().Context())
			defer parent.End()

			var tokenString string
			authHeader := c.Request().Header.Get("Authorization")

			_, span := monitoring.Tracer.Start(c.Request().Context(), "session middleware", trace.WithSpanKind(trace.SpanKindConsumer))
			defer span.End()

			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired session")
			}

			span.AddEvent("found token in header")
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")

			span.SetStatus(codes.Ok, "token located")
			span.SetAttributes(attribute.String("token", tokenString))

			return next(c)
		}
	}
}
