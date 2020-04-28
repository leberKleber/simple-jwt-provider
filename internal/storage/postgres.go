package storage

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	db *sql.DB
}

func New(dbHost string, dbPort int, dbUsername, dbPassword, dbName string) (*Storage, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=30", dbHost, dbPort, dbUsername, dbPassword, dbName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s Storage) Migrate(dbMigrationsPath string) error {
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver for database schema migration: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", dbMigrationsPath), "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create a migrate object for database schema migration: %w", err)
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

func (s Storage) Close() error {
	return s.db.Close()
}
