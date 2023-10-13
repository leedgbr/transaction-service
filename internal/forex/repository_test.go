package forex_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"transaction-service/internal/date"
	"transaction-service/internal/forex"
)

var (
	httpClient *forex.MockHttpClient
	repository forex.Repository
)

func setUpRepository() {
	httpClient = &forex.MockHttpClient{}
	repository = forex.NewTreasuryRepository(httpClient)
}

func TestRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("should send a request that sorts by record date, filters out records older than the "+
			"provided date and filters out records for countries other than the provided country (url encoded)", func(t *testing.T) {
			setUpRepository()
			httpClient.SetCannedResponse(http.StatusOK, `{"data": [{"record_date": "2020-08-01", "exchange_rate": "0.345"}]}`)
			dateOfOldestExchangeRate := date.NewInUTC(2023, time.April, 12)

			repository.FindByCountry(context.Background(), "United Kingdom", dateOfOldestExchangeRate)
			expectedURL := "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?sort=-record_date&format=json&filter=record_date:gte:2023-04-12,country:eq:United+Kingdom&page[size]=1&page[number]=1"
			assert.Equal(t, expectedURL, httpClient.Request.URL.String())
		})
		t.Run("should return record date and exchange rate when exchange record IS found", func(t *testing.T) {
			setUpRepository()
			httpClient.SetCannedResponse(http.StatusOK, `{"data": [{"record_date": "2020-08-01", "exchange_rate": "0.345"}]}`)
			dateOfOldestExchangeRate := date.NewInUTC(2023, time.April, 12)

			result, err := repository.FindByCountry(context.Background(), "United Kingdom", dateOfOldestExchangeRate)
			assert.Nil(t, err)
			expected := forex.Record{
				RecordDate:   forex.RecordDate{Time: date.NewInUTC(2020, time.August, 01)},
				ExchangeRate: forex.ExchangeRate{Value: 0.345},
			}
			assert.Equal(t, expected, result)
		})
		t.Run("should return empty record when exchange record IS NOT found", func(t *testing.T) {
			setUpRepository()
			httpClient.SetCannedResponse(http.StatusOK, `{"data": []}`)
			dateOfOldestExchangeRate := date.NewInUTC(2023, time.April, 12)

			result, err := repository.FindByCountry(context.Background(), "United Kingdom", dateOfOldestExchangeRate)
			assert.Nil(t, err)
			expected := forex.Record{}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("failure", func(t *testing.T) {
		setUpRepository()
		httpClient.SetCannedResponse(http.StatusInternalServerError, `*error-payload*`)
		dateOfOldestExchangeRate := time.Now().AddDate(0, -6, 0)

		result, err := repository.FindByCountry(context.Background(), "United Kingdom", dateOfOldestExchangeRate)
		assert.Equal(t, errors.New("http status 500 received from treasury api. response body: *error-payload*"), err)
		assert.Equal(t, forex.Record{}, result)
	})
}
