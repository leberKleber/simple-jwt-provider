package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
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
			expectedResponseCode: http.StatusCreated,
		},
		{
			name:                 "Missing email",
			requestBody:          `{"password": "s3cr3t"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"email must be set"}`,
		},
		{
			name:                 "Invalid JSON",
			requestBody:          `{"passwords3cr3t"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid JSON"}`,
		},
		{
			name:                 "Missing password",
			requestBody:          `{"email": "test.test@test.test"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"password must be set"}`,
		},
		{
			name:                 "User already exists",
			requestBody:          `{"email": "test.test@test.test", "password": "s3cr3t"}`,
			providerError:        internal.ErrUserAlreadyExists,
			expectedEMail:        "test.test@test.test",
			expectedPassword:     "s3cr3t",
			expectedResponseCode: http.StatusConflict,
			expectedResponseBody: `{"message":"user with given email already exists"}`,
		},
		{
			name:                 "Unexpected error",
			requestBody:          `{"email": "test.test@test.test", "password": "s3cr3t"}`,
			providerError:        errors.New("nope"),
			expectedEMail:        "test.test@test.test",
			expectedPassword:     "s3cr3t",
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenEMail string
			var givenPassword string
			var givenClaims map[string]interface{}

			toTest := NewServer(&ProviderMock{
				CreateUserFunc: func(email, password string, claims map[string]interface{}) error {
					givenEMail = email
					givenPassword = password
					givenClaims = claims

					return tt.providerError
				},
			}, true, "username", "password")
			testServer := httptest.NewServer(toTest.h)

			bb := bytes.NewReader([]byte(tt.requestBody))
			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/v1/admin/users", bb)
			if err != nil {
				t.Fatalf("Failed to build http request: %s", err)
			}
			req.SetBasicAuth("username", "password")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to call server cause: %s", err)
			}
			defer resp.Body.Close()

			if givenEMail != tt.expectedEMail {
				t.Errorf("Provider called with unexpected email. Given: %q, Expected: %q", givenEMail, tt.expectedEMail)
			}

			if givenPassword != tt.expectedPassword {
				t.Errorf("Provider called with unexpected password. Given: %q, Expected: %q", givenPassword, tt.expectedPassword)
			}

			if !reflect.DeepEqual(givenClaims, givenClaims) { //TODO check claims
				t.Errorf("Request respond with unexpected claims code. Expected: %d, Given: %d", tt.expectedResponseCode, resp.StatusCode)
			}

			if resp.StatusCode != tt.expectedResponseCode {
				t.Errorf("Request respond with unexpected status code. Expected: %d, Given: %d", tt.expectedResponseCode, resp.StatusCode)
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
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

func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name                 string
		providerError        error
		requestEmail         string
		expectedEncodedEmail string
		expectedResponseBody string
		expectedResponseCode int
	}{
		{
			name:                 "Happycase",
			requestEmail:         "info%40leberkleber.io",
			expectedEncodedEmail: "info@leberkleber.io",
			expectedResponseCode: http.StatusNoContent,
		},
		{
			name:                 "User not found",
			requestEmail:         "info%40leberkleber.io",
			providerError:        internal.ErrUserNotFound,
			expectedEncodedEmail: "info@leberkleber.io",
			expectedResponseCode: http.StatusNotFound,
			expectedResponseBody: `{"message":"user with given email doesnt already exists"}`,
		},
		{
			name:                 "Error while deletion",
			requestEmail:         "info%40leberkleber.io",
			providerError:        errors.New("nope"),
			expectedEncodedEmail: "info@leberkleber.io",
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenEMail string

			toTest := NewServer(&ProviderMock{
				DeleteUserFunc: func(email string) error {
					givenEMail = email
					return tt.providerError
				},
			}, true, "username", "password")
			testServer := httptest.NewServer(toTest.h)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/admin/users/%s", testServer.URL, tt.requestEmail), nil)
			if err != nil {
				t.Fatalf("Failed to build http request: %s", err)
			}
			req.SetBasicAuth("username", "password")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to call server cause: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedResponseCode {
				t.Errorf("Request respond with unexpected status code. Expected: %d, Given: %d", tt.expectedResponseCode, resp.StatusCode)
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
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

			if tt.expectedEncodedEmail != givenEMail {
				t.Errorf("Unexpected delete email. Expected: %q, Given: %q", tt.expectedEncodedEmail, givenEMail)
			}

			if !bytes.Equal(compactedRespBodyAsBytes, []byte(tt.expectedResponseBody)) {
				t.Errorf("Request response body is not as expected. Expected: %q, Given: %q", tt.expectedResponseBody, string(compactedRespBodyAsBytes))
			}
		})
	}
}
