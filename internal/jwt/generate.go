package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

var nowFunc = time.Now
var lifeTime = 4 * time.Hour

type Generator struct {
	privateKey *ecdsa.PrivateKey
}

type claims struct {
	jwt.StandardClaims
	EMail string `json:"email"`
}

func NewGenerator(key string) (*Generator, error) {
	key = strings.Replace(key, `\n`, "\n", -1) //TODO fix me

	blockPrv, _ := pem.Decode([]byte(key))
	if blockPrv == nil {
		return nil, errors.New("no valid public key found")
	}

	pKey, err := x509.ParseECPrivateKey(blockPrv.Bytes)
	if err != nil {
		return nil, err
	}

	return &Generator{
		privateKey: pKey,
	}, err
}

func (g *Generator) Generate(email string) (string, error) {

	now := nowFunc()
	t := jwt.NewWithClaims(jwt.SigningMethodES512, claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(lifeTime).Unix(),
			IssuedAt:  now.Unix(),
		},
		EMail: email,
	})
	signedToken, err := t.SignedString(g.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
