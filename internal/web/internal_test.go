package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAliveHandler(t *testing.T) {
	expectedResponseCode := http.StatusOK
	expectedResponseBody := `{"alive":true}`

	toTest := NewServer(nil, false, "", "")
	testServer := httptest.NewServer(toTest.h)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/v1/internal/alive", nil)
	if err != nil {
		t.Fatalf("Failed to build http request: %s", err)
	}

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
