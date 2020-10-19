package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/leberKleber/simple-jwt-provider/internal"
	"github.com/leberKleber/simple-jwt-provider/internal/jwt"
	"github.com/leberKleber/simple-jwt-provider/internal/mailer"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"github.com/leberKleber/simple-jwt-provider/internal/web"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}

	cfgAsString, err := conf.String(&cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Could not build config string")
	}
	fmt.Print(cfgAsString)
	logrus.Infof("Starting provider")

	s, err := storage.New(cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.Name)
	if err != nil {
		logrus.WithError(err).Fatal("Could not create storage")
	}

	err = s.Migrate(cfg.DB.MigrationsFolderPath)
	if err != nil {
		logrus.WithError(err).Fatal("Could not migrate database")
	}

	jwtGenerator, err := jwt.NewGenerator(cfg.JWT.PrivateKey, cfg.JWT.Audience, cfg.JWT.Issuer, cfg.JWT.Subject)
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

	provider := &internal.Provider{Storage: s, JWTGenerator: jwtGenerator, Mailer: m}
	server := web.NewServer(provider, cfg.AdminAPI.Enable, cfg.AdminAPI.Username, cfg.AdminAPI.Password)

	if err := server.ListenAndServe(cfg.ServerAddress); err != nil {
		logrus.WithError(err).Fatal("Failed to run server")
	}
}
