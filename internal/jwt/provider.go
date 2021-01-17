package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"strings"
	"time"
)

var x509ParseECPrivateKey = x509.ParseECPrivateKey

// Provider should be created via NewProvider and creates JWTs via Generate with static and custom claims
type Provider struct {
	jwtLifetime   time.Duration
	privateKey    *ecdsa.PrivateKey
	signingMethod *jwt.SigningMethodECDSA
	privateClaims struct {
		audience string
		issuer   string
		subject  string
	}
}

// NewProvider a Provider instance with the given jwt-configuration. Before instantiation the private key will be
// checked and parsed
func NewProvider(privateKey string, jwtLifetime time.Duration, jwtAudience, jwtIssuer, jwtSubject string) (*Provider, error) {
	privateKey = strings.Replace(privateKey, `\n`, "\n", -1) //TODO fix me (needed for start via ide)
	blockPrv, _ := pem.Decode([]byte(privateKey))
	if blockPrv == nil {
		return nil, errors.New("no valid private key found")
	}

	pKey, err := x509ParseECPrivateKey(blockPrv.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private-key: %w", err)
	}

	return &Provider{
		jwtLifetime:   jwtLifetime,
		privateKey:    pKey,
		signingMethod: jwt.SigningMethodES512,
		privateClaims: struct {
			audience string
			issuer   string
			subject  string
		}{
			audience: jwtAudience,
			issuer:   jwtIssuer,
			subject:  jwtSubject,
		},
	}, err
}

var checkSigningMethodKeyFunc = func(signingMethod jwt.SigningMethod, publicKey *ecdsa.PublicKey) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		tokenSigningMethod := reflect.TypeOf(token.Method)
		expectedSigningMethod := reflect.TypeOf(signingMethod)
		if tokenSigningMethod != expectedSigningMethod {
			return nil, fmt.Errorf("unexpected signing method %q, expected: %q", tokenSigningMethod, expectedSigningMethod)
		}

		return publicKey, nil
	}
}
