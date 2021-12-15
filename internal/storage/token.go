package storage

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// ErrTokenNotFound returned when no token could be found
var ErrTokenNotFound = errors.New("no token found")

// TokenTypeReset identifies a token as reset-token. Then it can only be used for password-reset
const TokenTypeReset string = "reset"

// TokenTypeRefresh identifies a token as refresh-token. Then it can only be used  for refresh
const TokenTypeRefresh string = "refresh"

// Token represent a persisted token
type Token struct {
	gorm.Model
	EMail string
	Token string
	Type  string
}

// CreateToken persists the given token in database. EMail must match to a users email. ID will be set automatically.
func (s Postgres) CreateToken(t *Token) error {
	res := s.db.Create(t)

	if res.Error != nil {
		return fmt.Errorf("failed to exec create token: %w", res.Error)
	}

	return nil
}

// TokensByEMailAndToken finds all tokens which matches the given email and token.
func (s Postgres) TokensByEMailAndToken(email, token string) ([]Token, error) {
	var tokens []Token
	res := s.db.Find(&tokens, &Token{EMail: email, Token: token})

	if res.Error != nil {
		return nil, fmt.Errorf("failed to exec select token stmt: %w", res.Error)
	}

	return tokens, nil
}

// DeleteToken deletes token with the given ID.
// return ErrTokenNotFound there is no token with the given ID
func (s Postgres) DeleteToken(id uint) error {
	res := s.db.Delete(&Token{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete token: %w", res.Error)
	}

	if res.RowsAffected < 1 {
		return ErrTokenNotFound
	}

	return nil
}
