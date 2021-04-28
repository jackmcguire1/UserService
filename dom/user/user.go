package user

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/jackmcguire1/UserService/pkg/utils"
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
	logEntry := log.WithField("user", utils.ToJSON(u))
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

		logEntry = logEntry.
			WithField("user-id", u.ID)

		logEntry.Debug("generated new uuid for user")
	}

	if u.CountryCode == "" || len(u.CountryCode) != 2 {
		return nil, fmt.Errorf("please enter a valid ISO ALPHA-2 country code", utils.ValidationErr)
	}
	if u.FirstName == "" {
		return nil, fmt.Errorf("%w - please enter a valid first name", utils.ValidationErr)
	}
	if u.LastName == "" {
		return nil, fmt.Errorf("%w - please enter a valid last name", utils.ValidationErr)
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

func (svc *service) DeleteUser(id string) error {
	return svc.Repo.DeleteUser(id)
}

func (svc *service) GetUsersByCountry(countryCode string) ([]*User, error) {
	logEntry := log.
		WithField("country-code", countryCode)

	logEntry.
		Info("call GetUsersByCountry")

	logEntry.Debug("querying get all users")
	users, err := svc.Repo.GetUsersByCountry(countryCode)
	if err != nil {
		logEntry.
			WithError(err).
			Error("failed to get all users from repository")

		return nil, err
	}

	logEntry.
		WithField("user-batch", utils.ToJSON(users)).
		Debug("got users from repository")

	return users, nil
}
