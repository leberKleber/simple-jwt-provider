package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"os"
)

type config struct {
	ServerAddress string `conf:"help:Server-address network-interface to bind on e.g.: '127.0.0.1:8080',default:0.0.0.0:80"`
	JWTPrivateKey string `conf:"env:JWT_PRIVATE_KEY,help:JWT PrivateKey ECDSA512,required,noprint"`
	DB            struct {
		Host                 string `conf:"help:Database-Host,required"`
		Port                 int    `conf:"help:Database-Port,default:5432"`
		Name                 string `conf:"help:Database-name,default:'simple-jwt-provider'"`
		Username             string `conf:"help:Database-Username,required"`
		Password             string `conf:"help:Database-Password,required,noprint"`
		MigrationsFolderPath string `conf:"help:Database Migrations Folder Path,required"`
	}
	AdminAPI struct {
		Enable   bool   `conf:"help:Enable admin API to manage stored users (true / false),default:false"`
		Username string `conf:"help:Basic Auth Username if enable-admin-api = true,required"`
		Password string `conf:"help:Basic Auth Password if enable-admin-api = true,required,noprint"`
	}
	Mail struct {
		TemplateFolderPath string `conf:"help:Path to mail-template folder,default:/mail-templates"`
		SMTPUsername       string `conf:"help:SMTP username to authorize with"`
		SMTPPassword       string `conf:"help:SMTP password to authorize with,noprint"`
		SMTPHost           string `conf:"help:SMTP host to connect to"`
		SMTPPort           int    `conf:"help:SMTP port to connect to"`
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
