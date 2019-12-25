package web

import (
	"encoding/json"
	"errors"
	"github.com/leberKleber/simple-jwt-provider/internal"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		EMail    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if requestBody.EMail == "" {
		writeError(w, http.StatusBadRequest, "EMail must be set")
		return
	}

	if requestBody.Password == "" {
		writeError(w, http.StatusBadRequest, "Password must be set")
		return
	}

	err = s.p.CreateUser(requestBody.EMail, requestBody.Password)
	if err != nil {
		if errors.Is(err, internal.ErrUserAlreadyExists) {
			writeError(w, http.StatusConflict, "User with given email already exists")
			return
		}

		logrus.WithError(err).Error("Failed to create user")
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}
