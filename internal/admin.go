package internal

import (
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

const blankedPassword = "**********"

var bcryptCost = 12
var ErrUserAlreadyExists = errors.New("user already exists")

type User struct {
	EMail    string
	Password string
	Claims   map[string]interface{}
}

/**
Creates new user with given email, password and claims.
return ErrUserAlreadyExists when user already exists
*/
func (p Provider) CreateUser(user User) error {
	securedPassword, err := bcryptPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to bcrypt password: %w", err)
	}

	err = p.Storage.CreateUser(storage.User{
		EMail:    user.EMail,
		Password: securedPassword,
		Claims:   user.Claims,
	})
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to query user with email %q: %w", user.EMail, err)
	}

	return nil
}

func (p Provider) GetUser(email string) (User, error) {
	user, err := p.Storage.User(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return User{}, ErrUserNotFound
		} else {
			return User{}, fmt.Errorf("failed to delete user with email %q: %w", email, err)
		}
	}

	return User{
		EMail:    user.EMail,
		Password: blankedPassword,
		Claims:   user.Claims,
	}, nil
}

// DeleteUser deletes user with given email.
// return ErrUserNotFound when user does not exist
// return ErrUserStillHasTokens when user still has tokens
func (p Provider) DeleteUser(email string) error {
	err := p.Storage.DeleteUser(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return ErrUserNotFound
		} else {
			return fmt.Errorf("failed to delete user with email %q: %w", email, err)
		}
	}

	return nil
}

func bcryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
}
