package user

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackmcguire1/UserService/pkg/utils"
)

type User struct {
	ID          string `json:"_id" bson:"_id"`
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	Email       string `json:"email" bson:"email"`
	NickName    string `json:"nickName" bson:"nickName"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Saved       string `json:"saved" bson:"saved"`
}

func (svc *service) GetUser(userID string) (*User, error) {
	logEntry := slog.With("user-id", userID)
	logEntry.Info("call GetUser")

	user, err := svc.Repo.GetUser(userID)
	if err != nil {
		logEntry.
			With("error", err).
			Error("failed to get user")

		return nil, err
	}

	return user, err
}

func (svc *service) PutUser(u *User) (*User, error) {
	logEntry := slog.With("user", utils.ToJSON(u))
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
				With("error", err).
				Error("failed to generate a new uuid V4")

			return nil, err
		}
		u.ID = guid.String()

		logEntry = logEntry.
			With("user-id", u.ID)

		logEntry.Debug("generated new uuid for user")
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	u.Saved = time.Now().Format(time.RFC3339)
	u.CountryCode = strings.ToUpper(u.CountryCode)

	logEntry.Debug("saving user to repository")
	err := svc.Repo.PutUser(u)
	if err != nil {
		logEntry.
			With("error", err).
			Error("failed to put user into repository")

		return nil, err
	}

	if svc.UserChannel != nil {
		svc.UserChannel <- &UserUpdate{
			User:   u,
			Status: "UPDATE",
		}
	}

	return u, err
}

func (svc *service) DeleteUser(id string) error {
	err := svc.Repo.DeleteUser(id)
	if err != nil {
		return err
	}

	if svc.UserChannel != nil {
		svc.UserChannel <- &UserUpdate{
			User:   &User{ID: id},
			Status: "DELETED",
		}
	}

	return err
}

func (svc *service) GetUsersByCountry(countryCode string) ([]*User, error) {
	logEntry := slog.
		With("country-code", countryCode)

	logEntry.
		Info("call GetUsersByCountry")

	logEntry.Debug("querying get all users")
	users, err := svc.Repo.GetUsersByCountry(countryCode)
	if err != nil {
		logEntry.
			With("error", err).
			Error("failed to get all users from repository by country")

		return nil, err
	}

	logEntry.
		With("user-batch", utils.ToJSON(users)).
		Debug("got users from repository")

	return users, nil
}

func (svc *service) GetAllUsers() ([]*User, error) {
	slog.
		Info("call GetAllUsers")

	users, err := svc.Repo.GetAllUsers()
	if err != nil {
		slog.
			With("error", err).
			Error("failed to get all users from repository")
	}

	slog.
		With("user-batch", utils.ToJSON(users)).
		Debug("got all users from repository")

	return users, err
}

func (u *User) Validate() error {

	if u.CountryCode == "" || len(u.CountryCode) != 2 {
		return fmt.Errorf("%w - please enter a valid ISO ALPHA-2 country code", utils.ValidationErr)
	}
	if u.FirstName == "" {
		return fmt.Errorf("%w - please enter a valid first name", utils.ValidationErr)
	}
	if u.LastName == "" {
		return fmt.Errorf("%w - please enter a valid last name", utils.ValidationErr)
	}
	if u.Email == "" || !strings.Contains(u.Email, "@") {
		return fmt.Errorf("%w - please enter a valid email", utils.ValidationErr)
	}

	return nil
}
