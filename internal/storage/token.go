package storage

import (
	"errors"
	"fmt"
	"time"
)

var ErrNoTokenHasBeenDeleted = errors.New("no token has been deleted")

const TokenTypeReset string = "reset"

type Token struct {
	ID        int64
	EMail     string
	Token     string
	Type      string
	CreatedAt time.Time
}

func (s Storage) CreateToken(t Token) (int64, error) {
	var id int64
	err := s.db.QueryRow(
		"INSERT INTO tokens (email, token, type, created_at) VALUES($1, $2, $3, $4) RETURNING id;",
		t.EMail, t.Token, t.Type, t.CreatedAt,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to exec stmt: %w", err)
	}

	return id, nil
}

func (s Storage) TokensByEMailAndToken(email, token string) ([]Token, error) {
	rows, err := s.db.Query("SELECT id, type, created_at FROM tokens WHERE email = $1 AND token = $2", email, token)
	if err != nil {
		return nil, fmt.Errorf("failed to exec select-token-stmt: %w", err)
	}
	defer rows.Close()

	var tokens []Token
	for rows.Next() {
		t := Token{
			Token: token,
		}
		err := rows.Scan(&t.ID, &t.Type, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan select-token-stmt result: %w", err)
		}

		tokens = append(tokens, t)
	}

	return tokens, nil
}

func (s Storage) DeleteToken(id int64) error {
	res, err := s.db.Exec("DELETE FROM tokens WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	i, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get num of affected row: %w", err)
	}
	if i != 1 {
		return ErrNoTokenHasBeenDeleted
	}

	return nil
}
