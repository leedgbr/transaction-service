package fixture

import (
	"fmt"
	"strings"
	"testing"
)

// NewClient returns a Client configured with the supplied baseURL.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

// Client is a very basic client for calling the transaction-service api.
type Client struct {
	baseURL string
}

// StoreTransaction calls the 'store transaction' operation with the supplied payload, returning the response status and
// body.  Should an error occur, the current test will be failed.
func (c *Client) StoreTransaction(t *testing.T, payload string) (int, string) {
	url := fmt.Sprintf("%s/transaction", c.baseURL)
	return Post(t, url, strings.NewReader(payload))
}

// FetchTransaction calls the 'fetch transaction' operation with the supplied transaction id and country, returning the
// response status and body.  Should an error occur, the current test will be failed.
func (c *Client) FetchTransaction(t *testing.T, id, country string) (int, string) {
	url := fmt.Sprintf("%s/transaction/%s?country=%s", c.baseURL, id, country)
	return Get(t, url)
}
