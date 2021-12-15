package main

import (
	"github.com/ardanlabs/conf"
	"github.com/leberKleber/simple-jwt-provider/internal"
	"github.com/leberKleber/simple-jwt-provider/internal/jwt"
	"github.com/leberKleber/simple-jwt-provider/internal/mailer"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"github.com/leberKleber/simple-jwt-provider/internal/web"
	"github.com/sirupsen/logrus"
	"net/http"

	// database migration
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// sql driver
	_ "github.com/lib/pq"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := newConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}

	logLvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse log-level")
	}
	logrus.SetLevel(logLvl)

	cfgAsString, err := conf.String(&cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Could not build config string")
	}
	logrus.WithField("configuration", cfgAsString).Info("Starting provider")

	s, err := storage.NewPostgres(cfg.Database.Type, cfg.Database.DSN)
	if err != nil {
		logrus.WithError(err).Fatal("Could not create storage")
	}

	jwtGenerator, err := jwt.NewProvider(cfg.JWT.PrivateKey, cfg.JWT.Lifetime, cfg.JWT.Audience, cfg.JWT.Issuer, cfg.JWT.Subject)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create jwt generator")
	}

	m, err := mailer.New(cfg.Mail.TemplatesFolderPath,
		cfg.Mail.SMTPUsername,
		cfg.Mail.SMTPPassword,
		cfg.Mail.SMTPHost,
		cfg.Mail.SMTPPort,
		cfg.Mail.TLS.InsecureSkipVerify,
		cfg.Mail.TLS.ServerName,
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create mailer")
	}

	provider := &internal.Provider{Storage: s, JWTProvider: jwtGenerator, Mailer: m}
	server := web.NewServer(provider, cfg.AdminAPI.Enable, cfg.AdminAPI.Username, cfg.AdminAPI.Password)

	err = server.ListenAndServe(cfg.ServerAddress)
	if err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("Failed to run server")
	}
}
