package forex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL    = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?sort=-record_date&format=json"
	pageSize   = 1
	pageNumber = 1
	dateFormat = "2006-01-02"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewTreasuryRepository creates a new TreasuryRepository with the supplied httpClient.
func NewTreasuryRepository(httpClient HttpClient) *TreasuryRepository {
	return &TreasuryRepository{
		httpClient: httpClient,
	}
}

// TreasuryRepository is responsible for making calls to the US Treasury Reporting of Exchange API
// (https://fiscaldata.treasury.gov/datasets/treasury-reporting-rates-exchange/treasury-reporting-rates-of-exchange)
type TreasuryRepository struct {
	httpClient HttpClient
}

// FindByCountry returns the most recent foreign exchange record for the specified country that is not older than the
// specified dateOfOldestRecord.
func (r *TreasuryRepository) FindByCountry(ctx context.Context, country string, dateOfOldestRecord time.Time) (Record, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, newURL(country, dateOfOldestRecord), nil)
	if err != nil {
		return Record{}, err
	}
	response, err := r.httpClient.Do(request)
	if err != nil {
		return Record{}, err
	}
	if response.StatusCode != http.StatusOK {
		body, err := read(response.Body)
		if err != nil {
			return Record{}, err
		}
		return Record{}, fmt.Errorf("http status %d received from treasury api. response body: %s", response.StatusCode, string(body))
	}
	unmarshalled, err := parseResponse(response, err)
	if err != nil {
		return Record{}, err
	}
	if len(unmarshalled.Data) <= 0 {
		return Record{}, nil
	}
	return unmarshalled.Data[0], nil
}

// parseResponse reads the http.Response into an APIResponse or returns an error.
func parseResponse(response *http.Response, err error) (APIResponse, error) {
	body, err := read(response.Body)
	if err != nil {
		return APIResponse{}, err
	}
	var unmarshalled APIResponse
	err = json.Unmarshal(body, &unmarshalled)
	if err != nil {
		return APIResponse{}, err
	}
	return unmarshalled, nil
}

func read(closer io.ReadCloser) ([]byte, error) {
	defer closer.Close()
	return io.ReadAll(closer)
}

// newURL creates the url to use to call the Treasury API, using the suppliec country and date of oldest record
func newURL(country string, dateOfOldestRecord time.Time) string {
	return fmt.Sprintf("%s&filter=record_date:gte:%s,country:eq:%s&page[size]=%d&page[number]=%d",
		baseURL, dateOfOldestRecord.Format(dateFormat), url.QueryEscape(country), pageSize, pageNumber)
}
