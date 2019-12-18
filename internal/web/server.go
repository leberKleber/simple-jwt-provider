package web

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Serve(address string) error {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/v1")
	v1.Path("/internal/alive").Methods(http.MethodGet).HandlerFunc(aliveHandler)

	return http.ListenAndServe(address, r)
}

func writeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(`{"message":"internal server error"}`))
	if err != nil {
		logrus.WithError(err).Error("Failed to write error response")
	}
}
