package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"os"
	"time"
)

var confUsage = conf.Usage

type config struct {
	LogLevel      string `conf:"env:LOG_LEVEL,help:Log-Level can be TRACE DEBUG INFO WARN ERROR FATAL or PANIC,default:INFO"`
	ServerAddress string `conf:"help:Server-address network-interface to bind on e.g.: '127.0.0.1:8080',default:0.0.0.0:80"`
	JWT           struct {
		Lifetime   time.Duration `conf:"env:JWT_LIFETIME,help:Lifetime of JWT,default:4h"`
		PrivateKey string        `conf:"env:JWT_PRIVATE_KEY,help:JWT PrivateKey ECDSA512,required,noprint"`
		Audience   string        `conf:"env:JWT_AUDIENCE,help:Audience private claim which will be applied in each JWT"`
		Issuer     string        `conf:"env:JWT_ISSUER,help:Issuer private claim which will be applied in each JWT"`
		Subject    string        `conf:"env:JWT_SUBJECT,help:Subject private claim which will be applied in each JWT"`
	}
	DSN      string `conf:"help:Data Source Name for persistence,default:host=localhost user=sjp password=sjp dbname=simple-jwt-provider port=5432 sslmode=disable"`
	AdminAPI struct {
		Enable   bool   `conf:"help:Enable admin API to manage stored users (true / false),default:false"`
		Username string `conf:"help:Basic Auth Username if enable-admin-api = true"`
		Password string `conf:"help:Basic Auth Password if enable-admin-api = true when is bcrypted prefix with 'bcrypt',noprint"`
	}
	Mail struct {
		TemplatesFolderPath string `conf:"help:Path to mail-templates folder,default:/mail-templates"`
		SMTPHost            string `conf:"env:MAIL_SMTP_HOST,help:SMTP host to connect to,required"`
		SMTPPort            int    `conf:"env:MAIL_SMTP_PORT,help:SMTP port to connect to,default:587"`
		SMTPUsername        string `conf:"env:MAIL_SMTP_USERNAME,help:SMTP username to authorize with,required"`
		SMTPPassword        string `conf:"env:MAIL_SMTP_PASSWORD,help:SMTP password to authorize with,required,noprint"`
		TLS                 struct {
			InsecureSkipVerify bool   `conf:"help:true if certificates should not be verified,default:false"`
			ServerName         string `conf:"help:name of the server who expose the certificate"`
		}
	}
}

func newConfig() (config, error) {
	cfg := config{}

	if origErr := conf.Parse(os.Environ(), "SJP", &cfg); origErr != nil {
		usage, err := confUsage("SJP", &cfg)
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
