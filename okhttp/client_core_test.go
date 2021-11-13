package okhttp

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestGetRequestHeaders(t *testing.T) {
	commonHeaders := make(http.Header)
	commonHeaders.Set("Content-Type", "application/json")
	commonHeaders.Set("User-Agent", "ok-http-client1")
	builder := clientBuilder{headers: commonHeaders}
	client := httpClient{builder: &builder}

	requestHeaders := make(http.Header)
	// Should be added
	requestHeaders.Set("X-request-id", "ABC-123")
	// Should be overwritten
	requestHeaders.Set("User-Agent", "ok-http-client2")
	finalHeaders := client.getRequestHeaders(requestHeaders)

	if len(finalHeaders) != 3 {
		t.Error("Expect 3 headers, got " + strconv.Itoa(len(finalHeaders)))
	}

	if finalHeaders.Get("Content-Type") != "application/json" {
		t.Error("Expect Content-Tupe header to be set as a common header")
	}

	if finalHeaders.Get("User-Agent") != "ok-http-client2" {
		t.Error("Expect User-Agent header to be overwritten by custom header")
	}

	if finalHeaders.Get("X-request-id") != "ABC-123" {
		t.Error("Expect X-request-id header to be set as a custom header")
	}
}

// One test for each return
func TestGetRequestBody(t *testing.T) {

	t.Run("nilBody", func(t *testing.T) {
		client := httpClient{}

		requestBody, err := client.getRequestBody("", nil)

		if err != nil {
			t.Error("Expected error to be nil")
		}

		if requestBody != nil {
			t.Error("Expected request body to be nil")
		}
	})

	t.Run("jsonBody", func(t *testing.T) {
		client := httpClient{}

		body := []string{"one", "two"}
		requestBody, err := client.getRequestBody("application/json", body)

		if err != nil {
			t.Error("Expected error to be nil")
		}

		if string(requestBody) != `["one","two"]` {
			t.Error("Expected request body be a stringified json array")
		}
	})

	t.Run("xmlBody", func(t *testing.T) {
		client := httpClient{}

		body := []string{"one", "two"}
		requestBody, err := client.getRequestBody("application/xml", body)

		if err != nil {
			t.Error("Expected error to be nil")
		}

		fmt.Println(string(requestBody))

		if string(requestBody) != `<string>one</string><string>two</string>` {
			t.Error("Expected request body be a stringified json array")
		}
	})

	t.Run("defaultBody", func(t *testing.T) {
		client := httpClient{}

		body := []string{"one", "two"}
		requestBody, err := client.getRequestBody("", body)

		if err != nil {
			t.Error("Expected error to be nil")
		}

		if string(requestBody) != `["one","two"]` {
			t.Error("Expected request body be a stringified json array")
		}
	})
}
