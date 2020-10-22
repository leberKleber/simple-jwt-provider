package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	tests := []struct {
		name                               string
		configuredUsername                 string
		configuredPassword                 string
		requestUsername                    string
		requestPassword                    string
		expectedNextHasBeenCalled          bool
		expectedUnauthorizedResponseHeader bool
	}{
		{
			name:                               "Happycase",
			configuredUsername:                 "username",
			configuredPassword:                 "password",
			requestUsername:                    "username",
			requestPassword:                    "password",
			expectedNextHasBeenCalled:          true,
			expectedUnauthorizedResponseHeader: false,
		},
		{
			name:                               "Missing auth header",
			configuredUsername:                 "username",
			configuredPassword:                 "password",
			expectedNextHasBeenCalled:          false,
			expectedUnauthorizedResponseHeader: true,
		},
		{
			name:                               "Invalid username",
			configuredUsername:                 "username",
			configuredPassword:                 "password",
			requestUsername:                    "nope",
			requestPassword:                    "password",
			expectedNextHasBeenCalled:          false,
			expectedUnauthorizedResponseHeader: true,
		},
		{
			name:                               "Invalid password",
			configuredUsername:                 "username",
			configuredPassword:                 "password",
			requestUsername:                    "username",
			requestPassword:                    "nope",
			expectedNextHasBeenCalled:          false,
			expectedUnauthorizedResponseHeader: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHasBeenCalled := false

			w := &httptest.ResponseRecorder{}
			r, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("Failed to create test request: %s", err)
			}
			if tt.requestUsername != "" || tt.requestPassword != "" {
				r.SetBasicAuth(tt.requestUsername, tt.requestPassword)
			}
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("done"))
				if err != nil {
					t.Fatalf("Could not write http request respons: %s", err)
				}
				nextHasBeenCalled = true
			})

			BasicAuth(tt.configuredUsername, tt.configuredPassword)(next).ServeHTTP(w, r)

			if tt.expectedNextHasBeenCalled != nextHasBeenCalled {
				t.Errorf("Call of next handler is not as expected. Given: %t, Exected: %t", nextHasBeenCalled, tt.expectedNextHasBeenCalled)
			}

			if tt.expectedUnauthorizedResponseHeader {
				response := w.Result()
				expectedStatusCode := http.StatusForbidden
				if response.StatusCode != expectedStatusCode {
					t.Errorf("Unexpected response code. Given: %d, Expected: %d", w.Code, expectedStatusCode)
				}

				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					t.Fatalf("Failed to read response body: %s", err)
				}

				expectedResponseBody := ""
				if string(body) != expectedResponseBody {
					t.Errorf("Unexpected response body value. \nGiven: %q \nExected: %q",
						string(body),
						expectedResponseBody)
				}
			}
		})
	}
}
