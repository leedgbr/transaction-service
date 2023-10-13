package transaction_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"transaction-service/internal/date"
	"transaction-service/internal/transaction"
)

var (
	router *gin.Engine
	rr     *httptest.ResponseRecorder
)

func setUpHandlerTest() {
	router = gin.Default()
	rr = httptest.NewRecorder()
}

func TestStoreHandler(t *testing.T) {
	setUpHandlerTest()
	mockStorer := &MockStorer{}
	transaction.ConfigureStoreHandler(router, mockStorer)

	mockStorer.On("Store", transaction.StoreRequest{
		Description:     stringPtr("*description*"),
		TransactionDate: stringPtr("2023-05-01"),
		AmountInCents:   intPtr(100),
	}).Return(transaction.StoreResponse{
		ID: "*txn-id*",
	}, nil)

	req := newPostRequest(t, "/transaction", `{
			"description": "*description*",
			"transactionDate": "2023-05-01",
			"amountInCents": 100
		}`)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"id": "*txn-id*"}`, rr.Body.String())
	mockStorer.AssertExpectations(t)
}

func TestFetchHandler(t *testing.T) {
	setUpHandlerTest()
	mockFetcher := &MockFetcher{}
	transaction.ConfigureFetchHandler(router, mockFetcher)

	mockFetcher.On("Fetch", mock.Anything, "*txn-id*", "*country*").
		Return(transaction.FetchResponse{
			Transaction: transaction.Response{
				ID:              "*txn-id*",
				Description:     "*description*",
				TransactionDate: &transaction.FormattedDate{Time: date.NewInUTC(2020, time.February, 15)},
				Amount: transaction.Amount{
					USDAmountInCents:       20,
					ConvertedAmountInCents: 30,
					ExchangeRate:           123.45,
				},
			},
		}, nil)

	req := newGetRequest(t, "/transaction/*txn-id*?country=*country*")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{
		"transaction": {
			"id": "*txn-id*", 
			"description": "*description*", 
			"transactionDate": "2020-02-15",
			"amount": {
				"usdAmountInCents": 20, 
				"convertedAmountInCents": 30, 
				"exchangeRate": 123.45
			}
		}
	}`, rr.Body.String())
	mockFetcher.AssertExpectations(t)
}

func newPostRequest(t *testing.T, url, body string) *http.Request {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	return req
}

func newGetRequest(t *testing.T, url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

type MockStorer struct {
	mock.Mock
}

func (m *MockStorer) Store(txn transaction.StoreRequest) (transaction.StoreResponse, error) {
	args := m.Called(txn)
	return args.Get(0).(transaction.StoreResponse), args.Error(1)
}

type MockFetcher struct {
	mock.Mock
}

func (m *MockFetcher) Fetch(ctx context.Context, transactionID, country string) (transaction.FetchResponse, error) {
	args := m.Called(ctx, transactionID, country)
	return args.Get(0).(transaction.FetchResponse), args.Error(1)
}
