package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	serverAddress := "leberKleber.io"
	setEnv(t, "SJP_SERVER_ADDRESS", serverAddress)
	jwtPrivateKey := "myJWTKey"
	setEnv(t, "SJP_JWT_PRIVATE_KEY", jwtPrivateKey)
	jwtAudience := "myJWTAudience"
	setEnv(t, "SJP_JWT_AUDIENCE", jwtAudience)
	jwtIssuer := "myJWTIssuer"
	setEnv(t, "SJP_JWT_ISSUER", jwtIssuer)
	jwtSubject := "myJWTSubject"
	setEnv(t, "SJP_JWT_SUBJECT", jwtSubject)
	dbHost := "myDBHost"
	setEnv(t, "SJP_DB_HOST", dbHost)
	expectedDBPort := 555
	dbPort := "555"
	setEnv(t, "SJP_DB_PORT", dbPort)
	dbName := "myDBName"
	setEnv(t, "SJP_DB_NAME", dbName)
	dbUsername := "myDBUsername"
	setEnv(t, "SJP_DB_USERNAME", dbUsername)
	dbPassword := "myDBPassword"
	setEnv(t, "SJP_DB_PASSWORD", dbPassword)
	dbMigrationsFolderPath := "myDBMigrationsFolderPath"
	setEnv(t, "SJP_DB_MIGRATIONS_FOLDER_PATH", dbMigrationsFolderPath)
	expectedAdminAPIEnable := true
	adminAPIEnable := "true"
	setEnv(t, "SJP_ADMIN_API_ENABLE", adminAPIEnable)
	adminAPIUsername := "myAdminAPIUsername"
	setEnv(t, "SJP_ADMIN_API_USERNAME", adminAPIUsername)
	adminAPIPassword := "myAdminAPIPassword"
	setEnv(t, "SJP_ADMIN_API_PASSWORD", adminAPIPassword)
	mailTemplatesFolderPath := "myAdminAPIMailTemplatesFolderPath"
	setEnv(t, "SJP_MAIL_TEMPLATES_FOLDER_PATH", mailTemplatesFolderPath)
	mailSMTPHost := "myMailSMTPHost"
	setEnv(t, "SJP_MAIL_SMTP_HOST", mailSMTPHost)
	expectedMailSMTPPort := 42
	mailSMTPPort := "42"
	setEnv(t, "SJP_MAIL_SMTP_PORT", mailSMTPPort)
	mailSMTPUsername := "myMailSMTPUsername"
	setEnv(t, "SJP_MAIL_SMTP_USERNAME", mailSMTPUsername)
	mailSMTPPassword := "myMailSMTPPassword"
	setEnv(t, "SJP_MAIL_SMTP_PASSWORD", mailSMTPPassword)
	expectedMailTLSServerName := true
	mailTLSInsecureSkipVerify := "true"
	setEnv(t, "SJP_MAIL_TLS_INSECURE_SKIP_VERIFY", mailTLSInsecureSkipVerify)
	mailTLSServerName := "myMailTLSServerName"
	setEnv(t, "SJP_MAIL_TLS_SERVER_NAME", mailTLSServerName)

	cfg, err := newConfig()
	if err != nil {
		t.Fatalf("Unexpected error while building new config cuase: %s", err)
	}

	fieldEqual(t, "serverAddress", cfg.ServerAddress, serverAddress)
	fieldEqual(t, "jwt>privateKey", cfg.JWT.PrivateKey, jwtPrivateKey)
	fieldEqual(t, "jwt>audience", cfg.JWT.Audience, jwtAudience)
	fieldEqual(t, "jwt>issuer", cfg.JWT.Issuer, jwtIssuer)
	fieldEqual(t, "jwt>subject", cfg.JWT.Subject, jwtSubject)
	fieldEqual(t, "db>host", cfg.DB.Host, dbHost)
	fieldEqual(t, "db>port", cfg.DB.Port, expectedDBPort)
	fieldEqual(t, "db>name", cfg.DB.Name, dbName)
	fieldEqual(t, "db>username", cfg.DB.Username, dbUsername)
	fieldEqual(t, "db>password", cfg.DB.Password, dbPassword)
	fieldEqual(t, "db>migrationsFolderPath", cfg.DB.MigrationsFolderPath, dbMigrationsFolderPath)
	// noinspection GoBoolExpressions
	fieldEqual(t, "adminAPI>enable", cfg.AdminAPI.Enable, expectedAdminAPIEnable)
	fieldEqual(t, "adminAPI>username", cfg.AdminAPI.Username, adminAPIUsername)
	fieldEqual(t, "adminAPI>password", cfg.AdminAPI.Password, adminAPIPassword)
	fieldEqual(t, "mail>templatesFolderPath", cfg.Mail.TemplatesFolderPath, mailTemplatesFolderPath)
	fieldEqual(t, "mail>smtpHost", cfg.Mail.SMTPHost, mailSMTPHost)
	fieldEqual(t, "mail>smtpPort", cfg.Mail.SMTPPort, expectedMailSMTPPort)
	fieldEqual(t, "mail>smtpUsername", cfg.Mail.SMTPUsername, mailSMTPUsername)
	fieldEqual(t, "mail>smtpPassword", cfg.Mail.SMTPPassword, mailSMTPPassword)
	// noinspection GoBoolExpressions
	fieldEqual(t, "mail>tls>insecureSkipVerify", cfg.Mail.TLS.InsecureSkipVerify, expectedMailTLSServerName)
	fieldEqual(t, "mail>tls>serverName", cfg.Mail.TLS.ServerName, mailTLSServerName)
}

