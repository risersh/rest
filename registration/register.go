package registration

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/risersh/notifications/email"
	"github.com/risersh/notifications/email/templates"
	"github.com/risersh/notifications/providers"
	"github.com/risersh/rest/conf"
	"github.com/risersh/rest/util/database"
	"github.com/risersh/rest/util/database/prisma/db"
	"github.com/risersh/rest/util/responses"
	"github.com/risersh/util/security"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	registered, err := database.DatabaseClient.Registration.FindFirst(
		db.Registration.Email.Equals(req.Email),
	).Exec(context.Background())

	if err != nil && err != db.ErrNotFound {
		return c.JSON(http.StatusInternalServerError, responses.ResponseError{
			Message: "an error occurred while processing the request",
		})
	} else if registered != nil {
		return c.JSON(http.StatusConflict, responses.ResponseError{
			Message: "registration already exists",
		})
	}

	existing, err := database.DatabaseClient.User.FindFirst(
		db.User.Email.Equals(req.Email),
	).Exec(context.Background())

	if err != nil && err != db.ErrNotFound {
		return c.JSON(http.StatusInternalServerError, responses.ResponseError{
			Message: "an error occurred while processing the request",
		})
	} else if existing != nil {
		return c.JSON(http.StatusConflict, responses.ResponseError{
			Message: "user already exists",
		})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	var code string

	if strings.Contains(req.Email, fmt.Sprintf("@%s", conf.Config.Branding.URL)) {
		code = conf.Config.Registration.Code
	} else {
		code = strings.ToUpper(security.GenerateRandomString(6))
	}

	user, err := database.DatabaseClient.Registration.CreateOne(
		db.Registration.Email.Set(req.Email),
		db.Registration.Password.Set(string(password)),
		db.Registration.Status.Set("PENDING"),
		db.Registration.Role.Set("OWNER"),
		db.Registration.Code.Set(code),
	).Exec(context.Background())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !strings.Contains(req.Email, conf.Config.Branding.URL) {
		err = email.SendTemplate(providers.ResendSendEmailArgs{
			To:      []string{req.Email},
			Subject: fmt.Sprintf("Welcome to %s!", conf.Config.Branding.Name),
		}, templates.RegistrationTemplate, map[string]string{
			"link": fmt.Sprintf("%s/registration/confirm?email=%s&token=%s", conf.Config.Branding.URL, req.Email, fmt.Sprint(code)),
			"code": fmt.Sprint(code),
		})
		if err != nil {
			multilog.Error("registration.register", "failed to send email", map[string]interface{}{
				"email": req.Email,
				"code":  fmt.Sprint(code),
				"link":  fmt.Sprintf("%s/registration/confirm?email=%s&token=%s", conf.Config.Branding.URL, req.Email, fmt.Sprint(code)),
				"error": err,
			})
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, user)
}
