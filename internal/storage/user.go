package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type User struct {
	EMail    string
	Password []byte
	Claims   map[string]interface{}
}

var ErrUserNotFound = errors.New("could not found user")
var ErrUserAlreadyExists = errors.New("user already exists")

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
			return user, ErrUserNotFound
		}

		return user, fmt.Errorf("failed to query user: %w", err)
	}

	err = json.Unmarshal(rawClaims, &user.Claims)
	if err != nil {
		return user, fmt.Errorf("failed to unmarshal user>claims: %w", err)
	}

	return user, nil
}

func (s *Storage) CreateUser(u User) error {
	stmt, err := s.db.Prepare("INSERT INTO users (email, password, claims) VALUES($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("failed to prepare stmt: %w", err)
	}

	rawClaims, err := json.Marshal(u.Claims)
	if err != nil {
		return fmt.Errorf("failed to marhsal user>claims: %w", err)
	}

	_, err = stmt.Exec(u.EMail, u.Password, rawClaims)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Constraint == "email_unique" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to exec stmt: %w", err)
	}

	return nil
}

func (s *Storage) UpdateUser(u User) error {
	stmt, err := s.db.Prepare("UPDATE users SET password = $2, claims = $3 WHERE email = $1;")
	if err != nil {
		return fmt.Errorf("failed to prepare stmt: %w", err)
	}

	rawClaims, err := json.Marshal(u.Claims)
	if err != nil {
		return fmt.Errorf("failed to marhsal user>claims: %w", err)
	}

	res, err := stmt.Exec(u.EMail, u.Password, rawClaims)
	if err != nil {
		return fmt.Errorf("failed to exec stmt: %w", err)
	}

	r, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get count of affected rows")
	}
	if r == 0 {
		return ErrUserNotFound
	}

	return nil
}
