package invitations

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/risersh/rest/rbac"
	"github.com/risersh/rest/util/database"
	"github.com/risersh/rest/util/database/prisma/db"
)

func Search(c echo.Context) error {
	if ok := rbac.HasPermission(c, rbac.AccountManage); !ok {
		return c.JSON(http.StatusForbidden, "Account Manage permission required")
	}

	result, err := database.DatabaseClient.Invitation.FindMany(
		db.Invitation.TenantID.Equals(c.Get("session").(*db.SessionModel).User().Tenant().ID),
		db.Invitation.Status.Not("deleted"),
	).OrderBy(
		db.Invitation.Created.Order(db.DESC),
	).Exec(context.Background())

	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	invitations := []Invitation{}

	for _, i := range result {
		invitations = append(invitations, Invitation{
			ID:      i.ID,
			Email:   i.Email,
			Role:    i.Role,
			Type:    i.Type,
			Context: i.Context,
			Status:  i.Status,
			Created: i.Created,
		})
	}

	return c.JSON(http.StatusOK, invitations)
}
