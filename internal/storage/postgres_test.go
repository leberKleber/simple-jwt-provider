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

func TestNew(t *testing.T) {
	tests := []struct {
		name                          string
		givenSSLModeEnabled           bool
		sqlOpenError                  error
		expectedSQLOpenDataSourceName string
		expectedError                 error
	}{
		{
			name:                          "Happycase sqlMode enabled",
			givenSSLModeEnabled:           true,
			expectedSQLOpenDataSourceName: "host=myHost port=1234 user=myUsername password=myPassword dbname=myName sslmode=enable connect_timeout=30",
		},
		{
			name:                          "Happycase sqlMode disable",
			givenSSLModeEnabled:           false,
			expectedSQLOpenDataSourceName: "host=myHost port=1234 user=myUsername password=myPassword dbname=myName sslmode=disable connect_timeout=30",
		},
		{
			name:          "Error while sql.Open",
			sqlOpenError:  errors.New("nope"),
			expectedError: errors.New("failed to open database connection: nope"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldSQLOpen := sqlOpen
			defer func() {
				sqlOpen = oldSQLOpen
			}()

			var sqlOpenDriverName, sqlOpenDataSourceName string
			sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
				sqlOpenDriverName = driverName
				sqlOpenDataSourceName = dataSourceName

				return nil, tt.sqlOpenError
			}

			dbHost := "myHost"
			dbPort := 1234
			dbUsername := "myUsername"
			dbPassword := "myPassword"
			dbName := "myName"

			s, err := New(dbHost, dbPort, dbUsername, dbPassword, dbName, tt.givenSSLModeEnabled)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("Unexpeted error. Given:\n%q\nExpected:\n%q", err, tt.expectedError)
			} else if err != nil {
				return
			}

			if s.dbName != dbName {
				t.Errorf("Unexpected storage.dbName. Given: %q, expected: %q", s.dbName, dbName)
			}

			expectedSQLOpenDriverName := "postgres"
			if sqlOpenDriverName != expectedSQLOpenDriverName {
				t.Errorf("Unexpected sqlOpen.DriverName Given: %q, expected: %q", sqlOpenDriverName, expectedSQLOpenDriverName)
			}

			if sqlOpenDataSourceName != tt.expectedSQLOpenDataSourceName {
				t.Errorf("Unexpected sqlOpen.DataSourceName Given: \n%q\nexpected: \n%q", sqlOpenDataSourceName, tt.expectedSQLOpenDataSourceName)
			}
		})
	}
}

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
