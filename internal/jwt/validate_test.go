package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"testing"
)

func TestProvider_IsTokenValid(t *testing.T) {
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

			isValid, claims, err := Provider{}.IsTokenValid(tt.givenToken)
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
