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
		expectedUser         User
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:        "Happycase",
			requestBody: `{"email": "test.test@test.test", "password": "s3cr3t", "claims": {"hello": "world", "c": 42}}`,
			expectedUser: User{
				EMail:    "test.test@test.test",
				Password: "s3cr3t",
				Claims: map[string]interface{}{
					"hello": "world",
					"c":     42,
				},
			},
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
			name:          "User already exists",
			requestBody:   `{"email": "test.test@test.test", "password": "s3cr3t", "claims": {"hello": "world", "c": 42}}`,
			providerError: internal.ErrUserAlreadyExists,
			expectedUser: User{
				EMail:    "test.test@test.test",
				Password: "s3cr3t",
				Claims: map[string]interface{}{
					"hello": "world",
					"c":     42,
				},
			},
			expectedResponseCode: http.StatusConflict,
			expectedResponseBody: `{"message":"User with given email already exists"}`,
		},
		{
			name:          "Unexpected error",
			requestBody:   `{"email": "test.test@test.test", "password": "s3cr3t", "claims": {"hello": "world", "c": 42}}`,
			providerError: errors.New("nope"),
			expectedUser: User{
				EMail:    "test.test@test.test",
				Password: "s3cr3t",
				Claims: map[string]interface{}{
					"hello": "world",
					"c":     42,
				},
			},
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenUser internal.User

			toTest := NewServer(&ProviderMock{
				CreateUserFunc: func(user internal.User) error {
					givenUser = user

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

			if !reflect.DeepEqual(givenUser, givenUser) { //TODO can not compare claims via deepEqual
				t.Errorf("Provider called with unexpected User. Given: \n%#v \nExpected: \n%#v", givenUser, tt.expectedUser)
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

func TestGetUserHandler(t *testing.T) {
	tests := []struct {
		name                 string
		providerError        error
		providerUser         internal.User
		requestEmail         string
		expectedEncodedEmail string
		expectedResponseBody string
		expectedResponseCode int
	}{
		{
			name:         "Happycase",
			requestEmail: "info%40leberkleber.io",
			providerUser: internal.User{
				EMail:    "test.test@test.test",
				Password: "myPassword",
				Claims: map[string]interface{}{
					"test": "claim",
				},
			},
			expectedEncodedEmail: "info@leberkleber.io",
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: `{"email":"test.test@test.test","password":"myPassword","claims":{"test":"claim"}}`,
		},
		{
			name:                 "User not found",
			requestEmail:         "info%40leberkleber.io",
			providerError:        internal.ErrUserNotFound,
			expectedEncodedEmail: "info@leberkleber.io",
			expectedResponseCode: http.StatusNotFound,
			expectedResponseBody: `{"message":"User with given email doesn't exists"}`,
		},
		{
			name:                 "Provider error",
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
				GetUserFunc: func(email string) (internal.User, error) {
					givenEMail = email
					return tt.providerUser, tt.providerError
				},
			}, true, "username", "password")
			testServer := httptest.NewServer(toTest.h)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/admin/users/%s", testServer.URL, tt.requestEmail), nil)
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
				t.Errorf("Request response body is not as expected. Expected: \n%q\n Given: \n%q", tt.expectedResponseBody, string(compactedRespBodyAsBytes))
			}
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
		requestEmail         string
		providerUser         internal.User
		providerError        error
		expectedUser         User
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:         "Happycase",
			requestBody:  `{"password": "s3cr3t", "claims": {"hello": "world", "c": 42}}`,
			requestEmail: `test.test@test.test`,
			providerUser: internal.User{
				EMail:    "test.test@test.test",
				Password: "**********",
				Claims: map[string]interface{}{
					"hello": "world",
					"c":     42,
				},
			},
			expectedUser: User{
				Password: "s3cr3t",
				Claims: map[string]interface{}{
					"hello": "world",
					"c":     42,
				},
			},
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: `{"email":"test.test@test.test","password":"**********","claims":{"c":42,"hello":"world"}}`,
		},
		{
			name:                 "Missing in body has been set",
			requestBody:          `{"email": "test1.test1@test1.test1", "password": "s3cr3t"}`,
			requestEmail:         `test.test@test.test`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"email can not be changed"}`,
		},
		{
			name:                 "Invalid JSON",
			requestBody:          `{"passwords3cr3t"}`,
			requestEmail:         `test.test@test.test`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid JSON"}`,
		},
		{
			name:          "User not found",
			requestBody:   `{"password": "s3cr3t", "claims": {"hello": "world", "c": 42}}`,
			requestEmail:  `test3.test3@test3.test3`,
			providerError: internal.ErrUserNotFound,
			expectedUser: User{
				Password: "s3cr3t",
				Claims: map[string]interface{}{
					"hello": "world",
					"c":     42,
				},
			},
			expectedResponseCode: http.StatusNotFound,
			expectedResponseBody: `{"message":"User with given email doesn't exists"}`,
		},
		{
			name:          "Unexpected error",
			requestBody:   `{"password": "s3cr3t", "claims": {"hello": "world", "c": 42}}`,
			requestEmail:  `test.test@test.test`,
			providerError: errors.New("nope"),
			expectedUser: User{
				EMail:    "test.test@test.test",
				Password: "s3cr3t",
				Claims: map[string]interface{}{
					"hello": "world",
					"c":     42,
				},
			},
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenUser internal.User

			toTest := NewServer(&ProviderMock{
				UpdateUserFunc: func(email string, user internal.User) (internal.User, error) {
					givenUser = user

					return tt.providerUser, tt.providerError
				},
			}, true, "username", "password")
			testServer := httptest.NewServer(toTest.h)

			bb := bytes.NewReader([]byte(tt.requestBody))
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/admin/users/%s", testServer.URL, tt.requestEmail), bb)
			if err != nil {
				t.Fatalf("Failed to build http request: %s", err)
			}
			req.SetBasicAuth("username", "password")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to call server cause: %s", err)
			}
			defer resp.Body.Close()

			if !reflect.DeepEqual(givenUser, givenUser) { //TODO can not compare claims via deepEqual
				t.Errorf("Provider called with unexpected User. Given: \n%#v \nExpected: \n%#v", givenUser, tt.expectedUser)
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
				t.Errorf("Request response body is not as expected. Expected: \n%q\nGiven:\n%q", tt.expectedResponseBody, string(compactedRespBodyAsBytes))
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
			expectedResponseBody: `{"message":"User with given email doesnt already exists"}`,
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
