package forex

import (
	"io"
	"net/http"
	"strings"
)

// MockHttpClient enables stubbing of http.Client's interface.  This allows us to test details of a call that would
// be made to an api via http.  It stores the request made, so that it can be asserted on later.  It also returns a
// canned http.Response.
type MockHttpClient struct {
	Request        *http.Request
	cannedResponse *http.Response
}

// SetCannedResponse sets the response to return when a request is received.
func (c *MockHttpClient) SetCannedResponse(status int, body string) {
	c.cannedResponse = newResponse(status, body)
}

// Do implements the http.Client's Do operation.  This stub implementation simply stores the request and returns the
// canned response.
func (c *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	c.Request = req
	return c.cannedResponse, nil
}

// newResponse is a convenience function for setting up a new http.Response.
func newResponse(status int, body string) *http.Response {
	response := &http.Response{
		Body:       io.NopCloser(strings.NewReader(body)),
		StatusCode: status,
	}
	return response
}
