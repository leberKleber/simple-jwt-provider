package storage

import (
	"errors"
	"fmt"
	"time"
)

// ErrTokenNotFound returned when no token could be found
var ErrTokenNotFound = errors.New("no token found")

// TokenTypeReset identifies a token as reset-token. Then it can only be used for password-reset
const TokenTypeReset string = "reset"

// TokenTypeRefresh identifies a token as refresh-token. Then it can only be used  for refresh
const TokenTypeRefresh string = "refresh"

// Token represent a persisted token
type Token struct {
	ID        int64
	EMail     string
	Token     string
	Type      string
	CreatedAt time.Time
}

// CreateToken persists the given token in database. EMail must match to a users email.
func (s Storage) CreateToken(t Token) (int64, error) {
	var id int64
	err := s.db.QueryRow(
		"INSERT INTO tokens (email, token, type, created_at) VALUES($1, $2, $3, $4) RETURNING id;",
		t.EMail, t.Token, t.Type, t.CreatedAt,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to exec create token stmt: %w", err)
	}

	return id, nil
}

// TokensByEMailAndToken finds all tokens which matches the given email and token.
func (s Storage) TokensByEMailAndToken(email, token string) ([]Token, error) {
	rows, err := s.db.Query("SELECT id, type, created_at FROM tokens WHERE email = $1 AND token = $2;", email, token)
	if err != nil {
		return nil, fmt.Errorf("failed to exec select token stmt: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var tokens []Token
	for rows.Next() {
		t := Token{
			Token: token,
			EMail: email,
		}
		err := rows.Scan(&t.ID, &t.Type, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan select token stmt result: %w", err)
		}

		tokens = append(tokens, t)
	}

	return tokens, nil
}

// DeleteToken deletes token with the given ID.
// return ErrTokenNotFound there is no token with the given ID
func (s Storage) DeleteToken(id int64) error {
	res, err := s.db.Exec("DELETE FROM tokens WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	i, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get num of affected row: %w", err)
	}
	if i < 1 {
		return ErrTokenNotFound
	}

	return nil
}
