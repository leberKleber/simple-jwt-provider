package storage

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbTypePostgres = "postgres"
const dbTypeSQLite = "sqlite"

var sqlOpen = gorm.Open

// Storage should be created via New and provides user and token database operation. Before access database Migrate should be called
type Storage struct {
	db *gorm.DB
}

// New opens a new sql connection with the given configuration
func New(dbType, dsn string) (*Storage, error) {
	dialector, err := buildDialector(dbType, dsn)
	if err != nil {
		return nil, err
	}

	db, err := sqlOpen(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.AutoMigrate(User{}, Token{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate persistence: %w", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func buildDialector(dbType, dsn string) (gorm.Dialector, error) {
	var dialector gorm.Dialector

	switch dbType {
	case dbTypePostgres:
		dialector = postgres.Open(dsn)
	case dbTypeSQLite:
		dialector = sqlite.Open(dsn)
	default:
		return nil, errors.New("unsupported database type")
	}

	return dialector, nil
}
