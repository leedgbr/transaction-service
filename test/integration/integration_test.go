package integration_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"transaction-service/test/integration/fixture"
)

var (
	testServer *fixture.TestServer
	client     *fixture.Client
)

func setUp(t *testing.T) {
	testServer = fixture.NewTestServer()
	testServer.Start(t)
	client = fixture.NewClient(testServer.BaseURL)
}

func tearDown() {
	testServer.Stop()
}

func TestStoreTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setUp(t)
		status, body := client.StoreTransaction(t, `{
			"description": "A holiday somewhere nice",
			"transactionDate": "2023-05-01",
			"amountInCents": 100
		}`)

		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{"id":"sequentialID-1"}`, body)
		tearDown()
	})
	t.Run("bad request error", func(t *testing.T) {
		setUp(t)
		status, body := client.StoreTransaction(t, `rubbish`)

		assert.Equal(t, http.StatusBadRequest, status)
		assert.JSONEq(t, `{"message": "BAD_REQUEST"}`, body)
		tearDown()
	})
	t.Run("business validation error", func(t *testing.T) {
		setUp(t)
		status, body := client.StoreTransaction(t, `{
			"transactionDate": "2023-05-01",
			"amountInCents": 100
		}`)

		assert.Equal(t, http.StatusUnprocessableEntity, status)
		assert.JSONEq(t, `{"fields":[{"fieldName": "description", "reason": "REQUIRED"}], "message": "VALIDATION_ERROR"}`, body)
		tearDown()
	})
}

func TestFetchTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setUp(t)
		client.StoreTransaction(t, `{
			"description": "A holiday somewhere nice",
			"transactionDate": "2023-05-01",
			"amountInCents": 100
		}`)
		txnID := "sequentialID-1"
		country := "United%20Kingdom"
		status, body := client.FetchTransaction(t, txnID, country)

		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{
			"transaction": {
				"id": "sequentialID-1", 
				"description": "A holiday somewhere nice",
				"transactionDate": "2023-05-01",
				"amount": {
					"convertedAmountInCents": 35,
					"exchangeRate": 0.345,
					"usdAmountInCents": 100
				}
			}
		}`, body)
		tearDown()
	})
	t.Run("business error", func(t *testing.T) {
		t.Run("validation error", func(t *testing.T) {
			setUp(t)
			client.StoreTransaction(t, `{
			"description": "A holiday somewhere nice",
			"transactionDate": "2023-05-01",
			"amountInCents": 100
		}`)
			txnID := "sequentialID-1"
			country := "a"
			status, body := client.FetchTransaction(t, txnID, country)

			assert.Equal(t, http.StatusUnprocessableEntity, status)
			assert.JSONEq(t, `{"fields":[{"fieldName": "country", "reason": "MIN_LENGTH"}], "message": "VALIDATION_ERROR"}`, body)
			tearDown()
		})
		t.Run("foreign exchange error", func(t *testing.T) {
			setUp(t)
			client.StoreTransaction(t, `{
			"description": "A holiday with no exchange rate",
			"transactionDate": "2022-01-02",
			"amountInCents": 100
		}`)
			txnID := "sequentialID-1"
			country := "United%20Kingdom"
			status, body := client.FetchTransaction(t, txnID, country)
			assert.Equal(t, http.StatusUnprocessableEntity, status)
			assert.JSONEq(t, `{"message": "UNABLE_TO_CONVERT_TO_TARGET_CURRENCY"}`, body)
			tearDown()
		})
	})
}
