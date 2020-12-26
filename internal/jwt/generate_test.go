package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"testing"
	"time"
)

var jwtPubKey = `-----BEGIN PUBLIC KEY-----
MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBQSa/dFpXRqz6aQQmx6sNpxl3mn8Z
0o+qgfgOxPAPxu+JppsCGqrX/6SeUI6kz3AFVABGBU8/9Ejzt7Ty9WJt1dEB+035
03+xLnmmyaj3bEhkerr229mDgPb8uDlPEl6f/Wv+Ma/eIIloCo8WJAe8YsviImbF
hAV1NK8+62/iMCfNj30=
-----END PUBLIC KEY-----
`
var jwtPrvKey = `-----BEGIN EC PRIVATE KEY-----
MIHcAgEBBEIASzDZeTVLxcE5KTAmwrKwFjzr5cDrA+tttx9XRUz0K7AlROtj7cMG
rHu/bdKj7lc2WaW8x/EOrU/FeCcsIL5nTH+gBwYFK4EEACOhgYkDgYYABAFBJr90
WldGrPppBCbHqw2nGXeafxnSj6qB+A7E8A/G74mmmwIaqtf/pJ5QjqTPcAVUAEYF
Tz/0SPO3tPL1Ym3V0QH7TfnTf7EueabJqPdsSGR6uvbb2YOA9vy4OU8SXp/9a/4x
r94giWgKjxYkB7xiy+IiZsWEBXU0rz7rb+IwJ82PfQ==
-----END EC PRIVATE KEY-----`

func TestNewGenerator_WithoutPrivateKey(t *testing.T) {
	_, err := NewGenerator("", 4*time.Hour, "audience", "issuer", "subject")

	expectedError := errors.New("no valid private key found")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Errorf("Unexpected error. Expected: %q, Given: %q", expectedError, err)
	}
}

func TestNewGenerator_InvalidPrivateKey(t *testing.T) {
	oldX509ParseECPrivateKey := x509ParseECPrivateKey
	defer func() { x509ParseECPrivateKey = oldX509ParseECPrivateKey }()

	x509ParseECPrivateKey = func(der []byte) (*ecdsa.PrivateKey, error) {
		return nil, errors.New("errrooooooorrrr")
	}

	_, err := NewGenerator(jwtPrvKey, 4*time.Hour, "audience", "issuer", "subject")

	expectedError := errors.New("failed to parse private-key: errrooooooorrrr")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Errorf("Unexpected error. Expected: %q, Given: %q", expectedError, err)
	}
}

func TestNewGenerator_GenerateAccessToken(t *testing.T) {
	g, err := NewGenerator(jwtPrvKey, 4*time.Hour, "audience", "issuer", "subject")
	if err != nil {
		t.Fatalf("failed to crreate new generator: %s", err)
	}

	generatedJWT, err := g.GenerateAccessToken("myMailAddress", map[string]interface{}{"myCustomClaim": "mialc"})
	if err != nil {
		t.Fatalf("failed to generate jwt: %s", err)
	}
	claims := validateJWT(t, generatedJWT)
	expectedJWTAudience := "audience"
	if claims["aud"] != expectedJWTAudience {
		t.Errorf("unexpected aud-privateClaim value. Expected: %q. Given: %q", expectedJWTAudience, claims["aud"])
	}

	expectedJWTIssuer := "issuer"
	if claims["iss"] != expectedJWTIssuer {
		t.Errorf("unexpected iss-privateClaim value. Expected: %q. Given: %q", expectedJWTIssuer, claims["iss"])
	}

	expectedJWTSubject := "subject"
	if claims["sub"] != expectedJWTSubject {
		t.Errorf("unexpected sub-privateClaim value. Expected: %q. Given: %q", expectedJWTSubject, claims["sub"])
	}

	if claims["id"] == "" {
		t.Error("jwt id has not been set")
	}

	if claims["exp"] == "" {
		t.Error("jwt exp has not been set")
	}

	if claims["iat"] == "" {
		t.Error("jwt iat has not been set")
	}

	if claims["nbf"] == "" {
		t.Error("jwt nbf has not been set")
	}

	expectedCustomClaim := "mialc"
	if claims["myCustomClaim"] != expectedCustomClaim {
		t.Errorf("unexpected email-privateClaim value. Expected: %q. Given: %q", expectedCustomClaim, claims["myCustomClaim"])
	}

	expectedJWTEMail := "myMailAddress"
	if claims["email"] != expectedJWTEMail {
		t.Errorf("unexpected email-privateClaim value. Expected: %q. Given: %q", expectedJWTEMail, claims["email"])
	}
}

