package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var nowFunc = time.Now
var lifeTime = 4 * time.Hour

type Generator struct {
	PrivateKey string
}

type claims struct {
	jwt.StandardClaims
	EMail string `json:"email"`
}

func (f *Generator) Generate(email string) string {
	now := nowFunc()
	return jwt.NewWithClaims(jwt.SigningMethodES512, claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(lifeTime).Unix(),
			IssuedAt:  now.Unix(),
		},
		EMail: email,
	}).Raw
}
