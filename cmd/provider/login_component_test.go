// +build component

package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	email := "info@leberkleber.io"
	password := "s3cr3t"

	createUser(t, email, password)
	token := loginUser(t, email, password)

	validateJWT(t, token)
}

func validateJWT(t *testing.T, tokenString string) {
	pubKey, err := decodeECDSApubKey(`-----BEGIN PUBLIC KEY-----
MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBQSa/dFpXRqz6aQQmx6sNpxl3mn8Z
0o+qgfgOxPAPxu+JppsCGqrX/6SeUI6kz3AFVABGBU8/9Ejzt7Ty9WJt1dEB+035
03+xLnmmyaj3bEhkerr229mDgPb8uDlPEl6f/Wv+Ma/eIIloCo8WJAe8YsviImbF
hAV1NK8+62/iMCfNj30=
-----END PUBLIC KEY-----
`)
	if err != nil {
		t.Fatalf("Failed to parse public key: %s", err)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return pubKey, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse jwt: %s", err)
	}

	if !token.Valid {
		t.Fatalf("Given token ist not valid. Token: %s", tokenString)
	}
}

func decodeECDSApubKey(pemEncodedPub string) (*ecdsa.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	if blockPub == nil {
		return nil, errors.New("No valid public key found")
	}
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, err
	}
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey, nil
}

func loginUser(t *testing.T, email, password string) string {
	resp, err := http.Post(
		"http://simple-jwt-provider/v1/auth/login",
		"application/json",
		bytes.NewReader([]byte(fmt.Sprintf(`{"email": %q, "password": %q}`, email, password))),
	)
	if err != nil {
		t.Fatalf("Failed to login cause: %s", err)
	}

	responseBody := struct {
		AccessToken  string `json:"access_token"`
		ErrorMessage string `json:"message"`
	}{}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Invalid response status code. Expected: %d, Given: %d, Body: %s", http.StatusOK, resp.StatusCode, responseBody.ErrorMessage)
	}

	return responseBody.AccessToken
}
