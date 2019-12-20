package web

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) aliveHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(`{"alive":true}`))
	if err != nil {
		logrus.WithError(err).Error("Failed to write alive response body")
		writeInternalServerError(w)
	}
}
