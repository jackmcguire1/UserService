package searchapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackmcguire1/UserService/api"
	"github.com/jackmcguire1/UserService/api/auth"
	"github.com/jackmcguire1/UserService/dom/user"
	"github.com/jackmcguire1/UserService/pkg/utils"
)

type SearchHandler struct {
	UserService user.UserService
	AuthHandler *auth.Handler
	Logger      *slog.Logger
}

func (h *SearchHandler) UsersByCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With,Origin,Accept")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	claims, err := h.AuthHandler.ValidateRequest(r)
	if err != nil {
		if errors.Is(err, auth.UnAuthorizedErr) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if errors.Is(err, auth.InvalidRequestErr) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !claims.IsAdmin {
		h.Logger.
			With("userID", claims.ID).
			With("error", "user is not administrator").
			Error("unauthenticated request")

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	logEntry := h.Logger.With("userID", claims.ID)

	logEntry.
		Info("search users by country request")

	type UserByCountryResponse struct {
		Users []*user.User `json:"users"`
	}

	w.Header().Add("Content-Type", "application/json")

	ccParams, ok := r.URL.Query()["cc"]
	if !ok || len(ccParams[0]) < 1 {
		logEntry.
			With("values", r.URL.Query()).
			Error("request does not contain 'cc' query parameter")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ToRAWJSON(api.HTTPError{Error: "missing 'cc'  query parameter"}))

		return
	}

	countryCode := strings.ToUpper(ccParams[0])
	if len(countryCode) != 2 {
		logEntry.
			With("country-code", countryCode).
			Error("request does not contain valid 'cc' query parameter")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ToRAWJSON(api.HTTPError{Error: "invalid 'cc'  query parameter - must be ISO ALPHA-2"}))

		return
	}

	logEntry.
		With("country-code", countryCode).
		Info("searching for users by country code")

	users, err := h.UserService.GetUsersByCountry(countryCode)
	if err != nil {
		logEntry.
			With("error", err).
			With("country-code", countryCode).
			Error("failed to get users by country code")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

		return
	}

	data, _ := json.MarshalIndent(&UserByCountryResponse{Users: users}, "", "\t")

	w.Write(data)
	w.WriteHeader(http.StatusOK)

	h.Logger.
		With("users", string(data)).
		Debug("returning users by country")

	return
}

func (h *SearchHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With,Origin,Accept")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	type AllUsers struct {
		Users []*user.User `json:"users"`
	}

	claims, err := h.AuthHandler.ValidateRequest(r)
	if err != nil {
		if errors.Is(err, auth.UnAuthorizedErr) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if errors.Is(err, auth.InvalidRequestErr) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !claims.IsAdmin {
		h.Logger.
			With("userID", claims.ID).
			With("error", "user is not administrator").
			Error("unauthenticated request")

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	h.Logger.
		Info("search all users")

	users, err := h.UserService.GetAllUsers()
	if err != nil {
		h.Logger.
			With("error", err).
			Error("failed to get all users")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

		return
	}

	data, _ := json.MarshalIndent(&AllUsers{Users: users}, "", "\t")

	w.Write(data)
	w.WriteHeader(http.StatusOK)

	h.Logger.
		With("users", string(data)).
		Debug("returning users")

	return
}
