// +build component

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type User struct {
	EMail    string                 `json:"email,omitempty"`
	Password string                 `json:"password,omitempty"`
	Claims   map[string]interface{} `json:"claims,omitempty"`
}

func createUser(t *testing.T, email, password string) {
	t.Helper()
	req, err := http.NewRequest(
		http.MethodPost,
		"http://simple-jwt-provider/v1/admin/users",
		bytes.NewReader([]byte(fmt.Sprintf(`{"email": %q, "password": %q, "claims": {"myCustomClaim": "customClaimValue"}}`, email, password))),
	)
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

func readUser(t *testing.T, email string) User {
	t.Helper()
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://simple-jwt-provider/v1/admin/users/%s", url.PathEscape(email)),
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create http request")
	}

	req.SetBasicAuth("username", "password")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create user cause: %s", err)
	}

	var responseBody User
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Error("Failed to read response body", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Invalid response status code. Expected: %d, Given: %d, Body: %#v", http.StatusOK, resp.StatusCode, responseBody)
	}

	return responseBody
}

func updateUser(t *testing.T, email string, newPassword string, newClaims map[string]interface{}) {
	t.Helper()

	requestBody := User{
		Password: newPassword,
		Claims:   newClaims,
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(requestBody)
	if err != nil {
		t.Fatal("failed to encode request body", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("http://simple-jwt-provider/v1/admin/users/%s", url.PathEscape(email)),
		&body,
	)
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Invalid response status code. Expected: %d, Given: %d, Body: %s", http.StatusOK, resp.StatusCode, respBody)
	}
}

func deleteUser(t *testing.T, email string) {
	t.Helper()
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("http://simple-jwt-provider/v1/admin/users/%s", url.QueryEscape(email)),
		nil,
	)
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

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Invalid response status code. Expected: %d, Given: %d, Body: %s", http.StatusOK, resp.StatusCode, respBody)
	}
}
