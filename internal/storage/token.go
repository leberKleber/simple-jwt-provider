package storage

import (
	"fmt"
	"time"
)

const TokenTypeReset string = "reset"

type Token struct {
	ID        int64
	EMail     string
	Token     string
	Type      string
	CreatedAt time.Time
}

func (s *Storage) CreateToken(t Token) (int64, error) {
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
