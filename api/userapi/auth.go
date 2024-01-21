package userapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jackmcguire1/UserService/pkg/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (handler *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With,Origin,Accept")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var loginReq *LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		handler.Logger.
			With("error", err).
			Error("failed to JSON decode login request")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usr, err := handler.UserService.GetUserByEmail(loginReq.Email)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sha := sha256.New()
	sha.Write([]byte(loginReq.Password))
	password := sha.Sum(nil)

	if !bytes.Equal(password, usr.Password) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tokenStr, err := handler.AuthHandler.SignClaims(usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: time.Now().UTC().Add(handler.AuthHandler.Expiry),
	})

	b, err := json.MarshalIndent(&LoginResponse{Token: tokenStr}, "", "\t")

	w.Write(b)
	w.WriteHeader(http.StatusOK)
	return
}
