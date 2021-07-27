package jwt

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"math/big"
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
	_, err := NewProvider("", 4*time.Hour, "audience", "issuer", "subject")

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

	_, err := NewProvider(jwtPrvKey, 4*time.Hour, "audience", "issuer", "subject")

	expectedError := errors.New("failed to parse private-key: errrooooooorrrr")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Errorf("Unexpected error. Expected: %q, Given: %q", expectedError, err)
	}
}

func TestCheckSigningMethodKeyFunc(t *testing.T) {
	tests := []struct {
		name               string
		givenSigningMethod jwt.SigningMethod
		givenPublicKey     *ecdsa.PublicKey
		givenToken         *jwt.Token
		expectedResponse   interface{}
		expectedErr        error
	}{
		{
			name:               "Happycase",
			givenSigningMethod: jwt.SigningMethodES512,
			givenPublicKey:     &ecdsa.PublicKey{X: big.NewInt(555), Y: big.NewInt(666)},
			givenToken: &jwt.Token{
				Method: jwt.SigningMethodES512,
			},
			expectedResponse: &ecdsa.PublicKey{X: big.NewInt(555), Y: big.NewInt(666)},
		}, {
			name:               "Unexpected signing method",
			givenSigningMethod: jwt.SigningMethodES512,
			givenPublicKey:     &ecdsa.PublicKey{X: big.NewInt(555), Y: big.NewInt(666)},
			givenToken: &jwt.Token{
				Method: jwt.SigningMethodPS256,
			},
			expectedErr:      errors.New("unexpected signing method \"*jwt.SigningMethodRSAPSS\", expected: \"*jwt.SigningMethodECDSA\""),
			expectedResponse: &ecdsa.PublicKey{X: big.NewInt(555), Y: big.NewInt(666)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtKeyFunc := checkSigningMethodKeyFunc(tt.givenSigningMethod, tt.givenPublicKey)

			resp, err := jwtKeyFunc(tt.givenToken)
			expectedResponseAsString := fmt.Sprint(tt.expectedResponse)
			respAsString := fmt.Sprint(resp)

			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Errorf("Unexpected error. \nExpected: %q\nGiven:\n%q", tt.expectedErr, err)
			} else if err != nil {
				return
			}

			if expectedResponseAsString != respAsString {
				t.Errorf("unexpected response. Given: %q, Expected: %q", respAsString, expectedResponseAsString)
			}
		})
	}
}
