// +build component

package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAliveResource(t *testing.T) {
	res, err := http.Get("http://simple-auth-provider/v1/internal/alive")
	if err != nil {
		t.Fatalf("Failed to call simple-auth-provider: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Response status code is not as expected. Expected %d, Given %d", http.StatusOK, res.StatusCode)
	}
	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read request body: %s", err)
	}
	expectedResponseBody := `{"alive":true}`
	if string(reqBody) != expectedResponseBody {
		t.Errorf("Response body is not as expected. Expected %q, Given %q", expectedResponseBody, reqBody)
	}
}
