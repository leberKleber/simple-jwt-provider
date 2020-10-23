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

var nowFunc = time.Now
var lifeTime = 4 * time.Hour

type Generator struct {
	privateKey    *ecdsa.PrivateKey
	privateClaims struct {
		audience string
		issuer   string
		subject  string
	}
}

func NewGenerator(privateKey, jwtAudience, jwtIssuer, jwtSubject string) (*Generator, error) {
	privateKey = strings.Replace(privateKey, `\n`, "\n", -1) //TODO fix me (needed for start via ide)
	blockPrv, _ := pem.Decode([]byte(privateKey))
	if blockPrv == nil {
		return nil, errors.New("no valid private key found")
	}

	pKey, err := x509.ParseECPrivateKey(blockPrv.Bytes)
	if err != nil {
		return nil, err
	}

	return &Generator{
		privateKey: pKey,
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

// Generate generates a valid jwt based on the Generator.privateKey. The jwt is issued to the given email and enriched
// with the given claims.
// 'userClaims' can be contain all json compatible types
func (g Generator) Generate(email string, userClaims map[string]interface{}) (string, error) {
	now := nowFunc()
	jwtID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt-id: %w", err)
	}

	claims := jwt.MapClaims{}
	if userClaims != nil {
		claims = userClaims
	}

	//standard claims by https://tools.ietf.org/html/rfc7519#section-4.1
	claims["aud"] = g.privateClaims.audience //Audience
	claims["exp"] = now.Add(lifeTime).Unix() //ExpiresAt
	claims["jit"] = jwtID                    //Id
	claims["iat"] = now.Unix()               //IssuedAt
	claims["iss"] = g.privateClaims.issuer   //Issuer
	claims["nbf"] = now.Unix()               //NotBefore
	claims["sub"] = g.privateClaims.subject  //Subject

	//public claims by https://www.iana.org/assignments/jwt/jwt.xhtml#claims
	claims["email"] = email //Recipient

	t := jwt.NewWithClaims(jwt.SigningMethodES512, claims)

	signedToken, err := t.SignedString(g.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
