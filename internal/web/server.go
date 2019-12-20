package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Provider interface {
	Login(email, password string) (string, error)
}

type Server struct {
	h http.Handler
	p Provider
}

func NewServer(p Provider) *Server {
	s := &Server{}

	r := mux.NewRouter()

	v1 := r.PathPrefix("/v1")
	v1.Path("/internal/alive").Methods(http.MethodGet).HandlerFunc(s.aliveHandler)

	v1Auth := r.PathPrefix("/auth")
	v1Auth.Path("/login").Methods(http.MethodPost).HandlerFunc(s.loginHandler)

	s.h = r
	s.p = p
	return s
}

func (s *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, s.h)
}

func writeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(`{"message":"internal server error"}`))
	if err != nil {
		logrus.WithError(err).Error("Failed to write error response")
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	b, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal json error response")
		writeInternalServerError(w)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(b)
	if err != nil {
		logrus.WithError(err).Error("Failed to write error response")
		writeInternalServerError(w)
		return
	}
}
