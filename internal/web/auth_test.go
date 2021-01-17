package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/leberKleber/simple-jwt-provider/internal"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
		providerAccessToken  string
		providerRefreshToken string
		providerError        error
		expectedEMail        string
		expectedPassword     string
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:                 "Happycase",
			requestBody:          `{"email": "test.test@test.test", "password": "s3cr3t"}`,
			expectedEMail:        "test.test@test.test",
			expectedPassword:     "s3cr3t",
			providerAccessToken:  "myAccessJWT",
			providerRefreshToken: "myRefreshJWT",
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: `{"access_token":"myAccessJWT","refresh_token":"myRefreshJWT"}`,
		},
		{
			name:                 "Invalid JSON",
			requestBody:          `{"password s3cr3t"}`,
			providerAccessToken:  "myAccessJWT",
			providerRefreshToken: "myRefreshJWT",
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid JSON"}`,
		},
		{
			name:                 "Missing Recipient",
			requestBody:          `{"password": "s3cr3t"}`,
			providerAccessToken:  "myNewJWT",
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"email must be set"}`,
		},
		{
			name:                 "Missing Password",
			requestBody:          `{"email": "test.test@test.test"}`,
			providerAccessToken:  "myNewJWT",
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"password must be set"}`,
		},
		{
			name:                 "Incorrect Password",
			requestBody:          `{"email": "test.test@test.test", "password": "n0p3"}`,
			providerAccessToken:  "myNewJWT",
			providerError:        internal.ErrIncorrectPassword,
			expectedEMail:        "test.test@test.test",
			expectedPassword:     "n0p3",
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid credentials"}`,
		},
		{
			name:                 "User not found",
			requestBody:          `{"email": "not.found@test.test", "password": "s3cr3t"}`,
			providerError:        internal.ErrUserNotFound,
			expectedEMail:        "not.found@test.test",
			expectedPassword:     "s3cr3t",
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid credentials"}`,
		},
		{
			name:                 "Unexpected error",
			requestBody:          `{"email": "not.found@test.test", "password": "s3cr3t"}`,
			providerError:        errors.New("nope"),
			expectedEMail:        "not.found@test.test",
			expectedPassword:     "s3cr3t",
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenEMail, givenPassword string

			toTest := NewServer(&ProviderMock{
				LoginFunc: func(email string, password string) (string, string, error) {
					givenEMail = email
					givenPassword = password

					return tt.providerAccessToken, tt.providerRefreshToken, tt.providerError
				},
			}, false, "", "")
			testServer := httptest.NewServer(toTest.h)

			bb := bytes.NewReader([]byte(tt.requestBody))
			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/v1/auth/login", bb)
			if err != nil {
				t.Fatalf("Failed to build http request: %s", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to call server cause: %s", err)
			}
			defer resp.Body.Close()

			expectedContentType := "application/json"
			givenContentType := resp.Header.Get("Content-Type")
			if expectedContentType != givenContentType {
				t.Errorf("Unexpected response content-type. Given: %q, Expected: %q", givenContentType, expectedContentType)
			}

			if resp.StatusCode != tt.expectedResponseCode {
				t.Errorf("Request respond with unexpected status code. Expected: %d, Given: %d", tt.expectedResponseCode, resp.StatusCode)
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
			}

			if givenEMail != tt.expectedEMail {
				t.Errorf("Provider called with unexpected email. Given: %q, Expected: %q", givenEMail, tt.expectedEMail)
			}

			if givenPassword != tt.expectedPassword {
				t.Errorf("Provider called with unexpected password. Given: %q, Expected: %q", givenPassword, tt.expectedPassword)
			}

			var compactedRespBodyAsBytes []byte
			if resp.ContentLength > 0 {
				compactedRespBody := &bytes.Buffer{}
				err = json.Compact(compactedRespBody, respBody)
				if err != nil {
					t.Fatalf("Failed to compact json: %s", err)
				}

				compactedRespBodyAsBytes = compactedRespBody.Bytes()
			}

			if !bytes.Equal(compactedRespBodyAsBytes, []byte(tt.expectedResponseBody)) {
				t.Errorf("Request response body is not as expected. Expected: \n%q, \nGiven: \n%q", tt.expectedResponseBody, string(compactedRespBodyAsBytes))
			}
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
		providerAccessToken  string
		providerRefreshToken string
		providerError        error
		expectedRefreshToken string
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:                 "Happycase",
			requestBody:          `{"refresh_token": "myOldRefreshToken"}`,
			expectedRefreshToken: "myOldRefreshToken",
			providerAccessToken:  "myAccessJWT",
			providerRefreshToken: "myRefreshJWT",
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: `{"access_token":"myAccessJWT","refresh_token":"myRefreshJWT"}`,
		},
		{
			name:                 "Invalid JSON",
			requestBody:          `{"refresh_token myOldRefreshToken"}`,
			providerAccessToken:  "myAccessJWT",
			providerRefreshToken: "myRefreshJWT",
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid JSON"}`,
		},
		{
			name:                 "Missing RefreshToken",
			requestBody:          `{}`,
			providerAccessToken:  "myNewJWT",
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"refresh_token must be set"}`,
		},
		{
			name:                 "User not found",
			requestBody:          `{"refresh_token": "myOldRefreshToken"}`,
			providerError:        internal.ErrUserNotFound,
			expectedRefreshToken: "myOldRefreshToken",
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid refresh-token and/or email"}`,
		},
		{
			name:                 "Invalid token",
			requestBody:          `{"refresh_token": "myOldRefreshToken"}`,
			providerError:        internal.ErrInvalidToken,
			expectedRefreshToken: "myOldRefreshToken",
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid refresh-token and/or email"}`,
		},
		{
			name:                 "Unexpected error",
			requestBody:          `{"refresh_token": "myOldRefreshToken"}`,
			providerError:        errors.New("nope"),
			expectedRefreshToken: "myOldRefreshToken",
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenRefreshToken string

			toTest := NewServer(&ProviderMock{
				RefreshFunc: func(refreshToken string) (string, string, error) {
					givenRefreshToken = refreshToken

					return tt.providerAccessToken, tt.providerRefreshToken, tt.providerError
				},
			}, false, "", "")
			testServer := httptest.NewServer(toTest.h)

			bb := bytes.NewReader([]byte(tt.requestBody))
			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/v1/auth/refresh", bb)
			if err != nil {
				t.Fatalf("Failed to build http request: %s", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to call server cause: %s", err)
			}
			defer resp.Body.Close()

			expectedContentType := "application/json"
			givenContentType := resp.Header.Get("Content-Type")
			if expectedContentType != givenContentType {
				t.Errorf("Unexpected response content-type. Given: %q, Expected: %q", givenContentType, expectedContentType)
			}

			if resp.StatusCode != tt.expectedResponseCode {
				t.Errorf("Request respond with unexpected status code. Expected: %d, Given: %d", tt.expectedResponseCode, resp.StatusCode)
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
			}

			if givenRefreshToken != tt.expectedRefreshToken {
				t.Errorf("Provider called with unexpected refresh-token. Given: %q, Expected: %q", givenRefreshToken, tt.expectedRefreshToken)
			}

			var compactedRespBodyAsBytes []byte
			if resp.ContentLength > 0 {
				compactedRespBody := &bytes.Buffer{}
				err = json.Compact(compactedRespBody, respBody)
				if err != nil {
					t.Fatalf("Failed to compact json: %s", err)
				}

				compactedRespBodyAsBytes = compactedRespBody.Bytes()
			}

			if !bytes.Equal(compactedRespBodyAsBytes, []byte(tt.expectedResponseBody)) {
				t.Errorf("Request response body is not as expected. Expected: \n%q, \nGiven: \n%q", tt.expectedResponseBody, string(compactedRespBodyAsBytes))
			}
		})
	}
}

