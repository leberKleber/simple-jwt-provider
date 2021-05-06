package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var sqlOpen = gorm.Open

// Storage should be created via New and provides user and token database operation. Before access database Migrate should be called
type Storage struct {
	db *gorm.DB
}

// New opens a new sql connection with the given configuration with a connection timeout of 30
func New(dsn string) (*Storage, error) {
	//dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := sqlOpen(postgres.Open(dsn), &gorm.Config{})
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
