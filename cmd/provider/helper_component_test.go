// +build component

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func createUser(t *testing.T, email, password string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "http://simple-jwt-provider/v1/admin/users", bytes.NewReader([]byte(fmt.Sprintf(`{"email": %q, "password": %q}`, email, password))))
	if err != nil {
		t.Fatalf("Failed to create http request")
	}

	req.SetBasicAuth("username", "password")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create user cause: %s", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body")
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Invalid response status code. Expected: %d, Given: %d, Body: %s", http.StatusOK, resp.StatusCode, respBody)
	}
}
