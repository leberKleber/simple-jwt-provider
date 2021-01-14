package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/leberKleber/simple-jwt-provider/internal"
	"github.com/leberKleber/simple-jwt-provider/internal/web/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Provider encapsulates internal.Provider to generate mocks
//go:generate moq -out provider_moq_test.go . Provider
type Provider interface {
	Login(email, password string) (string, string, error)
	Refresh(refreshToken string) (string, string, error)
	CreatePasswordResetRequest(email string) error
	ResetPassword(email, resetToken, password string) error
	CreateUser(user internal.User) error
	UpdateUser(email string, user internal.User) (internal.User, error)
	GetUser(email string) (internal.User, error)
	DeleteUser(email string) error
}

// Server should be created via NewServer and starts with ListenAndServe all http endpoints for this service.
type Server struct {
	h http.Handler
	p Provider
}

// NewServer returns a Server instance with configure http routs
func NewServer(p Provider, enableAdminAPI bool, adminAPIUsername, adminAPIPassword string) *Server {
	s := &Server{}
	r := mux.NewRouter()

	r.Use(contentTypeMiddleware)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)

	v1 := r.PathPrefix("/v1").Subrouter()
	v1.Path("/internal/alive").Methods(http.MethodGet).HandlerFunc(s.aliveHandler)
	v1.Path("/auth/login").Methods(http.MethodPost).HandlerFunc(s.loginHandler)
	v1.Path("/auth/refresh").Methods(http.MethodPost).HandlerFunc(s.refreshHandler)
	v1.Path("/auth/password-reset-request").Methods(http.MethodPost).HandlerFunc(s.passwordResetRequestHandler)
	v1.Path("/auth/password-reset").Methods(http.MethodPost).HandlerFunc(s.passwordResetHandler)

	if enableAdminAPI {
		adminAPI := v1.PathPrefix("/admin").Subrouter()
		adminAPI.Use(middleware.BasicAuth(adminAPIUsername, adminAPIPassword))

		adminAPI.Path("/users").Methods(http.MethodPost).HandlerFunc(s.createUserHandler)
		adminAPI.Path("/users/{email}").Methods(http.MethodGet).HandlerFunc(s.getUserHandler)
		adminAPI.Path("/users/{email}").Methods(http.MethodPut).HandlerFunc(s.updateUserHandler)
		adminAPI.Path("/users/{email}").Methods(http.MethodDelete).HandlerFunc(s.deleteUserHandler)
	}

	s.h = r
	s.p = p
	return s
}

// ListenAndServe wraps http.ListenAndServe
func (s *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, s.h)
}

func contentTypeMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

type errorResponseBody struct {
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	respBody, err := json.Marshal(errorResponseBody{
		Message: message,
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal json error response")
		writeInternalServerError(w)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(respBody)
	if err != nil {
		logrus.WithError(err).Error("Failed to write error response")
		writeInternalServerError(w)
		return
	}
}

func writeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(`{"message":"internal server error"}`))
	if err != nil {
		logrus.WithError(err).Error("Failed to write error response")
	}
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	writeError(w, http.StatusNotFound, "endpoint not found")
}

func methodNotAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}
