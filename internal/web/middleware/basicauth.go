package middleware

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

const bcryptedPasswordPrefix = "bcrypt:"

// BasicAuth builds a basic auth http.Handler middleware which blocks all unauthorized request and respond with a
// http status 403
func BasicAuth(username, password string) func(h http.Handler) http.Handler {
	passwordIsBcrypted := strings.HasPrefix(password, bcryptedPasswordPrefix)
	if passwordIsBcrypted {
		password = strings.Replace(password, bcryptedPasswordPrefix, "", 1)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok {
				unauthorized(w)
				return
			}

			if u != username {
				unauthorized(w)
				return
			}

			if passwordIsBcrypted {
				err := bcrypt.CompareHashAndPassword([]byte(password), []byte(p))
				if err != nil {
					unauthorized(w)
					return
				}
			} else {
				if p != password {
					unauthorized(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	_, err := w.Write([]byte(`{"message": "forbidden"}`))
	if err != nil {
		logrus.WithError(err).Error("Failed to write unauthorized http response body")
	}
}
