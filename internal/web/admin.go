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

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		EMail    string                 `json:"email"`
		Password string                 `json:"password"`
		Claims   map[string]interface{} `json:"claims"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if requestBody.EMail == "" {
		writeError(w, http.StatusBadRequest, "email must be set")
		return
	}

	if requestBody.Password == "" {
		writeError(w, http.StatusBadRequest, "password must be set")
		return
	}

	err = s.p.CreateUser(requestBody.EMail, requestBody.Password, requestBody.Claims)
	if err != nil {
		if errors.Is(err, internal.ErrUserAlreadyExists) {
			writeError(w, http.StatusConflict, "user with given email already exists")
			return
		}

		logrus.WithError(err).Error("Failed to create user")
		writeInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
			writeError(w, http.StatusNotFound, "user with given email doesnt already exists")
			return
		} else {
			logrus.WithError(err).Error("Failed to delete user")
			writeInternalServerError(w)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
