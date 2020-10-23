package storage

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sirupsen/logrus"
)

var postgresWithInstance = postgres.WithInstance
var migrateNewWithDatabaseInstance = func(sourceURL string, databaseName string, databaseInstance database.Driver) (migration, error) {
	m, e := migrate.NewWithDatabaseInstance(sourceURL, databaseName, databaseInstance)
	return m, e
}

type Storage struct {
	db     *sql.DB
	dbName string
}

//go:generate moq -out migration_moq_test.go . migration
type migration interface {
	Up() error
}

// New opens a new sql connection with the given configuration with a connection timeout of 30
func New(dbHost string, dbPort int, dbUsername, dbPassword, dbName string, sslModeEnabled bool) (*Storage, error) {
	sslMode := "disable"
	if sslModeEnabled {
		sslMode = "enable"
	}

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=30", dbHost, dbPort, dbUsername, dbPassword, dbName, sslMode),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return &Storage{
		db:     db,
		dbName: dbName,
	}, nil
}

// Migrate executes all sql migration files from the configures db-migrations folder. Should always be called before
// start
func (s Storage) Migrate(dbMigrationsPath string) error {
	driver, err := postgresWithInstance(s.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver for database schema migration: %w", err)
	}

	m, err := migrateNewWithDatabaseInstance(fmt.Sprintf("file://%s", dbMigrationsPath), s.dbName, driver)
	if err != nil {
		return fmt.Errorf("failed to create a migrate for database schema migration: %w", err)
	}

	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("failed to executed database schema migration: %w", err)
		}
		logrus.Info("no database schema changes")
		return nil
	}

	logrus.Info("executed database schema migration successfully")
	return nil
}

// Close warps sql.DB.Close
func (s Storage) Close() error {
	return s.db.Close()
}
