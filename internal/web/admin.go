package web

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/leberKleber/simple-jwt-provider/internal"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type User struct {
	EMail    string                 `json:"email"`
	Password string                 `json:"password"`
	Claims   map[string]interface{} `json:"claims"`
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if user.EMail == "" {
		writeError(w, http.StatusBadRequest, "email must be set")
		return
	}

	if user.Password == "" {
		writeError(w, http.StatusBadRequest, "password must be set")
		return
	}

	err = s.p.CreateUser(internal.User{
		EMail:    user.EMail,
		Password: user.Password,
		Claims:   user.Claims,
	})
	if err != nil {
		if errors.Is(err, internal.ErrUserAlreadyExists) {
			writeError(w, http.StatusConflict, "User with given email already exists")
			return
		}

		logrus.WithError(err).Error("Failed to create User")
		writeInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	email, err := url.PathUnescape(mux.Vars(r)["email"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not unescape email")
		return
	}

	//when email has not been set 'notFoundHandler' handler will be used

	user, err := s.p.GetUser(email)
	if err != nil {
		if errors.Is(err, internal.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "User with given email doesn't exists")
			return
		} else {
			logrus.WithError(err).Error("Failed to get User")
			writeInternalServerError(w)
			return
		}
	}

	err = json.NewEncoder(w).Encode(User{
		EMail:    user.EMail,
		Password: user.Password,
		Claims:   user.Claims,
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to encode User")
		writeInternalServerError(w)
		return
	}
}

func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	email, err := url.PathUnescape(mux.Vars(r)["email"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not unescape email")
		return
	}

	//when email has not been set 'notFoundHandler' handler will be used

	err = s.p.DeleteUser(email)
	if err != nil {
		if errors.Is(err, internal.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "User with given email doesnt already exists")
			return
		} else {
			logrus.WithError(err).Error("Failed to delete User")
			writeInternalServerError(w)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
