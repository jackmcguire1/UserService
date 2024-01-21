package userapi

import (
	"log/slog"

	"github.com/jackmcguire1/UserService/api/auth"
	"github.com/jackmcguire1/UserService/dom/user"
)

type UserHandler struct {
	UserService user.UserService
	Logger      *slog.Logger
	AuthHandler *auth.Handler
}
