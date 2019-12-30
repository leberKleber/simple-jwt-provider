package main

import (
	"github.com/leberKleber/simple-jwt-provider/internal"
	"github.com/leberKleber/simple-jwt-provider/internal/jwt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"github.com/leberKleber/simple-jwt-provider/internal/web"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}

	s, err := storage.New(cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseUsername, cfg.DatabasePassword, cfg.DatabaseName)
	if err != nil {
		logrus.WithError(err).Fatal("Could not create storage")
	}

	err = s.Migrate(cfg.DatabaseMigrationsFilePath)
	if err != nil {
		logrus.WithError(err).Fatal("Could not migrate database")
	}

	jwtGenerator, err := jwt.NewGenerator(cfg.JWTPrivateKey)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create jwt generator")
	}

	provider := &internal.Provider{Storage: s, JWTGenerator: jwtGenerator}
	server := web.NewServer(provider, cfg.EnableAdminAPI, cfg.AdminAPIUsername, cfg.AdminAPIPassword)

	logrus.WithField("config", cfg).Info("Start provider")

	if err := server.ListenAndServe(cfg.ServerAddress); err != nil {
		logrus.WithError(err).Fatal("Failed to run server")
	}
}
