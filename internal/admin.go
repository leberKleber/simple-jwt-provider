package internal

import (
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var bcryptCost = 10
var ErrUserAlreadyExists = fmt.Errorf("user already exists")

/**
Creates new user with given email and password
return ErrUserAllreadyExists when user already exists
*/
func (p *Provider) CreateUser(email, password string) error {
	securedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return fmt.Errorf("failed to bcrypt password: %w", err)
	}

	err = p.Storage.CreateUser(storage.User{
		EMail:    email,
		Password: securedPassword,
	})
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to query user with email %q: %w", email, err)
	}

	return nil
}
