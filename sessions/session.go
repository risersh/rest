package sessions

import (
	"context"
	"errors"
	"log"

	"github.com/risersh/rest/util/database"
	"github.com/risersh/rest/util/database/prisma/db"
)

// SessionContextKey is a string type that represents the key for the session in the context.
type SessionContextKey string

// SessionContextKey is a string type that represents the key for the session in the context.
const SessionContext SessionContextKey = "session"

// Session is a struct that represents a session.
type Session struct {
	User db.UserModel
}

type SessionStatus string

var (
	SessionStatusActive   SessionStatus = "ACTIVE"
	SessionStatusInactive SessionStatus = "INACTIVE"
)

type SessionClaims struct {
	ID string
}

func NewSession(userID string) (Session, error) {
	// session, err := database.DatabaseClient.Session.CreateOne(
	// 	db.Session.Status.Set("ACTIVE"),
	// 	db.Session.User.Link(db.User.ID.Equals(userID)),
	// ).With(
	// 	db.Session.User.Fetch().With(
	// 		db.User.Tenant.Fetch(),
	// 		db.User.Roles.Fetch().With(db.Role.Permission.Fetch()),
	// 	),
	// ).Exec(context.Background())
	session, err := database.DatabaseClient.Session.CreateOne(
		db.Session.Status.Set("ACTIVE"),
		db.Session.User.Link(db.User.ID.Equals(userID)),
	).Exec(context.Background())

	if err != nil {
		return Session{}, err
	}

	return Session{
		User: *session.User(),
	}, nil
}

func GetSession(id string) (*db.SessionModel, error) {
	session, err := database.DatabaseClient.Session.FindUnique(
		db.Session.ID.Equals(id),
	).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	if session.Status != string(SessionStatusActive) {
		return nil, errors.New("session is not active")
	}

	return session, nil
}

func HydrateSession(id string) (*db.SessionModel, error) {
	session, err := database.DatabaseClient.
		Session.
		FindUnique(
			db.Session.ID.Equals(id),
		).
		With(
			db.Session.User.Fetch().With(
				db.User.Tenant.Fetch(),
			),
		).
		Exec(context.Background())

	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	if session.User().Tenant().Status != "active" {
		return nil, errors.New("tenant is not active")
	}

	return session, nil
}
