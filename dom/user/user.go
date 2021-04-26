package user

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
)

type User struct {
	ID          string
	FirstName   string
	LastName    string
	CountryCode string
	Saved       string
}

func (svc *service) GetUser(userID string) (*User, error) {
	logEntry := log.WithField("user-id", userID)
	logEntry.Info("call GetUser")

	user, err := svc.Repo.GetUser(userID)
	if err != nil {
		logEntry.
			WithError(err).
			Error("failed to get user")

		return nil, err
	}

	return user, err
}

func (svc *service) PutUser(u *User) (*User, error) {
	logEntry := log.WithField("user", u)
	logEntry.Info("call PutUser")

	if u == nil {
		logEntry.Error("user struct not init")

		return nil, fmt.Errorf("user struct was nil")
	}

	if u.ID == "" {
		logEntry.Warn("no userID has been defined, generating new")

		guid, err := uuid.NewUUID()
		if err != nil {
			logEntry.
				WithError(err).
				Error("failed to generate a new uuid V4")

			return nil, err
		}
		u.ID = guid.String()

		logEntry.
			WithField("user-id", u.ID).
			Debug("uuid generated for new user")
	}

	u.Saved = time.Now().Format(time.RFC3339)

	logEntry.Debug("saving user to repository")
	err := svc.Repo.PutUser(u)
	if err != nil {
		logEntry.
			WithError(err).
			Error("failed to put user into repository")

		return nil, err
	}

	return u, err
}

func (svc *service) GetUsersByCountry(countryCode string, cursor string, limit int) ([]*User, string, error) {
	logEntry := log.
		WithField("country-code", countryCode)

	logEntry.
		WithField("cursor", cursor).
		WithField("limit", limit).
		Info("call GetUsersByCountry")

	logEntry.Debug("querying get all users")
	users, bookmark, err := svc.Repo.GetAllUsers(cursor, 100)
	if err != nil {
		logEntry.
			WithError(err).
			Error("failed to get all users from repository")

		return nil, "", err
	}

	logEntry.
		WithField("user-batch", users).
		Debug("got get all users after query")

	logEntry.
		WithField("len-users", len(users)).
		Debug("got users from repository")

	usersByCountry := []*User{}
	for _, user := range users {
		if user.CountryCode == countryCode {
			usersByCountry = append(usersByCountry, user)
		}
	}

	return usersByCountry, bookmark, nil
}
