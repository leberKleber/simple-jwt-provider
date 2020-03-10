package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"os"
)

type config struct {
	ServerAddress    string `conf:"help:Server-address network-interface to bind on e.g.: '127.0.0.1:8080',default:0.0.0.0:80"`
	JWTPrivateKey    string `conf:"env:JWT_PRIVATE_KEY,help:JWT PrivateKey ECDSA512,required,noprint"`
	PasswordResetURL string `conf:"help: External URL to password-reset api. Password-reset-token would be replaced by %s,default:localhost:8080/password-reset?token=%s"`
	DB               struct {
		Host                 string `conf:"help:Database-Host,required"`
		Port                 int    `conf:"help:Database-Port,default:5432"`
		Name                 string `conf:"help:Database-name,default:'simple-jwt-provider'"`
		Username             string `conf:"help:Database-Username"`
		Password             string `conf:"help:Database-Password,noprint"`
		MigrationsFolderPath string `conf:"help:Database Migrations Folder Path,default:/db-migrations"`
	}
	AdminAPI struct {
		Enable   bool   `conf:"help:Enable admin API to manage stored users (true / false),default:false"`
		Username string `conf:"help:Basic Auth Username if enable-admin-api = true"`
		Password string `conf:"help:Basic Auth Password if enable-admin-api = true,noprint"`
	}
	Mail struct {
		TemplatesFolderPath string `conf:"help:Path to mail-templates folder,default:/mail-templates"`
		SMTPUsername        string `conf:"env:MAIL_SMTP_USERNAME,help:SMTP username to authorize with,required"`
		SMTPPassword        string `conf:"env:MAIL_SMTP_PASSWORD,help:SMTP password to authorize with,required,noprint"`
		SMTPHost            string `conf:"env:MAIL_SMTP_HOST,help:SMTP host to connect to,required"`
		SMTPPort            int    `conf:"env:MAIL_SMTP_PORT,help:SMTP port to connect to,default:587"`
		TLS                 struct {
			InsecureSkipVerify bool   `conf:"help:true if certificates should not be verified,default:false"`
			ServerName         string `conf:"help:name of the server who expose the certificate"`
		}
	}
}

func newConfig() (config, error) {
	cfg := config{}

	if origErr := conf.Parse(os.Environ(), "SJP", &cfg); origErr != nil {
		usage, err := conf.Usage("SJP", &cfg)
		if err != nil {
			return cfg, err
		}
		fmt.Println(usage)
		return cfg, origErr
	}

	if cfg.AdminAPI.Enable && (cfg.AdminAPI.Password == "" || cfg.AdminAPI.Username == "") {
		return cfg, errors.New("admin-api-password and admin-api-username must be set if api has been enabled")
	}

	return cfg, nil
}