func TestPasswordResetRequestHandler(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
		providerError        error
		expectedEMail        string
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:                 "Happycase",
			requestBody:          `{"email": "test.test@test.test"}`,
			expectedEMail:        "test.test@test.test",
			expectedResponseCode: http.StatusCreated,
		},
		{
			name:                 "Invalid JSON",
			requestBody:          `{"email test.test@test.test"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid JSON"}`,
		},
		{
			name:                 "Missing email",
			requestBody:          `{}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"email must be set"}`,
		},
		{
			name:                 "User not found",
			requestBody:          `{"email": "test.test@test.test"}`,
			providerError:        internal.ErrUserNotFound,
			expectedEMail:        "test.test@test.test",
			expectedResponseCode: http.StatusCreated,
		},
		{
			name:                 "Unexpected error",
			requestBody:          `{"email": "test.test@test.test"}`,
			providerError:        errors.New("error no 42"),
			expectedEMail:        "test.test@test.test",
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenEMail string

			toTest := NewServer(&ProviderMock{
				CreatePasswordResetRequestFunc: func(email string) error {
					givenEMail = email
					return tt.providerError
				},
			}, false, "", "")
			testServer := httptest.NewServer(toTest.h)

			bb := bytes.NewReader([]byte(tt.requestBody))
			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/v1/auth/password-reset-request", bb)
			if err != nil {
				t.Fatalf("Failed to build http request: %s", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to call server cause: %s", err)
			}
			defer resp.Body.Close()

			expectedContentType := "application/json"
			givenContentType := resp.Header.Get("Content-Type")
			if expectedContentType != givenContentType {
				t.Errorf("Unexpected response content-type. Given: %q, Expected: %q", givenContentType, expectedContentType)
			}

			if resp.StatusCode != tt.expectedResponseCode {
				t.Errorf("Request respond with unexpected status code. Expected: %d, Given: %d", tt.expectedResponseCode, resp.StatusCode)
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
			}

			if givenEMail != tt.expectedEMail {
				t.Errorf("Provider called with unexpected email. Given: %q, Expected: %q", givenEMail, tt.expectedEMail)
			}

			var compactedRespBodyAsBytes []byte
			if resp.ContentLength > 0 {
				compactedRespBody := &bytes.Buffer{}
				err = json.Compact(compactedRespBody, respBody)
				if err != nil {
					t.Fatalf("Failed to compact json: %s", err)
				}

				compactedRespBodyAsBytes = compactedRespBody.Bytes()
			}

			if !bytes.Equal(compactedRespBodyAsBytes, []byte(tt.expectedResponseBody)) {
				t.Errorf("Request response body is not as expected. Expected: %q, Given: %q", tt.expectedResponseBody, string(compactedRespBodyAsBytes))
			}
		})
	}
}

