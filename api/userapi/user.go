package userapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apex/log"
	"github.com/jackmcguire1/UserService/dom/user"
	"github.com/jackmcguire1/UserService/pkg/utils"
)

type UserHandler struct {
	UserService user.UserService
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.
		WithField("raw-request", r).
		Info("got new request")

	w.Header().Add("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:

		userParams, ok := r.URL.Query()["id"]
		if !ok || len(userParams[0]) < 1 {
			log.
				WithField("values", r.URL.Query()).
				Error("request does not contain 'id' query parameter")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing 'id'  query parameter"))

			return
		}
		userId := userParams[0]

		userResponse, err := h.getUser(userId)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				log.
					WithField("user-id", userId).
					WithError(err).
					Warn("user does not exist ")

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{}"))

				return
			}

			log.
				WithField("user-id", userId).
				WithError(err).
				Error("failed to get user")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		w.Write(userResponse)
		w.WriteHeader(http.StatusOK)

		return

	case http.MethodPost:
		reqData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.
				WithError(err).
				Error("failed to get read data from request body")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		log.
			WithField("raw-body", string(reqData)).
			Info("got body from request")

		var user *user.User
		err = json.Unmarshal(reqData, &user)
		if err != nil {
			log.
				WithField("body", string(reqData)).
				Error("failed to get user data from request body")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		if user.ID == "" {
			err := fmt.Errorf("missing user id")

			log.
				WithError(err).
				WithField("user", user).
				Error("failed to update user")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing user id"))

			return
		}

		userResponse, err := h.UpdateUser(user)
		if err != nil {
			if errors.Is(err, utils.ValidationErr) {
				log.
					WithError(err).
					WithField("user", user).
					Error("failed to update user")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))

				return
			}

			log.
				WithError(err).
				WithField("user", user).
				Error("failed to update user")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		w.Write(userResponse)
		w.WriteHeader(http.StatusOK)

		return
	case http.MethodPut:
		reqData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.
				WithError(err).
				Error("failed to get read data from request body")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		log.
			WithField("raw-body", string(reqData)).
			Info("got body from request")

		var user *user.User
		err = json.Unmarshal(reqData, &user)
		if err != nil {
			log.
				WithField("body", string(reqData)).
				WithError(err).
				Error("failed to unmarshal user data from request body")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		userResponse, err := h.createUser(user)
		if err != nil {

			if errors.Is(err, utils.AlreadyExists) {
				log.
					WithError(err).
					WithField("user", user).
					Warn("failed to create user")

				w.Write([]byte("user already exists"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if errors.Is(err, utils.ValidationErr) {
				log.
					WithError(err).
					WithField("user", user).
					Error("failed to update user")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))

				return
			}

			log.
				WithError(err).
				WithField("user", user).
				Error("failed to create user")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		w.Write(userResponse)
		w.WriteHeader(http.StatusCreated)

		return

	case http.MethodDelete:
		type DeleteResponse struct {
			Deleted bool
			Message string
		}

		userParams, ok := r.URL.Query()["id"]
		if !ok || len(userParams[0]) < 1 {
			log.
				WithField("url-values", r.URL.Query()).
				Error("request does not contain 'id' query parameter")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing 'id'  query parameter"))

			return
		}
		userId := userParams[0]

		log.
			WithField("user-id", userId).
			Info("got user to delete")

		err := h.UserService.DeleteUser(userId)
		if err != nil {

			if errors.Is(err, utils.ErrNotFound) {
				log.
					WithError(err).
					WithField("user-id", userId).
					Warn("user does not exist")

				w.Write([]byte(utils.ToJSON(&DeleteResponse{
					Deleted: false,
					Message: "user does not exist",
				})))
				w.WriteHeader(http.StatusOK)

				return
			}

			log.
				WithError(err).
				WithField("user-id", userId).
				Error("failed to delete user")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		log.
			WithField("user-id", userId).
			Debug("deleted user successfully")

		w.Write([]byte(utils.ToJSON(&DeleteResponse{
			Deleted: true,
			Message: "success",
		})))
		w.WriteHeader(http.StatusOK)
		return

	default:
		err := fmt.Errorf("unsupported HTTP method")
		log.
			WithError(err).
			WithField("http-method", r.Method).
			Error("unsupported HTTP method requested")

		w.WriteHeader(http.StatusBadRequest)
	}

	return
}

func (h *UserHandler) getUser(userId string) ([]byte, error) {
	logEntry := log.WithField("user-id", userId)
	logEntry.Info("call getUser - API")

	usr, err := h.UserService.GetUser(userId)
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(usr, "", "\t")
	if err != nil {
		return nil, err
	}

	log.
		WithField("user", string(b)).
		Info("got user")

	return b, err
}

func (h *UserHandler) UpdateUser(usr *user.User) ([]byte, error) {
	logEntry := log.WithField("user", usr)
	logEntry.Info("call putUser - API")

	usr, err := h.UserService.PutUser(usr)
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(usr, "", "\t")
	if err != nil {
		return nil, err
	}

	log.
		WithField("user", string(b)).
		Info("got updated user")

	return b, err
}

func (h *UserHandler) createUser(usr *user.User) ([]byte, error) {

	logEntry := log.WithField("user", usr)
	logEntry.Info("call getUser - API")

	if usr.ID != "" {
		existingUser, err := h.UserService.GetUser(usr.ID)
		if err != nil && !errors.Is(err, utils.ErrNotFound) {
			return nil, err
		}

		if existingUser != nil {
			return nil, fmt.Errorf("user already exists %w", utils.AlreadyExists)
		}
	}

	usr, err := h.UserService.PutUser(usr)
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(usr, "", "\t")
	if err != nil {
		return nil, err
	}

	log.
		WithField("user", string(b)).
		Info("got new user")

	return b, err
}
