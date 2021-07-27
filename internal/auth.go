package internal

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

// ErrIncorrectPassword returned when user authentication failed cause incorrect password
var ErrIncorrectPassword = errors.New("password incorrect")

// ErrUserNotFound returned when requested user not found
var ErrUserNotFound = errors.New("user not found")

// ErrNoValidTokenFound returned when requested user has no valid token
var ErrNoValidTokenFound = errors.New("no valid token found")

// ErrInvalidToken returned when the give token is not valid
var ErrInvalidToken = errors.New("given token is invalid")

// ErrTokenNotParsable returned when the give token is not parsable
var ErrTokenNotParsable = errors.New("given token is not parsable")

// Login checks email / password combination and return a new access and refresh token if correct.
// return ErrIncorrectPassword when password is incorrect
// return ErrUserNotFound when user not found
func (p Provider) Login(email, password string) (accessToken, refreshToken string, err error) {
	u, err := p.Storage.User(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", "", ErrUserNotFound
		}
		return "", "", fmt.Errorf("failed to find user with email %q: %w", email, err)
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return "", "", ErrIncorrectPassword
	}

	accessToken, err = p.JWTProvider.GenerateAccessToken(email, u.Claims)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access-token: %w", err)
	}

	refreshToken, jwtID, err := p.JWTProvider.GenerateRefreshToken(email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh-token: %w", err)
	}

	err = p.Storage.CreateToken(&storage.Token{
		EMail: email,
		Token: jwtID,
		Type:  storage.TokenTypeRefresh,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to persist refresh-token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// Refresh checks user and token validity and return a new access and refresh token if everything is valid
// return ErrTokenNotParsable when the token is not parsable
// return ErrInvalidToken when the token is not valid
// return ErrUserNotFound when the referred user could not be found
func (p Provider) Refresh(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	isValid, claims, err := p.JWTProvider.IsTokenValid(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("%w: %s", ErrTokenNotParsable, err)
	}

	if !isValid {
		return "", "", ErrInvalidToken
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", "", errors.New("email claim is not parsable as string")
	}

	tokenID, ok := claims["jit"].(string)
	if !ok {
		return "", "", errors.New("jit claim is not parsable as string")
	}

	//TODO do Storage.TokensByEMailAndToken and Storage.DeleteToken in transaction
	tokens, err := p.Storage.TokensByEMailAndToken(email, tokenID)
	if err != nil {
		return "", "", fmt.Errorf("failed to find refresh-tokens: %w", err)
	}

	var t *storage.Token
	for _, token := range tokens {
		if token.Type == storage.TokenTypeRefresh {
			// TODO check lifetime
			t = &token
			break
		}
	}

	if t == nil {
		return "", "", ErrNoValidTokenFound
	}

	err = p.Storage.DeleteToken(t.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to delete refresh-token: %w", err)
	}

	u, err := p.Storage.User(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", "", ErrUserNotFound
		}
		return "", "", fmt.Errorf("failed to find user with email %q: %w", email, err)
	}

	newAccessToken, err = p.JWTProvider.GenerateAccessToken(email, u.Claims)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access-token: %w", err)
	}

	newRefreshToken, jwtID, err := p.JWTProvider.GenerateRefreshToken(email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh-token: %w", err)
	}

	err = p.Storage.CreateToken(&storage.Token{
		EMail: email,
		Token: jwtID,
		Type:  storage.TokenTypeRefresh,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to persist refresh-token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// CreatePasswordResetRequest send a password-reset-request email to the give address.
// return ErrUserNotFound when user does not exists
func (p Provider) CreatePasswordResetRequest(email string) error {
	u, err := p.Storage.User(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user with email %q: %w", email, err)
	}

	t, err := generateHEXToken()
	if err != nil {
		return fmt.Errorf("failed to generate password reset token: %w", err)
	}

	err = p.Storage.CreateToken(&storage.Token{
		EMail: email,
		Token: t,
		Type:  storage.TokenTypeReset,
	})
	if err != nil {
		return fmt.Errorf("failed to create password reset token for email %q: %w", email, err)
	}

	err = p.Mailer.SendPasswordResetRequestEMail(email, t, u.Claims)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

// ResetPassword resets the password of the given account if the reset token is correct.
// return ErrNoValidTokenFound no valid token could be found
func (p *Provider) ResetPassword(email, resetToken, newPassword string) error {
	//TODO do Storage.TokensByEMailAndToken and Storage.DeleteToken in transaction
	tokens, err := p.Storage.TokensByEMailAndToken(email, resetToken)
	if err != nil {
		return fmt.Errorf("failed to find reset-tokens: %w", err)
	}

	var t *storage.Token
	for _, token := range tokens {
		if token.Type == storage.TokenTypeReset {
			// TODO check lifetime
			t = &token
			break
		}
	}

	if t == nil {
		return ErrNoValidTokenFound
	}

	u, err := p.Storage.User(email)
	if err != nil {
		return fmt.Errorf("failed to find user with email %q: %w", email, err)
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

// generate 64 char long hex token  (32 bytes == 64 hex chars)
var generateHEXToken = func() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return fmt.Sprintf("%x", b), err
}
