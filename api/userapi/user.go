package userapi

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jackmcguire1/UserService/api"
	"github.com/jackmcguire1/UserService/dom/user"
	"github.com/jackmcguire1/UserService/pkg/utils"
)

type CreateUserRequest struct {
	ID          string `json:"_id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	NickName    string `json:"nickName"`
	CountryCode string `json:"countryCode"`
	Saved       string `json:"saved"`
	Password    string `json:"password"`
	IsAdmin     bool   `json:"isAdmin"`
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With,Origin,Accept")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	h.Logger.
		With("raw-request", r).
		Debug("got new request")

	switch r.Method {
	case http.MethodGet:

		userParams, ok := r.URL.Query()["id"]
		if !ok || len(userParams[0]) < 1 {
			h.Logger.
				With("values", r.URL.Query()).
				Error("request does not contain 'id' query parameter")

			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.ToRAWJSON(api.HTTPError{Error: "missing 'id'  query parameter"}))

			return
		}
		userId := userParams[0]

		userResponse, err := h.getUser(userId)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				h.Logger.
					With("user-id", userId).
					With("error", err).
					Warn("user does not exist ")

				w.WriteHeader(http.StatusNotFound)
				w.Write(utils.ToRAWJSON(api.HTTPError{Error: "user not found"}))

				return
			}

			h.Logger.
				With("user-id", userId).
				With("error", err).
				Error("failed to get user")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(userResponse)

		return

	case http.MethodPost:
		reqData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.Logger.
				With("error", err).
				Error("failed to get read data from request body")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

			return
		}

		h.Logger.
			With("raw-body", string(reqData)).
			Info("got body from request")

		var user *user.User
		err = json.Unmarshal(reqData, &user)
		if err != nil {
			h.Logger.
				With("body", string(reqData)).
				Error("failed to get user data from request body")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

			return
		}

		if user.ID == "" {
			err := fmt.Errorf("missing user id")

			h.Logger.
				With("error", err).
				With("user", user).
				Error("failed to update user")

			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

			return
		}

		userResponse, err := h.UpdateUser(user)
		if err != nil {
			if errors.Is(err, utils.ValidationErr) {
				h.Logger.
					With("error", err).
					With("user", user).
					Error("failed to update user")

				w.WriteHeader(http.StatusBadRequest)
				w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

				return
			}

			if errors.Is(err, utils.AlreadyExists) {
				h.Logger.
					With("error", err).
					With("user", user).
					Error("failed to update user because of conflict")

				w.WriteHeader(http.StatusConflict)
				w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

				return
			}

			h.Logger.
				With("error", err).
				With("user", user).
				Error("failed to update user")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(userResponse)

		return
	case http.MethodPut:
		var user *CreateUserRequest
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			h.Logger.
				With("error", err).
				Error("failed to unmarshal user data from request body")

			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

			return
		}

		userResponse, err := h.createUser(user)
		if err != nil {
			if errors.Is(err, utils.AlreadyExists) {
				h.Logger.
					With("error", err).
					With("user", user).
					Warn("failed to create user")

				w.WriteHeader(http.StatusConflict)
				w.Write(utils.ToRAWJSON(api.HTTPError{Error: "user already exists"}))

				return
			}

			if errors.Is(err, utils.ValidationErr) {
				h.Logger.
					With("error", err).
					With("user", user).
					Error("failed to update user")

				w.WriteHeader(http.StatusBadRequest)
				w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

				return
			}

			h.Logger.
				With("error", err).
				With("user", user).
				Error("failed to create user")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(userResponse)

		h.Logger.
			With("response", string(userResponse)).
			With("user-id", user.ID).
			Debug("returning user")
		return

	case http.MethodDelete:
		type DeleteResponse struct {
			Deleted bool   `json:"deleted"`
			Message string `json:"message"`
		}

		userParams, ok := r.URL.Query()["id"]
		if !ok || len(userParams[0]) < 1 {
			h.Logger.
				With("url-values", r.URL.Query()).
				Error("request does not contain 'id' query parameter")

			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.ToRAWJSON(api.HTTPError{Error: "missing 'id'  query parameter"}))

			return
		}
		userId := userParams[0]

		h.Logger.
			With("user-id", userId).
			Info("got user to delete")

		err := h.UserService.DeleteUser(userId)
		if err != nil {

			if errors.Is(err, utils.ErrNotFound) {
				h.Logger.
					With("error", err).
					With("user-id", userId).
					Warn("user does not exist")

				w.WriteHeader(http.StatusBadRequest)
				w.Write(utils.ToRAWJSON(api.HTTPError{Error: "user does not exist"}))

				return
			}

			h.Logger.
				With("error", err).
				With("user-id", userId).
				Error("failed to delete user")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ToRAWJSON(api.HTTPError{Error: err.Error()}))

			return
		}

		h.Logger.
			With("user-id", userId).
			Debug("deleted user successfully")

		resp := utils.ToRAWJSON(&DeleteResponse{
			Deleted: true,
			Message: "success",
		})
		w.WriteHeader(http.StatusOK)
		w.Write(resp)

		h.Logger.
			With("response", string(resp)).
			With("user-id", userId).
			Debug("deleted user successfully")
		return

	default:
		err := fmt.Errorf("unsupported HTTP method")
		h.Logger.
			With("error", err).
			With("http-method", r.Method).
			Error("unsupported HTTP method requested")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ToRAWJSON(api.HTTPError{Error: "unsupported HTTP METHOD"}))
	}

	return
}

func (h *UserHandler) getUser(userId string) ([]byte, error) {
	logEntry := h.Logger.With("user-id", userId)
	logEntry.Info("call getUser - API")

	usr, err := h.UserService.GetUser(userId)
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(usr, "", "\t")
	if err != nil {
		return nil, err
	}

	h.Logger.
		With("user", string(b)).
		Info("got user")

	return b, err
}

func (h *UserHandler) UpdateUser(usr *user.User) ([]byte, error) {
	logEntry := h.Logger.With("user", usr)
	logEntry.Info("call UpdateUser - API")

	usr, err := h.UserService.PutUser(usr)
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(usr, "", "\t")
	if err != nil {
		return nil, err
	}

	h.Logger.
		With("user", string(b)).
		Info("got updated user")

	return b, err
}

func (h *UserHandler) createUser(usr *CreateUserRequest) ([]byte, error) {
	logEntry := h.Logger.With("user", usr)
	logEntry.Info("call createUser - API")

	if usr.ID != "" {
		existingUser, err := h.UserService.GetUser(usr.ID)
		if err != nil && !errors.Is(err, utils.ErrNotFound) {
			return nil, err
		}

		if existingUser != nil {
			return nil, fmt.Errorf("user already exists err: %w", utils.AlreadyExists)
		}
	}

	sha := sha256.New()
	sha.Write([]byte(usr.Password))
	password := sha.Sum(nil)

	newUser, err := h.UserService.PutUser(&user.User{
		ID:          usr.ID,
		FirstName:   usr.FirstName,
		LastName:    usr.LastName,
		Email:       usr.Email,
		NickName:    usr.NickName,
		CountryCode: usr.CountryCode,
		Password:    password,
		IsAdmin:     usr.IsAdmin,
	})
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(newUser, "", "\t")
	if err != nil {
		return nil, err
	}

	h.Logger.
		With("user", string(b)).
		Info("got new user")

	return b, err
}
