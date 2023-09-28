package searchapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/jackmcguire1/UserService/api"
	"github.com/jackmcguire1/UserService/dom/user"
	"github.com/jackmcguire1/UserService/pkg/utils"
)

type SearchHandler struct {
	UserService user.UserService
}

func (h *SearchHandler) UsersByCountry(w http.ResponseWriter, r *http.Request) {
	type UserByCountryResponse struct {
		Users []*user.User `json:"users"`
	}

	w.Header().Add("Content-Type", "application/json")

	ccParams, ok := r.URL.Query()["cc"]
	if !ok || len(ccParams[0]) < 1 {
		log.
			WithField("values", r.URL.Query()).
			Error("request does not contain 'cc' query parameter")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ToRAWJSON(api.HTTPError{Error: "missing 'cc'  query parameter"}))

		return
	}

	countryCode := strings.ToUpper(ccParams[0])
	if len(countryCode) != 2 {
		log.
			WithField("country-code", countryCode).
			Error("request does not contain valid 'cc' query parameter")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ToRAWJSON(api.HTTPError{Error: "invalid 'cc'  query parameter - must be ISO ALPHA-2"}))

		return
	}

	log.
		WithField("country-code", countryCode).
		Info("searching for users by country code")

	users, err := h.UserService.GetUsersByCountry(countryCode)
	if err != nil {
		log.
			WithError(err).
			WithField("country-code", countryCode).
			Error("failed to get users by country code")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

		return
	}

	data, _ := json.MarshalIndent(&UserByCountryResponse{Users: users}, "", "\t")

	w.Write(data)
	w.WriteHeader(http.StatusOK)

	return
}

func (h *SearchHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	type AllUsers struct {
		Users []*user.User `json:"users"`
	}

	w.Header().Add("Content-Type", "application/json")

	log.
		Info("searching for users by country code")

	users, err := h.UserService.GetAllUsers()
	if err != nil {
		log.
			WithError(err).
			Error("failed to get all users")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

		return
	}

	data, _ := json.MarshalIndent(&AllUsers{Users: users}, "", "\t")

	w.Write(data)
	w.WriteHeader(http.StatusOK)

	return
}
