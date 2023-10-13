package fixture

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

// BaseURL returns the base url of the local test server on the supplied port number.
func BaseURL(port int) string {
	return fmt.Sprintf("http://localhost:%d", port)
}

// Post performs a http post operation with the supplied url and body, returning that response status and body.
func Post(t *testing.T, url string, body io.Reader) (int, string) {
	response, err := http.Post(url, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}
	return response.StatusCode, string(responseBody)
}

// Get performs a http get operation with the supplied url, returning that response status and body.
func Get(t *testing.T, url string) (int, string) {
	response, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return response.StatusCode, string(responseBody)
}
