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

func Delete(c echo.Context) error {
	if ok := rbac.HasPermission(c, rbac.AccountManage); !ok {
		return c.JSON(http.StatusForbidden, "Account Manage permission required")
	}

	id := c.Param("id")

	result, err := database.DatabaseClient.Invitation.FindMany(
		db.Invitation.ID.Equals(id),
		db.Invitation.Tenant.Where(
			db.Tenant.ID.Equals(c.Get("session").(*db.SessionModel).User().Tenant().ID),
		),
		db.Invitation.Status.Not("deleted"),
	).Update(
		db.Invitation.Status.Set("deleted"),
	).Exec(context.Background())

	if err != nil {
		log.Printf("error: %v", err)

		if ok := db.IsErrNotFound(err); ok {
			return c.JSON(http.StatusNotFound, "Invitation not found")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, result)
}
