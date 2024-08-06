package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/risersh/rest/conf"
	"github.com/risersh/rest/monitoring"
	"github.com/risersh/rest/sessions"
	"github.com/risersh/rest/util/database"
	"github.com/risersh/rest/util/database/prisma/db"
	"github.com/risersh/rest/util/responses"
	"github.com/risersh/util/security"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c echo.Context) error {
	parent, child := monitoring.NewSpanWithParent(c.Request().Context(), "login")
	defer parent.End()
	defer child.End()

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		child.AddEvent("bind.error", trace.WithAttributes(attribute.String("error", err.Error())))
		return responses.HandleException(c, http.StatusInternalServerError, "bad request", err)
	}

	child.AddEvent("get user record", trace.WithAttributes(attribute.String("email", req.Email), attribute.String("password", req.Password)))
	user, err := database.DatabaseClient.User.FindFirst(db.User.Email.Equals(req.Email)).Exec(context.Background())

	if err != nil {
		child.SetStatus(codes.Error, "could not find user")
		if db.IsErrNotFound(err) {
			return c.JSON(http.StatusUnauthorized, "invalid email address or password")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		child.SetStatus(codes.Error, "invalid password")
		return c.JSON(http.StatusUnauthorized, "invalid email address or password")
	}

	session, err := database.DatabaseClient.Session.CreateOne(
		db.Session.Status.Set("ACTIVE"),
		db.Session.User.Link(db.User.ID.Equals(user.ID)),
	).Exec(context.Background())

	if err != nil {
		child.SetStatus(codes.Error, fmt.Sprintf("could not create session: %s", err))
		return responses.HandleException(c, http.StatusInternalServerError, "could not create session", err)
	}

	_, span := monitoring.Tracer.Start(c.Request().Context(), "generate jwt")
	defer span.End()

	token := security.PasetoSign(sessions.PrivateKey, sessions.SessionClaims{
		ID: session.ID,
	}, time.Now().Add(time.Duration(conf.Config.Session.Duration)*time.Hour))

	span.AddEvent("generated jwt", trace.WithAttributes(attribute.String("jwt", token)))

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
