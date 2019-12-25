package internal

import (
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var ErrIncorrectPassword = fmt.Errorf("password incorrect")
var ErrUserNotFound = fmt.Errorf("user not found")

/**
Check email / password combination and return a new jwt if correct.
return ErrIncorrectPassword when password is incorrect
return UserNotFoundErr when user not found
*/
func (p *Provider) Login(email, password string) (string, error) {
	u, err := p.Storage.User(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("failed to query user with email %q: %w", email, err)
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return "", ErrIncorrectPassword
	}

	return p.JWTGenerator.Generate(email)
}
