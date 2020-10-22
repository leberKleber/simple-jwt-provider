package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	expectedResponseCode := http.StatusForbidden
	expectedResponseBody := `{"message":"forbidden"}`

	toTest := NewServer(nil, true, "un", "pw")
	testServer := httptest.NewServer(toTest.h)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/v1/admin/users", nil)
	if err != nil {
		t.Fatalf("Failed to build http request: %s", err)
	}

	req.SetBasicAuth("invalid", "invalid")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to call server cause: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedResponseCode {
		t.Errorf("Request respond with unexpected status code. Expected: %d, Given: %d", expectedResponseCode, resp.StatusCode)
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

	if !bytes.Equal(compactedRespBodyAsBytes, []byte(expectedResponseBody)) {
		t.Errorf("Request response body is not as expected. Expected: %q, Given: %q", expectedResponseBody, string(compactedRespBodyAsBytes))
	}
}

func TestNotFoundHandler(t *testing.T) {
	expectedResponseCode := http.StatusNotFound
	expectedResponseBody := `{"message":"endpoint not found"}`

	toTest := NewServer(nil, false, "", "")
	testServer := httptest.NewServer(toTest.h)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/unexpected/endpoint", nil)
	if err != nil {
		t.Fatalf("Failed to build http request: %s", err)
	}

	req.SetBasicAuth("invalid", "invalid")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to call server cause: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedResponseCode {
		t.Errorf("Request respond with unexpected status code. Expected: %d, Given: %d", expectedResponseCode, resp.StatusCode)
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

	if !bytes.Equal(compactedRespBodyAsBytes, []byte(expectedResponseBody)) {
		t.Errorf("Request response body is not as expected. Expected: %q, Given: %q", expectedResponseBody, string(compactedRespBodyAsBytes))
	}
}

func TestMethodNotAllowedHandler(t *testing.T) {
	expectedResponseCode := http.StatusMethodNotAllowed
	expectedResponseBody := `{"message":"method not allowed"}`

	toTest := NewServer(nil, false, "", "")
	testServer := httptest.NewServer(toTest.h)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/v1/auth/password-reset-request", nil)
	if err != nil {
		t.Fatalf("Failed to build http request: %s", err)
	}

	req.SetBasicAuth("invalid", "invalid")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to call server cause: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedResponseCode {
		t.Errorf("Request respond with unexpected status code. Expected: %d, Given: %d", expectedResponseCode, resp.StatusCode)
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

	if !bytes.Equal(compactedRespBodyAsBytes, []byte(expectedResponseBody)) {
		t.Errorf("Request response body is not as expected. Expected: %q, Given: %q", expectedResponseBody, string(compactedRespBodyAsBytes))
	}
}
