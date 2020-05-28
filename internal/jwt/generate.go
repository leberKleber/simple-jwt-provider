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
		return nil, errors.New("no valid public key found")
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

func (g Generator) Generate(email string) (string, error) {
	now := nowFunc()
	jwtID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt-id: %w", err)
	}

	t := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
		//standard claims by https://tools.ietf.org/html/rfc7519#section-4.1
		"aud": g.privateClaims.audience, //Audience
		"exp": now.Add(lifeTime).Unix(), //ExpiresAt
		"jit": jwtID,                    //Id
		"iat": now.Unix(),               //IssuedAt
		"iss": g.privateClaims.issuer,   //Issuer
		"nbf": now.Unix(),               //NotBefore
		"sub": g.privateClaims.subject,  //Subject

		//public claims by https://www.iana.org/assignments/jwt/jwt.xhtml#claims
		"email": email, //EMail
	})

	signedToken, err := t.SignedString(g.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
