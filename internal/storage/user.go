package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

// User represent a persisted user
type User struct {
	EMail    string
	Password []byte
	Claims   map[string]interface{}
}

// ErrUserNotFound returned when requested user not found
var ErrUserNotFound = errors.New("user not found")

// ErrUserAlreadyExists returned when given user already exists
var ErrUserAlreadyExists = errors.New("user already exists")

// CreateUser persists the given user in database
// return ErrUserNotFound when user not found
// return ErrUserAlreadyExists when user already exists
func (s *Storage) CreateUser(u User) error {
	rawClaims, err := json.Marshal(u.Claims)
	if err != nil {
		return fmt.Errorf("failed to marhsal user>claims: %w", err)
	}

	_, err = s.db.Exec("INSERT INTO users (email, password, claims) VALUES($1, $2, $3);", u.EMail, u.Password, rawClaims)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Constraint == "email_unique" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to exec create user stmt: %w", err)
	}

	return nil
}

// User finds the user identified by email
// return ErrUserNotFound when user not found
func (s *Storage) User(email string) (User, error) {
	user := User{
		EMail: email,
	}
	var rawClaims []byte
	err := s.db.QueryRow(
		"SELECT password, claims FROM users WHERE email = $1;",
		email,
	).Scan(&user.Password, &rawClaims)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrUserNotFound
		}

		return User{}, fmt.Errorf("failed to query user: %w", err)
	}

	err = json.Unmarshal(rawClaims, &user.Claims)
	if err != nil {
		return User{}, fmt.Errorf("failed to unmarshal user>claims: %w", err)
	}

	return user, nil
}

// UpdateUser updates all properties (excluding email) from the given user which will be identified by email
// return ErrUserNotFound when user not found
func (s *Storage) UpdateUser(u User) error {
	rawClaims, err := json.Marshal(u.Claims)
	if err != nil {
		return fmt.Errorf("failed to marhsal user>claims: %w", err)
	}

	resp, err := s.db.Exec("UPDATE users SET password = $2, claims = $3 WHERE email = $1;", u.EMail, u.Password, rawClaims)
	if err != nil {
		return fmt.Errorf("failed to exec update user stmt: %w", err)
	}

	ra, err := resp.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get count of affected rows: %w", err)
	}
	if ra == 0 {
		return ErrUserNotFound
	}

	return nil
}

// DeleteUser deletes the user with the given email and all corresponding tokes in one transaction.
// return ErrUserNotFound when user not found
func (s *Storage) DeleteUser(email string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin delete transaction: %w", err)
	}

	_, err = tx.Exec("DELETE FROM tokens WHERE email = $1;", email)
	if err != nil {
		return fmt.Errorf("failed to exec delete tokens from user stmt: %w", err)
	}

	resp, err := tx.Exec("DELETE FROM users WHERE email = $1;", email)
	if err != nil {
		return fmt.Errorf("failed to exec delete user stmt: %w", err)
	}

	ra, err := resp.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get count of affected rows: %w", err)
	}
	if ra == 0 {
		return ErrUserNotFound
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit delete user transaction: %w", err)
	}

	return nil
}
