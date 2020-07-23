package internal

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrIncorrectPassword = errors.New("password incorrect")
var ErrUserNotFound = errors.New("user not found")
var ErrNoValidTokenFound = errors.New("no valid token found")
var nowFunc = time.Now

/**
Check email / password combination and return a new jwt if correct.
return ErrIncorrectPassword when password is incorrect
return UserNotFoundErr when user not found
*/
func (p Provider) Login(email, password string) (string, error) {
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

	return p.JWTGenerator.Generate(email, u.Claims)
}

/**
CreatePasswordResetRequest send a password-reset-request email to the give address.
return ErrUserNotFound when user does not exists
*/
func (p Provider) CreatePasswordResetRequest(email string) error {
	u, err := p.Storage.User(email)
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

	err = p.Mailer.SendPasswordResetRequestEMail(email, t, u.Claims)
	if err != nil {
		return fmt.Errorf("failed to send password-reset-email: %w", err)
	}

	return nil
}

/**
ResetPassword resets the password of the given account if the reset token is correct.
*/
func (p *Provider) ResetPassword(email, resetToken, newPassword string) error {
	tokens, err := p.Storage.TokensByEMailAndToken(email, resetToken)
	if err != nil {
		return fmt.Errorf("faild to find all avalilable tokens: %w", err)
	}

	var t *storage.Token
	for _, token := range tokens {
		if token.Type == storage.TokenTypeReset {
			//TODO check lifetime
			t = &token
			break
		}
	}

	if t == nil {
		return ErrNoValidTokenFound
	}

	u, err := p.Storage.User(email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	securedPassword, err := bcryptPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to bcrypt password: %w", err)
	}
	u.Password = securedPassword

	err = p.Storage.UpdateUser(u)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	err = p.Storage.DeleteToken(t.ID)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

//generate 64 char long hex token  (32 bytes == 64 hex chars)
func generateHEXToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return fmt.Sprintf("%x", b), err
}
