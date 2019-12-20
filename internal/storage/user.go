package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	EMail    string
	Password []byte
}

var ErrUserNotFound = errors.New("could not found user")

func (s *Storage) User(email string) (User, error) {
	user := User{
		EMail: email,
	}
	err := s.db.QueryRow(
		"SELECT users.password FROM users WHERE users.email = $1;",
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