func TestNewGenerator_GenerateAccessToken_FailedToGenerateUUID(t *testing.T) {
	oldUUIDNewRandom := uuidNewRandom
	defer func() { uuidNewRandom = oldUUIDNewRandom }()

	uuidNewRandom = func() (uuid.UUID, error) {
		return uuid.UUID{}, errors.New("nope")
	}

	_, err := Generator{}.GenerateAccessToken("my.email.de", nil)

	expectedError := errors.New("failed to generate jwt-id: nope")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Fatalf("unexpected error. Expected: %q. Gven:: %q", expectedError, err)
	}
}

func TestNewGenerator_GenerateRefreshToken(t *testing.T) {
	g, err := NewGenerator(jwtPrvKey, 4*time.Hour, "audience", "issuer", "subject")
	if err != nil {
		t.Fatalf("failed to crreate new generator: %s", err)
	}

	generatedJWT, err := g.GenerateRefreshToken("myMailAddress")
	if err != nil {
		t.Fatalf("failed to generate jwt: %s", err)
	}
	claims := validateJWT(t, generatedJWT)
	expectedJWTAudience := "audience"
	if claims["aud"] != expectedJWTAudience {
		t.Errorf("unexpected aud-privateClaim value. Expected: %q. Given: %q", expectedJWTAudience, claims["aud"])
	}

	expectedJWTIssuer := "issuer"
	if claims["iss"] != expectedJWTIssuer {
		t.Errorf("unexpected iss-privateClaim value. Expected: %q. Given: %q", expectedJWTIssuer, claims["iss"])
	}

	expectedJWTSubject := "subject"
	if claims["sub"] != expectedJWTSubject {
		t.Errorf("unexpected sub-privateClaim value. Expected: %q. Given: %q", expectedJWTSubject, claims["sub"])
	}

	if claims["refresh"] != true {
		t.Errorf("unexpected refresh value. Expected: %t. Given: %q", true, claims["refresh"])
	}

	if claims["id"] == "" {
		t.Error("jwt id has not been set")
	}

	if claims["exp"] == "" {
		t.Error("jwt exp has not been set")
	}

	if claims["iat"] == "" {
		t.Error("jwt iat has not been set")
	}

	if claims["nbf"] == "" {
		t.Error("jwt nbf has not been set")
	}

	expectedJWTEMail := "myMailAddress"
	if claims["email"] != expectedJWTEMail {
		t.Errorf("unexpected email-privateClaim value. Expected: %q. Given: %q", expectedJWTEMail, claims["email"])
	}
}

func TestNewGenerator_GenerateRefreshToken_FailedToGenerateUUID(t *testing.T) {
	oldUUIDNewRandom := uuidNewRandom
	defer func() { uuidNewRandom = oldUUIDNewRandom }()

	uuidNewRandom = func() (uuid.UUID, error) {
		return uuid.UUID{}, errors.New("nope")
	}

	_, err := Generator{}.GenerateRefreshToken("my.email.de")

	expectedError := errors.New("failed to generate jwt-id: nope")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Fatalf("unexpected error. Expected: %q. Gven:: %q", expectedError, err)
	}
}

func validateJWT(t *testing.T, tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	pubKey, err := decodeECDSApubKey(jwtPubKey)
	if err != nil {
		t.Fatalf("Failed to parse public key: %s", err)
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return pubKey, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse jwt: %s", err)
	}

	if !token.Valid {
		t.Fatalf("Given token ist not valid. Token: %s", tokenString)
	}

	return claims
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
