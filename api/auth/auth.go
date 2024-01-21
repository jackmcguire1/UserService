package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackmcguire1/UserService/dom/user"
)

const (
	AUTH_HEADER = "Auth"
)

var (
	UnAuthorizedErr   = fmt.Errorf("Unauthorized")
	InvalidRequestErr = fmt.Errorf("BadRequest")
)

type Handler struct {
	JWTSecret []byte
	Expiry    time.Duration
}

func (handler *Handler) SignClaims(usr *user.User) (string, error) {
	claims := &user.Claims{
		IsAdmin: usr.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: usr.ID,
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(handler.Expiry)),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(handler.JWTSecret)
}

func (handler *Handler) ValidateJWT(token string) (usrClaim *user.Claims, err error) {
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match

	usrClaim = &user.Claims{}

	tkn, err := jwt.ParseWithClaims(token, usrClaim, func(token *jwt.Token) (any, error) {
		return handler.JWTSecret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, UnAuthorizedErr
		}
		return nil, InvalidRequestErr
	}

	if !tkn.Valid {
		return nil, UnAuthorizedErr
	}

	return
}

func (handler *Handler) ValidateRequest(r *http.Request) (*user.Claims, error) {
	authHeader := r.Header.Get(AUTH_HEADER)

	items := strings.Split(authHeader, "Bearer ")
	if len(items) != 2 {
		return nil, InvalidRequestErr
	}

	return handler.ValidateJWT(items[1])
}
