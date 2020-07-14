package internal

import (
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
)

//go:generate moq -out storage_moq_test.go . Storage
type Storage interface {
	User(email string) (storage.User, error)
	CreateUser(user storage.User) error
	UpdateUser(user storage.User) error
	CreateToken(t storage.Token) (int64, error)
	TokensByEMailAndToken(email, token string) ([]storage.Token, error)
	DeleteToken(id int64) error
}

//go:generate moq -out jwt_generator_moq_test.go . JWTGenerator
type JWTGenerator interface {
	Generate(email string, userClaims map[string]interface{}) (string, error)
}

//go:generate moq -out mailer_moq_test.go . Mailer
type Mailer interface {
	SendPasswordResetRequestEMail(recipient, passwordResetLink string) error
}

type Provider struct {
	Storage      Storage
	JWTGenerator JWTGenerator
	Mailer       Mailer
}
