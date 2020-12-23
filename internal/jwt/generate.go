package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"strings"
	"time"
)

var timeNow = time.Now
var uuidNewRandom = uuid.NewRandom
var x509ParseECPrivateKey = x509.ParseECPrivateKey

// Generator should be created via NewGenerator and creates JWTs via Generate with static and custom claims
type Generator struct {
	jwtLifetime   time.Duration
	privateKey    *ecdsa.PrivateKey
	privateClaims struct {
		audience string
		issuer   string
		subject  string
	}
}

// NewGenerator a Generator instance with the given jwt-configuration. Before instantiation the private key will be
// checked and parsed
func NewGenerator(privateKey string, jwtLifetime time.Duration, jwtAudience, jwtIssuer, jwtSubject string) (*Generator, error) {
	privateKey = strings.Replace(privateKey, `\n`, "\n", -1) //TODO fix me (needed for start via ide)
	blockPrv, _ := pem.Decode([]byte(privateKey))
	if blockPrv == nil {
		return nil, errors.New("no valid private key found")
	}

	pKey, err := x509ParseECPrivateKey(blockPrv.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private-key: %w", err)
	}

	return &Generator{
		jwtLifetime: jwtLifetime,
		privateKey:  pKey,
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

// GenerateAccessToken generates a valid access-jwt based on the Generator.privateKey. The jwt is issued to the given email and enriched
// with the given claims.
// 'userClaims' can be contain all json compatible types
func (g Generator) GenerateAccessToken(email string, userClaims map[string]interface{}) (string, error) {
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
	claims["aud"] = g.privateClaims.audience      //Audience
	claims["exp"] = now.Add(g.jwtLifetime).Unix() //ExpiresAt
	claims["jit"] = jwtID                         //Id
	claims["iat"] = now.Unix()                    //IssuedAt
	claims["iss"] = g.privateClaims.issuer        //Issuer
	claims["nbf"] = now.Unix()                    //NotBefore
	claims["sub"] = g.privateClaims.subject       //Subject

	// public claims by https://www.iana.org/assignments/jwt/jwt.xhtml#claims
	claims["email"] = email // Preferred e-mail address

	token, err := jwt.NewWithClaims(jwt.SigningMethodES512, claims).SignedString(g.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return token, nil
}

func (g Generator) GenerateRefreshToken(email string) (string, error) {
	return "<refresh_token>", nil //TODO generate
}
