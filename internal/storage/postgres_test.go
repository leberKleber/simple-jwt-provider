package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"testing"
)

func TestStorage_Migrate(t *testing.T) {
	tests := []struct {
		name                                      string
		postgresWithInstanceReturnError           error
		migrateNewWithDatabaseInstanceReturnError error
		migrateUpError                            error
		expectedError                             error
	}{
		{
			name: "Happycase",
		},
		{
			name:                            "postgres instance error",
			postgresWithInstanceReturnError: errors.New("nope"),
			expectedError:                   errors.New("failed to create driver for database schema migration: nope"),
		},
		{
			name: "new migration error",
			migrateNewWithDatabaseInstanceReturnError: errors.New("nope"),
			expectedError: errors.New("failed to create a migrate for database schema migration: nope"),
		},
		{
			name:           "migrate up error",
			migrateUpError: errors.New("nope"),
			expectedError:  errors.New("failed to executed database schema migration: nope"),
		},
		{
			name:           "no migration change",
			migrateUpError: migrate.ErrNoChange,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldPostgresWithInstance := postgresWithInstance
			oldMigrateNewWithDatabaseInstance := migrateNewWithDatabaseInstance
			defer func() {
				postgresWithInstance = oldPostgresWithInstance
				migrateNewWithDatabaseInstance = oldMigrateNewWithDatabaseInstance
			}()

			postgresWithInstance = func(instance *sql.DB, config *postgres.Config) (database.Driver, error) {
				return nil, tt.postgresWithInstanceReturnError
			}
			migrateNewWithDatabaseInstance = func(sourceURL string, databaseName string, databaseInstance database.Driver) (migration, error) {
				return &migrationMock{
					UpFunc: func() error {
						return tt.migrateUpError
					},
				}, tt.migrateNewWithDatabaseInstanceReturnError
			}

			s := Storage{}

			err := s.Migrate("pathToDBMigrationFolder")
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("unexpected error. Expected:\n%q\nGiven:\n%q", tt.expectedError, err)
			}
		})
	}
}