func TestNewConfigWithAdminAPIConstraint(t *testing.T) {
	cleanupEnvs(t)

	setEnv(t, "SJP_SERVER_ADDRESS", "leberKleber.io")
	setEnv(t, "SJP_JWT_PRIVATE_KEY", "myJWTKey")
	setEnv(t, "SJP_JWT_AUDIENCE", "myJWTAudience")
	setEnv(t, "SJP_JWT_ISSUER", "myJWTIssuer")
	setEnv(t, "SJP_JWT_SUBJECT", "myJWTSubject")
	setEnv(t, "SJP_DB_HOST", "myDBHost")
	setEnv(t, "SJP_DB_PORT", "555")
	setEnv(t, "SJP_DB_NAME", "myDBName")
	setEnv(t, "SJP_DB_USERNAME", "myDBUsername")
	setEnv(t, "SJP_DB_PASSWORD", "myDBPassword")
	setEnv(t, "SJP_DB_MIGRATIONS_FOLDER_PATH", "myDBMigrationsFolderPath")
	setEnv(t, "SJP_MAIL_TEMPLATES_FOLDER_PATH", "myAdminAPIMailTemplatesFolderPath")
	setEnv(t, "SJP_MAIL_SMTP_HOST", "myMailSMTPHost")
	setEnv(t, "SJP_MAIL_SMTP_PORT", "42")
	setEnv(t, "SJP_MAIL_SMTP_USERNAME", "myMailSMTPUsername")
	setEnv(t, "SJP_MAIL_SMTP_PASSWORD", "myMailSMTPPassword")
	setEnv(t, "SJP_MAIL_TLS_INSECURE_SKIP_VERIFY", "true")
	setEnv(t, "SJP_MAIL_TLS_SERVER_NAME", "myMailTLSServerName")

	// without username
	setEnv(t, "SJP_ADMIN_API_ENABLE", "true")
	setEnv(t, "SJP_ADMIN_API_USERNAME", "")
	setEnv(t, "SJP_ADMIN_API_PASSWORD", "myAdminAPIPassword")

	_, err := newConfig()
	expectedError := errors.New("admin-api-password and admin-api-username must be set if api has been enabled")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {

	}

	// without password
	setEnv(t, "SJP_ADMIN_API_ENABLE", "true")
	setEnv(t, "SJP_ADMIN_API_USERNAME", "myAdminAPIUsername")
	setEnv(t, "SJP_ADMIN_API_PASSWORD", "")

	_, err = newConfig()
	expectedError = errors.New("admin-api-password and admin-api-username must be set if api has been enabled")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Fatalf("returned error is not as expected. Expected:\n%s\nGiven:\n%s", expectedError, err)
	}

	cleanupEnvs(t)
}

