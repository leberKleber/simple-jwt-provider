// +build component

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestRefreshPassword(t *testing.T) {
	email := "refresh_test@leberkleber.io"
	password := "s3cr3t"

	createUser(t, email, password)
	_, refreshToken, authorized := loginUser(t, email, password)
	if !authorized {
		t.Fatal("failed to auth user")
	}

	accessToken, newRefreshToken := refresh(t, refreshToken)

	validateJWT(t, accessToken)
	validateJWT(t, newRefreshToken)

	accessToken, newRefreshToken = refresh(t, newRefreshToken)

	validateJWT(t, accessToken)
	validateJWT(t, newRefreshToken)
}

func refresh(t *testing.T, refreshToken string) (string, string) {
	t.Helper()
	resp, err := http.Post(
		"http://simple-jwt-provider/v1/auth/refresh",
		"application/json",
		bytes.NewReader([]byte(fmt.Sprintf(`{"refresh_token": %q}`, refreshToken))),
	)
	if err != nil {
		t.Fatalf("Failed to refresh with response: %v cause: %s", resp, err)
	}

	responseBody := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ErrorMessage string `json:"message"`
	}{}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return "", ""
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Invalid response status code. Expected: %d, Given: %d, Body: %s", http.StatusOK, resp.StatusCode, responseBody.ErrorMessage)
	}

	return responseBody.AccessToken, responseBody.RefreshToken

}
