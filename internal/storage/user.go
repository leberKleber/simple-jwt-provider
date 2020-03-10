package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type User struct {
	EMail    string
	Password []byte
}

var ErrUserNotFound = errors.New("could not found user")
var ErrUserAlreadyExists = errors.New("user already exists")

func (s *Storage) User(email string) (User, error) {
	user := User{
		EMail: email,
	}
	err := s.db.QueryRow(
		"SELECT password FROM users WHERE email = $1;",
		email,
	).Scan(&user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserNotFound
		}

		return user, fmt.Errorf("failed to query user cause: %w", err)
	}

	return user, nil
}

func (s *Storage) CreateUser(u User) error {
	stmt, err := s.db.Prepare("INSERT INTO users (email, password) VALUES($1, $2)")
	if err != nil {
		return fmt.Errorf("failed to prepare stmt: %w", err)
	}

	_, err = stmt.Exec(u.EMail, u.Password)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Constraint == "email_unique" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to exec stmt: %w", err)
	}

	return nil
}

func (s *Storage) UpdateUser(u User) error {
	stmt, err := s.db.Prepare("UPDATE users SET password = $2 WHERE email = $1;")
	if err != nil {
		return fmt.Errorf("failed to prepare stmt: %w", err)
	}

	res, err := stmt.Exec(u.EMail, u.Password)
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
