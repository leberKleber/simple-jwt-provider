package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

var parseFunc = jwt.ParseWithClaims

// IsTokenValid validates the given token with the in NewProvider configured privateKey.PublicKeys and return
// isValid indicator, token-claims (when token is valid) and an error when present
// return
func (p Provider) IsTokenValid(tokenAsString string) (isValid bool, claims jwt.MapClaims, err error) {
	token, err := parseFunc(tokenAsString, &claims, checkSigningMethodKeyFunc(p.signingMethod, &p.privateKey.PublicKey))
	if err != nil {
		return false, nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return false, nil, errors.New("token is not valid")
	}

	return true, claims, nil
}
