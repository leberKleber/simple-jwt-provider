package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func BasicAuth(username, password string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok {
				unauthorized(w)
				return
			}

			if u != username || p != password {
				unauthorized(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte(`{"message": "Unauthorized"}`))
	if err != nil {
		logrus.WithError(err).Error("Failed to write unauthorized http response body")
	}
}
