package invitations

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/risersh/rest/util/database"
	"github.com/risersh/rest/util/database/prisma/db"
	"github.com/risersh/rest/util/responses"
)

type ConfirmInvitationRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,min=6,max=6"`
}

func Confirm(c echo.Context) error {
	var req ConfirmInvitationRequest

	if err := c.Bind(&req); err != nil {
		return responses.HandleException(c, http.StatusInternalServerError, "invalid request", err)
	}

	validate := validator.New()
	err := validate.Struct(req)

	if err != nil {
		return responses.HandleException(c, http.StatusInternalServerError, "invalid request", err)
	}

	invitation, err := database.DatabaseClient.Invitation.FindFirst(
		db.Invitation.Email.Equals(req.Email),
		db.Invitation.Code.Equals(req.Code),
	).Exec(context.Background())

	if err != nil && err == db.ErrNotFound {
		return responses.HandleException(c, http.StatusNotFound, "could not find registration", err)
	}

	var user *db.UserModel

	// Check if user exists.
	user, err = database.DatabaseClient.User.FindFirst(
		db.User.Email.Equals(req.Email),
	).Exec(context.Background())

	// If user does not exist, create one.
	if err != nil && err == db.ErrNotFound {
		user, err = database.DatabaseClient.User.CreateOne(
			db.User.Email.Set(req.Email),
			db.User.Password.Set(""),
			db.User.Status.Set("ACTIVE"),
			db.User.Tenant.Link(db.Tenant.ID.Equals(invitation.TenantID)),
		).Exec(context.Background())

		if err != nil {
			return responses.HandleException(c, http.StatusInternalServerError, "could not create user", err)
		}
	}

	// Update invitation status to confirmed.
	_, err = database.DatabaseClient.Invitation.FindUnique(
		db.Invitation.ID.Equals(invitation.ID),
	).Update(
		db.Invitation.User.Link(db.User.ID.Equals(user.ID)),
		db.Invitation.Status.Set("confirmed"),
	).Exec(context.Background())

	if err != nil {
		return responses.HandleException(c, http.StatusInternalServerError, "could not update invitation", err)
	}

	return c.JSON(http.StatusOK, "success")
}
