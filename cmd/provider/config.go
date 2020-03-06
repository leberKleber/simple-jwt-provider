package main

import (
	"errors"
	"github.com/alexflint/go-arg"
)

type config struct {
	db
	adminAPI
	mail
	ServerAddress string `arg:"--server-address,env:SERVER_ADDRESS,help:Server-address network-interface to bind on e.g.: '127.0.0.1:8080' default ':80'"`
	JWTPrivateKey string `arg:"--jwt-private-key,env:JWT_PRIVATE_KEY,help:JWT PrivateKey ECDSA512,required"`
}

type db struct {
	DatabaseHost               string `arg:"--database-host,env:DATABASE_HOST,help:Database-Host"`
	DatabasePort               int    `arg:"--database-port,env:DATABASE_PORT,help:Database-Port default: '5432'"`
	DatabaseName               string `arg:"--database-name,env:DATABASE_NAME,help:Database-name default: 'simple-auth-provider'"`
	DatabaseUsername           string `arg:"--database-username,env:DATABASE_USERNAME,help:Database-Username"`
	DatabasePassword           string `arg:"--database-password,env:DATABASE_PASSWORD,help:Database-Password"`
	DatabaseMigrationsFilePath string `arg:"--database-migrations-file-path,env:DATABASE_MIGRATIONS_FILE_PATH,required,help:Database Migrations File Path"`
}

type adminAPI struct {
	EnableAdminAPI   bool   `arg:"--enable-admin-api,env:ENABLE_ADMIN_API,help:Enable admin API to manage stored users (true / false) default: 'true'"`
	AdminAPIUsername string `arg:"--admin-api-username,env:ADMIN_API_USERNAME,help:Basic Auth Username if enable-admin-api = true"`
	AdminAPIPassword string `arg:"--admin-api-password,env:ADMIN_API_PASSWORD,help:Basic Auth Password if enable-admin-api = true"`
}

type mail struct {
	MailFromAddress string `arg:"--mail-from-address,env:MAIL_FROM_ADDRESS,help:Mail address where the mails send from"`
	SMTPUsername    string `arg:"--smtp-username,env:SMTP_USERNAME,help:SMTP password to authorize with"`
	SMTPPassword    string `arg:"--smtp-password,env:SMTP_PASSWORD,help:SMTP password to authorize with"`
	SMTPHost        string `arg:"--smtp-host,env:SMTP_HOST,help:SMTP host to connect to"`
	SMTPPort        int    `arg:"--smtp-port,env:SMTP_PORT,help:SMTP port to connect to"`
}

func newConfig() (config, error) {
	cfg := config{
		ServerAddress: ":80",
		db: db{
			DatabasePort: 5432,
			DatabaseName: "simple-auth-provider",
		},
		adminAPI: adminAPI{
			EnableAdminAPI: false,
		},
	}
	err := arg.Parse(&cfg)

	if cfg.EnableAdminAPI && (cfg.AdminAPIPassword == "" || cfg.AdminAPIUsername == "") {
		return cfg, errors.New("admin-api-password and admin-api-username must be set if api has been enabled")
	}

	return cfg, err
}
