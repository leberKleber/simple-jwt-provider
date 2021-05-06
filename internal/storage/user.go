package storage

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// User represent a persisted user
type User struct {
	gorm.Model
	EMail    string `gorm:"uniqueIndex:unique_email"`
	Password []byte
	Claims   Claims
}

// ErrUserNotFound returned when requested user not found
var ErrUserNotFound = errors.New("user not found")

// ErrUserAlreadyExists returned when given user already exists
var ErrUserAlreadyExists = errors.New("user already exists")

// CreateUser persists the given user in database
// return ErrUserNotFound when user not found
// return ErrUserAlreadyExists when user already exists
func (s *Storage) CreateUser(u User) error {
	res := s.db.Create(&u)
	if res.Error != nil {
		pqErr, ok := res.Error.(pq.Error)
		if ok && pqErr.Constraint == "unique_email" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to exec create user stmt: %w", res.Error)
	}

	return nil
}

// User finds the user identified by email
// return ErrUserNotFound when user not found
func (s *Storage) User(email string) (User, error) {
	var user User

	err := s.db.First(&user, User{EMail: email}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, ErrUserNotFound
	} else if err != nil {
		return User{}, fmt.Errorf("failed to query user: %w", err)
	}

	return user, nil
}

// UpdateUser updates all properties (excluding email) from the given user which will be identified by email
// return ErrUserNotFound when user not found
func (s *Storage) UpdateUser(u User) error {
	res := s.db.Updates(u)
	if res.Error != nil {
		return fmt.Errorf("failed to exec update user stmt: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// DeleteUser deletes the user with the given email and all corresponding tokes in one transaction.
// return ErrUserNotFound when user not found
func (s *Storage) DeleteUser(email string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.db.Delete(&Token{}, Token{EMail: email}).Error
		if err != nil {
			return fmt.Errorf("failed to exec delete tokens from user stmt: %w", err)
		}

		res := s.db.Delete(&User{}, User{EMail: email})
		if res.Error != nil {
			return fmt.Errorf("failed to exec delete user stmt: %w", res.Error)
		}

		if res.RowsAffected == 0 {
			return ErrUserNotFound
		}

		return nil
	})

	return err
}
