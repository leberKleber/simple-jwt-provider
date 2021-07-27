package jwt

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"reflect"
	"testing"
	"time"
)

func TestProvider_Error_IsTokenValid(t *testing.T) {
	tests := []struct {
		name            string
		givenToken      string
		parseFuncErr    error
		parseFuncToken  *jwt.Token
		parseFuncClaims jwt.MapClaims
		expectedJWT     string
		expectedIsValid bool
		expectedClaims  jwt.MapClaims
		expectedErr     error
	}{
		{
			name:            "Happycase",
			givenToken:      "myToken",
			parseFuncClaims: jwt.MapClaims{"my": "claim"},
			parseFuncToken:  &jwt.Token{Valid: true},
			expectedJWT:     "myToken",
			expectedClaims:  jwt.MapClaims{"my": "claim"},
		}, {
			name:         "parse error",
			givenToken:   "myToken",
			parseFuncErr: errors.New("my error"),
			expectedErr:  errors.New("failed to parse token: my error"),
			expectedJWT:  "myToken",
		}, {
			name:            "invalid token",
			givenToken:      "myToken",
			parseFuncClaims: jwt.MapClaims{"my": "claim"},
			parseFuncToken:  &jwt.Token{Valid: false},
			expectedJWT:     "myToken",
			expectedClaims:  jwt.MapClaims{"my": "claim"},
			expectedErr:     errors.New("token is not valid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldParseFunc := parseFunc
			defer func() {
				parseFunc = oldParseFunc
			}()

			parseFunc = func(tokenString string, c jwt.Claims, _ jwt.Keyfunc) (token *jwt.Token, e error) {
				if tokenString != tt.expectedJWT {
					t.Errorf("Unexpected parseFunc>token. Expected: %s. Given: %s", tt.expectedJWT, tokenString)
				}

				cc := c.(*jwt.MapClaims)
				*cc = tt.parseFuncClaims

				return tt.parseFuncToken, tt.parseFuncErr
			}

			isValid, claims, err := Provider{privateKey: &ecdsa.PrivateKey{}}.IsTokenValid(tt.givenToken)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("Unexpected error. Expected: %q. Given: %q", tt.expectedErr, err)
			} else if err != nil {
				return
			}

			if isValid == tt.expectedIsValid {
				t.Errorf("Unexpected response isValid. Expected: %t. Given: %t", tt.expectedIsValid, isValid)
			}

			if !reflect.DeepEqual(claims, tt.expectedClaims) {
				t.Errorf("Unexpected response claims. Expected: %#v. Given: %#v", tt.expectedClaims, claims)
			}
		})
	}
}

func TestProvider_IsTokenValid(t *testing.T) {
	email := "my.mail@test.de"

	provider, err := NewProvider(jwtPrvKey, time.Minute, "audience", "issuer", "subject")
	if err != nil {
		t.Fatal("failed to create provider", err)
	}

	token, jwtID, err := provider.GenerateRefreshToken(email)
	if err != nil {
		t.Fatal("failed to generate test refresh-token", err)
	}
	if jwtID == "" {
		t.Error("generate returns no jwtID")
	}

	isValid, claims, err := provider.IsTokenValid(token)
	if err != nil {
		t.Fatal("failed to validate token", err)
	}

	if !isValid {
		t.Error("token is not valid")
	}

	claimEmail, ok := claims["email"].(string)
	if !ok {
		t.Fatalf("email is not parsable as string. Claims: %#v", claims)
	}

	if email != claimEmail {
		t.Errorf("claims>email is not as expected. Expected: %q, Given: %q", email, claimEmail)
	}
}
