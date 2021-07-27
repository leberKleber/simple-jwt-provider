package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestGenerator_GenerateAccessToken(t *testing.T) {
	g, err := NewProvider(jwtPrvKey, 4*time.Hour, "audience", "issuer", "subject")
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

func TestGenerator_GenerateAccessToken_FailedToGenerateUUID(t *testing.T) {
	oldUUIDNewRandom := uuidNewRandom
	defer func() { uuidNewRandom = oldUUIDNewRandom }()

	uuidNewRandom = func() (uuid.UUID, error) {
		return uuid.UUID{}, errors.New("nope")
	}

	_, err := Provider{}.GenerateAccessToken("my.email.de", nil)

	expectedError := errors.New("failed to generate jwt-id: nope")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Fatalf("unexpected error. Expected: %q. Gven:: %q", expectedError, err)
	}
}

func TestGenerator_GenerateAccessToken_FailedToSignToken(t *testing.T) {
	p, err := NewProvider(jwtPrvKey, 4*time.Hour, "audience", "issuer", "subject")
	if err != nil {
		t.Fatalf("failed to crreate new generator: %s", err)
	}

	_, err = p.GenerateAccessToken("my.email.de", map[string]interface{}{
		"unmarshableClaim": make(chan string),
	})

	expectedError := errors.New("failed to sign access-token: json: unsupported type: chan string")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Fatalf("unexpected error. Expected: %q. Gven:: %q", expectedError, err)
	}
}

func TestGenerator_GenerateRefreshToken(t *testing.T) {
	g, err := NewProvider(jwtPrvKey, 4*time.Hour, "audience", "issuer", "subject")
	if err != nil {
		t.Fatalf("failed to crreate new generator: %s", err)
	}

	generatedJWT, jwtID, err := g.GenerateRefreshToken("myMailAddress")
	if err != nil {
		t.Fatalf("failed to generate jwt: %s", err)
	}
	if jwtID == "" {
		t.Error("generate returns no jwtID")
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

	expectedJWTEMail := "myMailAddress"
	if claims["email"] != expectedJWTEMail {
		t.Errorf("unexpected email-privateClaim value. Expected: %q. Given: %q", expectedJWTEMail, claims["email"])
	}
}

func TestGenerator_GenerateRefreshToken_FailedToGenerateUUID(t *testing.T) {
	oldUUIDNewRandom := uuidNewRandom
	defer func() { uuidNewRandom = oldUUIDNewRandom }()

	uuidNewRandom = func() (uuid.UUID, error) {
		return uuid.UUID{}, errors.New("nope")
	}

	_, _, err := Provider{}.GenerateRefreshToken("my.email.de")
	expectedError := errors.New("failed to generate jwt-id: nope")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Fatalf("unexpected error. Expected: %q. Gven:: %q", expectedError, err)
	}
}

func validateJWT(t *testing.T, tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	pubKey, err := decodeECDSAPubKey(jwtPubKey)
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

func decodeECDSAPubKey(pemEncodedPub string) (*ecdsa.PublicKey, error) {
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
