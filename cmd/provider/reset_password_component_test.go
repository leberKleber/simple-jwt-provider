// +build component

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestResetPassword(t *testing.T) {
	email := "reset_test@leberkleber.io"
	password := "s3cr3t"
	newPassword := "t3rc3s"

	createUser(t, email, password)
	createPasswordResetRequest(t, email)
	token := findPasswordResetTokenFromMailAndVerifyContent(t, email)
	resetPassword(t, email, token, newPassword)

	loginUser(t, email, newPassword)
}

func findPasswordResetTokenFromMailAndVerifyContent(t *testing.T, email string) string {
	resp, err := http.Get("http://mail-server:8025/api/v2/messages")
	if err != nil {
		t.Fatalf("Failed to login cause: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Invalid response status code. Expected: %d, Given: %d", http.StatusOK, resp.StatusCode)
	}

	var mailhogRes MailhogResponse

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&mailhogRes)
	if err != nil {
		t.Fatalf("failed to encode smtp-server api-response: %s", err)
	}

	var respMail MailhogResponseItemRaw
	respMailFound := false
	for _, r := range mailhogRes.Items {
		for i := range r.Raw.To {
			if r.Raw.To[i] == email {
				respMail = r.Raw
				respMailFound = true
				break
			}
		}
	}

	if !respMailFound {
		t.Fatal("could not find mail body")
	}

	expectedCustomCalimValue := "customClaimValue"
	if !strings.Contains(respMail.Data, "customClaimValue") {
		t.Errorf("email body dosent contains custom claim value %q: \n%q", expectedCustomCalimValue, respMail.Data)
	}

	reg, err := regexp.Compile("([a-f0-9]{64})")
	if err != nil {
		t.Fatal("could not compile regex")
	}

	res := reg.Find([]byte(respMail.Data))
	if len(res) == 0 {
		t.Fatalf("no reset token found. Mail content %q", respMail.Data)
	}

	return string(res)
}

func createPasswordResetRequest(t *testing.T, email string) {
	resp, err := http.Post(
		"http://simple-jwt-provider/v1/auth/password-reset-request",
		"application/json",
		bytes.NewReader([]byte(fmt.Sprintf(`{"email": %q}`, email))),
	)
	if err != nil {
		t.Fatalf("Failed to create password-reset-request cause: %s", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Invalid response status code. Expected: %d, Given: %d", http.StatusCreated, resp.StatusCode)
	}

	return
}

func resetPassword(t *testing.T, email, token, newPassword string) {
	resp, err := http.Post(
		"http://simple-jwt-provider/v1/auth/password-reset",
		"application/json",
		bytes.NewReader([]byte(fmt.Sprintf(`{"email": %q, "reset_token":%q, "password": %q}`, email, token, newPassword))),
	)
	if err != nil {
		t.Fatalf("Failed to login cause: %s", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Invalid response status code. Expected: %d, Given: %d", http.StatusNoContent, resp.StatusCode)
	}
}

type MailhogResponse struct {
	Items []MailhogResponseItem `json:"items"`
}
type MailhogResponseItem struct {
	Raw MailhogResponseItemRaw `json:"Raw"`
}

type MailhogResponseItemRaw struct {
	To   []string `json:"To"`
	Data string   `json:"Data"`
}