func TestNewConfigCfgLibErrorHandling(t *testing.T) {
	cleanupEnvs(t)

	// PrivateKey cfg must be set but is not
	_, err := newConfig()

	expectedError := errors.New("required field PrivateKey is missing value")
	if fmt.Sprint(expectedError) != fmt.Sprint(err) {
		t.Fatalf("returned error is not as expected. Expected:\n%s\nGiven:\n%s", expectedError, err)
	}
}

func TestNewConfigUnableToGenerateUsage(t *testing.T) {
	oldConfigUsage := confUsage
	defer func() {
		confUsage = oldConfigUsage
	}()

	confUsage = func(namespace string, v interface{}) (string, error) {
		return "", errors.New("failed to generate usage")
	}

	cleanupEnvs(t)

	// PrivateKey cfg must be set but is not
	_, err := newConfig()

	expectedError := errors.New("failed to generate usage")
	if fmt.Sprint(expectedError) != fmt.Sprint(err) {
		t.Fatalf("returned error is not as expected. Expected:\n%s\nGiven:\n%s", expectedError, err)
	}

}

func setEnv(t *testing.T, key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		t.Fatalf("failed to set env variable %q cause: %s", key, err)
	}
}

func unsetEnv(t *testing.T, key string) {
	err := os.Unsetenv(key)
	if err != nil {
		t.Fatalf("failed to unset env variable %q cause: %s", key, err)
	}
}

func fieldEqual(t *testing.T, name string, cfgValue, expectedValue interface{}) {
	if !reflect.DeepEqual(cfgValue, expectedValue) {
		t.Errorf("unexpected cfg-value in field %q. Given: %s, Expected: %s", name, cfgValue, expectedValue)
	}
}

func cleanupEnvs(t *testing.T) {
	unsetEnv(t, "SJP_SERVER_ADDRESS")
	unsetEnv(t, "SJP_JWT_PRIVATE_KEY")
	unsetEnv(t, "SJP_JWT_AUDIENCE")
	unsetEnv(t, "SJP_JWT_ISSUER")
	unsetEnv(t, "SJP_JWT_SUBJECT")
	unsetEnv(t, "SJP_DB_HOST")
	unsetEnv(t, "SJP_DB_PORT")
	unsetEnv(t, "SJP_DB_NAME")
	unsetEnv(t, "SJP_DB_USERNAME")
	unsetEnv(t, "SJP_DB_PASSWORD")
	unsetEnv(t, "SJP_DB_MIGRATIONS_FOLDER_PATH")
	unsetEnv(t, "SJP_MAIL_TEMPLATES_FOLDER_PATH")
	unsetEnv(t, "SJP_MAIL_SMTP_HOST")
	unsetEnv(t, "SJP_MAIL_SMTP_PORT")
	unsetEnv(t, "SJP_MAIL_SMTP_USERNAME")
	unsetEnv(t, "SJP_MAIL_SMTP_PASSWORD")
	unsetEnv(t, "SJP_MAIL_TLS_INSECURE_SKIP_VERIFY")
	unsetEnv(t, "SJP_MAIL_TLS_SERVER_NAME")
	unsetEnv(t, "SJP_ADMIN_API_ENABLE")
	unsetEnv(t, "SJP_ADMIN_API_USERNAME")
	unsetEnv(t, "SJP_ADMIN_API_PASSWORD")
}
