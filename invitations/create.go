package invitations

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/risersh/rest/rbac"
	"github.com/risersh/rest/util/database"
	"github.com/risersh/rest/util/database/prisma/db"

	"github.com/nvr-ai/go-util/security"
)

type CreateInvitationRequest struct {
	Type    string `json:"type" validate:"required,min=1,max=30"`
	Context string `json:"context" validate:"max=25"`
	Email   string `json:"email" validate:"max=255"`
	Role    string `json:"role" validate:"max=255"`
	Message string `json:"message" validate:"max=255"`
}

func Create(c echo.Context) error {
	if ok := rbac.HasPermission(c, rbac.AccountManage); !ok {
		return c.JSON(http.StatusForbidden, "Account Manage permission required")
	}

	tenant := c.Get("session").(*db.SessionModel).User().Tenant()
	var create CreateInvitationRequest
	if err := c.Bind(&create); err != nil {
		return c.JSON(http.StatusPreconditionFailed, err)
	}

	validate := validator.New()

	if err := validate.Struct(create); err != nil {
		log.Printf("validation error: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	code := strings.ToUpper(security.GenerateRandomString(6))

	if create.Type == string(InvitationTypeTenant) {
		create.Context = tenant.ID
	}

	_, err := database.DatabaseClient.Invitation.CreateOne(
		db.Invitation.Tenant.Link(db.Tenant.ID.Equals(tenant.ID)),
		db.Invitation.Email.Set(create.Email),
		db.Invitation.Role.Set(create.Role),
		db.Invitation.Code.Set(code),
		db.Invitation.Message.Set(create.Message),
		db.Invitation.Type.Set(create.Type),
		db.Invitation.Context.Set(create.Context),
	).Exec(context.Background())

	if err != nil {
		log.Printf("error: %v", err)

		if info, ok := db.IsErrUniqueConstraint(err); ok {
			return c.JSON(http.StatusConflict, info)
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, nil)
}
