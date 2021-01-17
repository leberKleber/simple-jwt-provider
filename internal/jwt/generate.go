package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

var timeNow = time.Now
var uuidNewRandom = uuid.NewRandom

const refreshTokenLifetime = 7 * 24 * time.Hour

// GenerateAccessToken generates a valid access-jwt based on the Provider.privateKey. The jwt is issued to the given email and enriched
// with the given claims.
// 'userClaims' can be contain all json compatible types
func (p Provider) GenerateAccessToken(email string, userClaims map[string]interface{}) (string, error) {
	now := timeNow()
	jwtID, err := uuidNewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt-id: %w", err)
	}

	claims := jwt.MapClaims{}
	if userClaims != nil {
		claims = userClaims
	}

	// standard claims by https://tools.ietf.org/html/rfc7519#section-4.1
	claims["aud"] = p.privateClaims.audience      //Audience
	claims["exp"] = now.Add(p.jwtLifetime).Unix() //ExpiresAt
	claims["jit"] = jwtID.String()                //Id
	claims["iat"] = now.Unix()                    //IssuedAt
	claims["iss"] = p.privateClaims.issuer        //Issuer
	claims["nbf"] = now.Unix()                    //NotBefore
	claims["sub"] = p.privateClaims.subject       //Subject

	// public claims by https://www.iana.org/assignments/jwt/jwt.xhtml#claims
	claims["email"] = email // Preferred e-mail address

	token, err := jwt.NewWithClaims(p.signingMethod, claims).SignedString(p.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access-token: %w", err)
	}

	return token, nil
}

// GenerateRefreshToken generates a valid refresh-jwt based on the Provider.privateKey. The jwt is issued to the given email.
func (p Provider) GenerateRefreshToken(email string) (string, string, error) {
	now := timeNow()
	jwtID, err := uuidNewRandom()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate jwt-id: %w", err)
	}

	claims := jwt.MapClaims{}

	// standard claims by https://tools.ietf.org/html/rfc7519#section-4.1
	claims["aud"] = p.privateClaims.audience             //Audience
	claims["exp"] = now.Add(refreshTokenLifetime).Unix() //ExpiresAt
	claims["jit"] = jwtID.String()                       //Id
	claims["iat"] = now.Unix()                           //IssuedAt
	claims["iss"] = p.privateClaims.issuer               //Issuer
	claims["nbf"] = now.Unix()                           //NotBefore
	claims["sub"] = p.privateClaims.subject              //Subject

	// public claims by https://www.iana.org/assignments/jwt/jwt.xhtml#claims
	claims["email"] = email // Preferred e-mail address

	token, err := jwt.NewWithClaims(p.signingMethod, claims).SignedString(p.privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh-token: %w", err)
	}

	return token, jwtID.String(), nil
}
