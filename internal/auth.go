package internal

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrIncorrectPassword = fmt.Errorf("password incorrect")
var ErrUserNotFound = fmt.Errorf("user not found")
var nowFunc = time.Now

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

/**
CreatePasswordResetRequest send a password-reset-request email to the give address.
return ErrUserNotFound when user does not exists
*/
func (p *Provider) CreatePasswordResetRequest(email string) error {
	_, err := p.Storage.User(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to query user with email %q: %w", email, err)
	}

	t, err := generateHEXToken()
	if err != nil {
		return fmt.Errorf("failed to generate password-reset-token")
	}

	_, err = p.Storage.CreateToken(storage.Token{
		EMail:     email,
		Token:     t,
		Type:      storage.TokenTypeReset,
		CreatedAt: nowFunc(),
	})
	if err != nil {
		return fmt.Errorf("failed to create password-reset-token for email %q: %w", email, err)
	}

	err = p.Mailer.SendPasswordResetRequestEMail(email, fmt.Sprintf(p.PasswordResetURLFmt, t))
	if err != nil {
		return fmt.Errorf("failed to send password-reset-email: %w", err)
	}

	return nil
}

//generate 64 char long hex token  (32 bytes == 64 hex chars)
func generateHEXToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return fmt.Sprintf("%x", b), err
}
