package internal

import (
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
)

type Storage interface {
	User(email string) (storage.User, error)
	CreateUser(user storage.User) error
}

type JWTGenerator interface {
	Generate(email string) (string, error)
}

type Provider struct {
	Storage      Storage
	JWTGenerator JWTGenerator
}
