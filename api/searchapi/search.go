package searchapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/jackmcguire1/UserService/dom/user"
)

type SearchHandler struct {
	UserService user.UserService
}

func (h *SearchHandler) UsersByCountry(w http.ResponseWriter, r *http.Request) {
	type UserByCountryResponse struct {
		Users []*user.User
	}

	ccParams, ok := r.URL.Query()["cc"]
	if !ok || len(ccParams[0]) < 2 {
		log.
			WithField("values", r.URL.Query()).
			Error("request does not contain 'cc' query parameter")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing 'cc'  query parameter"))

		return
	}
	countryCode := strings.ToUpper(ccParams[0])

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
		w.Write([]byte(err.Error()))
		return
	}

	data, _ := json.MarshalIndent(&UserByCountryResponse{Users: users}, "", "\t")

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
