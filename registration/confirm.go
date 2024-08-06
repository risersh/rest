package registration

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/risersh/rest/conf"
	"github.com/risersh/rest/sessions"
	"github.com/risersh/rest/util/database"
	"github.com/risersh/rest/util/database/prisma/db"
	"github.com/risersh/rest/util/responses"
	"github.com/risersh/util/security"
)

type ConfirmRegistrationRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,min=6,max=6"`
}

type ResponseError struct {
	Message string `json:"message"`
}

func ConfirmRegistration(c echo.Context) error {
	var req ConfirmRegistrationRequest

	if err := c.Bind(&req); err != nil {
		return responses.HandleException(c, http.StatusInternalServerError, "invalid request", err)
	}

	validate := validator.New()
	err := validate.Struct(req)

	if err != nil {
		return responses.HandleException(c, http.StatusInternalServerError, "invalid request", err)
	}

	registration, err := database.DatabaseClient.Registration.FindFirst(
		db.Registration.Email.Equals(req.Email),
		db.Registration.Code.Equals(req.Code),
	).Exec(context.Background())

	log.Printf("registration: %v", registration)
	if err != nil && err == db.ErrNotFound {
		return responses.HandleException(c, http.StatusInternalServerError, "could not find registration", err)
	}

	// database.DatabaseClient.Registration.FindUnique(
	// 	db.Registration.ID.Equals(registration.ID),
	// ).Delete().Exec(context.Background())

	tenant, err := database.DatabaseClient.Tenant.CreateOne(
		db.Tenant.Name.Set("Default"),
		db.Tenant.Status.Set("active"),
	).Exec(context.Background())

	if err != nil {
		return responses.HandleException(c, http.StatusInternalServerError, "could not create tenant", err)
	}

	log.Printf("tenant: %v", tenant)

	if err != nil {
		log.Printf("error: %v", err)
		return responses.HandleException(c, http.StatusInternalServerError, "could not find role", err)
	}

	user, err := database.DatabaseClient.User.CreateOne(
		db.User.Email.Set(req.Email),
		db.User.Password.Set(registration.Password),
		db.User.Status.Set("ACTIVE"),
		db.User.Tenant.Link(db.Tenant.ID.Equals(tenant.ID)),
	).Exec(context.Background())

	if err != nil {
		return responses.HandleException(c, http.StatusInternalServerError, "could not create user", err)
	}

	session, err := database.DatabaseClient.Session.CreateOne(
		db.Session.Status.Set("ACTIVE"),
		db.Session.User.Link(db.User.ID.Equals(user.ID)),
	).Exec(context.Background())

	if err != nil {
		return responses.HandleException(c, http.StatusInternalServerError, "could not create session", err)
	}

	token := security.PasetoSign(sessions.PrivateKey, sessions.SessionClaims{
		ID: session.ID,
	}, time.Now().Add(time.Duration(conf.Config.Session.Duration)*time.Hour))

	if err != nil {
		return responses.HandleException(c, http.StatusInternalServerError, "could not sign jwt", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