func TestPasswordResetHandler(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
		providerError        error
		expectedEMail        string
		expectedResetToken   string
		expectedPassword     string
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:                 "Happycase",
			requestBody:          `{"email":"test.test@test.test","password": "new_s3cr3t","reset_token": "myResetToken"}`,
			expectedEMail:        "test.test@test.test",
			expectedPassword:     "new_s3cr3t",
			expectedResetToken:   "myResetToken",
			expectedResponseCode: http.StatusNoContent,
		},
		{
			name:                 "Invalid JSON",
			requestBody:          `{"email test.test@test.test}"`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid JSON"}`,
		},
		{
			name:                 "Missing email",
			requestBody:          `{"password": "new_s3cr3t","reset_token": "myResetToken"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"email must be set"}`,
		},
		{
			name:                 "Missing reset-token",
			requestBody:          `{"email":"test.test@test.test","password": "new_s3cr3t"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"reset-token must be set"}`,
		},
		{
			name:                 "Missing password",
			requestBody:          `{"email":"test.test@test.test","reset_token": "myResetToken"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"password must be set"}`,
		},
		{
			name:                 "Invalid token",
			requestBody:          `{"email":"test.test@test.test","password": "new_s3cr3t","reset_token": "invalidResetToken"}`,
			providerError:        internal.ErrNoValidTokenFound,
			expectedEMail:        "test.test@test.test",
			expectedPassword:     "new_s3cr3t",
			expectedResetToken:   "invalidResetToken",
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"reset-token is invalid or token email combination is not correct"}`,
		},
		{
			name:                 "Unexpected error",
			requestBody:          `{"email":"test.test@test.test","password": "new_s3cr3t","reset_token": "myResetToken"}`,
			providerError:        errors.New("computer says nooooo"),
			expectedEMail:        "test.test@test.test",
			expectedPassword:     "new_s3cr3t",
			expectedResetToken:   "myResetToken",
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenEMail, givenResetToken, givenPassword string

			toTest := NewServer(&ProviderMock{
				ResetPasswordFunc: func(email string, resetToken string, password string) error {
					givenEMail = email
					givenResetToken = resetToken
					givenPassword = password
					return tt.providerError
				},
			}, false, "", "")
			testServer := httptest.NewServer(toTest.h)

			bb := bytes.NewReader([]byte(tt.requestBody))
			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/v1/auth/password-reset", bb)
			if err != nil {
				t.Fatalf("Failed to build http request: %s", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to call server cause: %s", err)
			}
			defer resp.Body.Close()

			expectedContentType := "application/json"
			givenContentType := resp.Header.Get("Content-Type")
			if expectedContentType != givenContentType {
				t.Errorf("Unexpected response content-type. Given: %q, Expected: %q", givenContentType, expectedContentType)
			}

			if resp.StatusCode != tt.expectedResponseCode {
				t.Errorf("Request respond with unexpected status code. Expected: %d, Given: %d", tt.expectedResponseCode, resp.StatusCode)
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
			}

			if givenEMail != tt.expectedEMail {
				t.Errorf("Provider called with unexpected email. Given: %q, Expected: %q", givenEMail, tt.expectedEMail)
			}

			if givenResetToken != tt.expectedResetToken {
				t.Errorf("Provider called with unexpected reset-token. Given: %q, Expected: %q", givenResetToken, tt.expectedResetToken)
			}

			if givenPassword != tt.expectedPassword {
				t.Errorf("Provider called with unexpected password. Given: %q, Expected: %q", givenPassword, tt.expectedPassword)
			}

			var compactedRespBodyAsBytes []byte
			if resp.ContentLength > 0 {
				compactedRespBody := &bytes.Buffer{}
				err = json.Compact(compactedRespBody, respBody)
				if err != nil {
					t.Fatalf("Failed to compact json: %s", err)
				}

				compactedRespBodyAsBytes = compactedRespBody.Bytes()
			}

			if !bytes.Equal(compactedRespBodyAsBytes, []byte(tt.expectedResponseBody)) {
				t.Errorf("Request response body is not as expected. Expected: %q, Given: %q", tt.expectedResponseBody, string(compactedRespBodyAsBytes))
			}
		})
	}
}
