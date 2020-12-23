package internal

import (
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
)

// Storage encapsulates storage.Storage to generate mocks
//go:generate moq -out storage_moq_test.go . Storage
type Storage interface {
	User(email string) (storage.User, error)
	CreateUser(user storage.User) error
	UpdateUser(user storage.User) error
	DeleteUser(email string) error
	CreateToken(t storage.Token) (int64, error)
	TokensByEMailAndToken(email, token string) ([]storage.Token, error)
	DeleteToken(id int64) error
}

// JWTGenerator encapsulates jwt.Generator to generate mocks
//go:generate moq -out jwt_generator_moq_test.go . JWTGenerator
type JWTGenerator interface {
	GenerateAccessToken(email string, userClaims map[string]interface{}) (string, error)
	GenerateRefreshToken(email string) (string, error)
}

// Mailer encapsulates mailer.Mailer to generate mocks
//go:generate moq -out mailer_moq_test.go . Mailer
type Mailer interface {
	SendPasswordResetRequestEMail(recipient, passwordResetToken string, claims map[string]interface{}) error
}

// Provider provides all necessary interfaces for use in internal
type Provider struct {
	Storage      Storage
	JWTGenerator JWTGenerator
	Mailer       Mailer
}
