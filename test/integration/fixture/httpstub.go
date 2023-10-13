package fixture

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	noExchangeRecordTreasuryURL = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?sort=-record_date&format=json&filter=record_date:gte:2021-07-02,country:eq:United+Kingdom&page[size]=1&page[number]=1"
	noExchangeRateTreasuryBody  = `{"data": []}`

	treasuryURL  = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?sort=-record_date&format=json&filter=record_date:gte:2022-11-01,country:eq:United+Kingdom&page[size]=1&page[number]=1"
	treasuryBody = `{"data": [{"record_date": "2020-08-01", "exchange_rate": "0.345"}]}`
)

// NewStubHttpClient creates a StubHttpClient configured with a stub response.
func NewStubHttpClient() *StubHttpClient {
	return &StubHttpClient{
		requestDetails: map[RequestDetails]*http.Response{
			{
				Method: http.MethodGet,
				URL:    treasuryURL,
			}: newResponse(http.StatusOK, treasuryBody),
			{
				Method: http.MethodGet,
				URL:    noExchangeRecordTreasuryURL,
			}: newResponse(http.StatusOK, noExchangeRateTreasuryBody),
		},
	}
}

// StubHttpClient is a simple stub implementation of a http client.  It is used to avoid a dependency on the real
// Treasury API when running integration tests.  There are various reasons for doing this, including the fact that we
// would like these tests to be repeatable and stable.
type StubHttpClient struct {
	requestDetails map[RequestDetails]*http.Response
}

// Do returns the configured http.Response that matches the provided http.Request.
func (c *StubHttpClient) Do(req *http.Request) (*http.Response, error) {
	details, err := newRequestDetails(req)
	if err != nil {
		return nil, err
	}
	response, ok := c.requestDetails[details]
	if !ok {
		return nil, fmt.Errorf("http client stub missing canned response for: %+v", details)
	}
	return response, nil
}

// RequestDetails represents the details on which incoming requests will be matched.
type RequestDetails struct {
	Method string
	URL    string
}

// newRequestDetails is a convenience function for creating a new RequestDetails with the supplied http.Request.
func newRequestDetails(req *http.Request) (RequestDetails, error) {
	if req == nil {
		return RequestDetails{}, errors.New("no request provided")
	}
	return RequestDetails{
		Method: req.Method,
		URL:    req.URL.String(),
	}, nil

}

// newResponse creates a new http.Response with the supplied http status and body.
func newResponse(status int, body string) *http.Response {
	response := &http.Response{
		Body:       io.NopCloser(strings.NewReader(body)),
		StatusCode: status,
	}
	return response
}
