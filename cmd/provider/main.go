package main

import (
	"github.com/leberKleber/simple-jwt-provider/internal/web"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}

	logrus.WithField("config", cfg).Info("Start provider")
	if err := web.Serve(cfg.ServerAddress); err != nil {
		logrus.WithError(err).Fatal("Failed to run server")
	}
}
