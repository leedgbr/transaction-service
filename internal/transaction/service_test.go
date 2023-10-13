package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"transaction-service/internal/business"
	"transaction-service/internal/date"
	"transaction-service/internal/forex"
	"transaction-service/internal/transaction"
)

var (
	ctx       context.Context
	mockForEx MockForEx
	mockRepo  MockRepository
	service   *transaction.RepositoryService
)

func TestServiceStore(t *testing.T) {
	t.Run("success - should return the generated id of the stored transaction", func(t *testing.T) {
		setUp()
		mockRepo.On("Save", transaction.Entity{
			Description:     "*description*",
			TransactionDate: date.NewInUTC(2022, time.October, 1),
			AmountInCents:   345,
		}).Return(transaction.Entity{
			ID:              "*saved*",
			Description:     "*description*",
			TransactionDate: date.NewInUTC(2022, time.October, 1),
			AmountInCents:   345,
		}, nil)

		response, err := service.Store(transaction.StoreRequest{
			Description:     stringPtr("*description*"),
			TransactionDate: stringPtr("2022-10-01"),
			AmountInCents:   intPtr(345),
		})

		assert.Nil(t, err)
		assert.Equal(t, transaction.StoreResponse{ID: "*saved*"}, response)
		mockRepo.AssertExpectations(t)
		mockForEx.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("should return a validation error when the request does not meet the business validation rules", func(t *testing.T) {
			setUp()

			response, err := service.Store(transaction.StoreRequest{})
			expectedErr := &business.Error{
				Fields: []business.FieldError{
					{
						FieldName: "description",
						Reason:    "REQUIRED",
					},
					{
						FieldName: "transactionDate",
						Reason:    "REQUIRED",
					},
					{
						FieldName: "amountInCents",
						Reason:    "REQUIRED",
					},
				},
				Message: "VALIDATION_ERROR",
			}
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, transaction.StoreResponse{}, response)
			mockRepo.AssertExpectations(t)
			mockForEx.AssertExpectations(t)
		})

		t.Run("should return an error when there is a problem with the transaction repository", func(t *testing.T) {
			setUp()
			mockRepo.On("Save", transaction.Entity{
				Description:     "*description*",
				TransactionDate: date.NewInUTC(2022, time.October, 1),
				AmountInCents:   345,
			}).Return(transaction.Entity{}, errors.New("problem"))

			response, err := service.Store(transaction.StoreRequest{
				Description:     stringPtr("*description*"),
				TransactionDate: stringPtr("2022-10-01"),
				AmountInCents:   intPtr(345),
			})

			assert.Equal(t, errors.New("problem"), err)
			assert.Equal(t, transaction.StoreResponse{}, response)
			mockRepo.AssertExpectations(t)
			mockForEx.AssertExpectations(t)
		})
	})
}

func TestServiceFetch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("should return the fetched transaction details with the requested currency conversion for the supplied country", func(t *testing.T) {
			setUp()
			mockRepo.On("FindByID", "*txn-id*").
				Return(transaction.Entity{
					ID:              "*txn-id*",
					Description:     "*description*",
					TransactionDate: date.NewInUTC(2022, time.May, 12),
					AmountInCents:   543,
				}, nil)
			mockForEx.On("Convert", ctx, "*country*", mock.Anything, 543).
				Return(forex.ConversionResult{
					Amount:       1234,
					ExchangeRate: 0.456,
				}, nil)

			response, err := service.Fetch(ctx, "*txn-id*", "*country*")

			assert.Nil(t, err)
			expectedResponse := transaction.FetchResponse{
				Transaction: transaction.Response{
					ID:          "*txn-id*",
					Description: "*description*",
					TransactionDate: &transaction.FormattedDate{
						Time: date.NewInUTC(2022, time.May, 12),
					},
					Amount: transaction.Amount{
						USDAmountInCents:       543,
						ConvertedAmountInCents: 1234,
						ExchangeRate:           0.456,
					},
				},
			}
			assert.Equal(t, expectedResponse, response)
			mockRepo.AssertExpectations(t)
			mockForEx.AssertExpectations(t)
		})
		t.Run("should request a foreign exchange rate that was recorded within six months of the transaction date", func(t *testing.T) {
			setUp()
			mockRepo.On("FindByID", mock.Anything).
				Return(transaction.Entity{
					TransactionDate: date.NewInUTC(2022, time.May, 12),
				}, nil)
			mockForEx.On("Convert", ctx, mock.Anything, date.NewInUTC(2021, time.November, 12), mock.Anything).
				Return(forex.ConversionResult{}, nil)

			service.Fetch(ctx, "*txn-id*", "*country*")
			mockForEx.AssertExpectations(t)
		})
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("should return a validation error when the input does not satisfy the business rules", func(t *testing.T) {
			setUp()
			response, err := service.Fetch(ctx, "*txn-id*", "")

			expectedErr := &business.Error{
				Fields: []business.FieldError{
					{
						FieldName: "country",
						Reason:    "MIN_LENGTH",
					},
				},
				Message: "VALIDATION_ERROR",
			}
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, transaction.FetchResponse{}, response)
			mockRepo.AssertExpectations(t)
			mockForEx.AssertExpectations(t)
		})

		t.Run("should return an error when a transaction with the supplied id cannot be found", func(t *testing.T) {
			setUp()
			mockRepo.On("FindByID", "*txn-id*").
				Return(transaction.Entity{}, nil)

			response, err := service.Fetch(ctx, "*txn-id*", "*country*")

			expectedErr := &business.Error{
				Message: "TRANSACTION_NOT_FOUND",
			}
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, transaction.FetchResponse{}, response)
			mockRepo.AssertExpectations(t)
			mockForEx.AssertExpectations(t)
		})

		t.Run("should return an error when there is a problem with the foreign exchange conversion", func(t *testing.T) {
			setUp()
			mockRepo.On("FindByID", "*txn-id*").
				Return(transaction.Entity{
					ID:              "*txn-id*",
					Description:     "*description*",
					TransactionDate: date.NewInUTC(2022, time.May, 12),
					AmountInCents:   543,
				}, nil)
			mockForEx.On("Convert", ctx, "*country*", mock.Anything, 543).
				Return(forex.ConversionResult{}, errors.New("problem"))

			response, err := service.Fetch(ctx, "*txn-id*", "*country*")

			assert.Equal(t, errors.New("problem"), err)
			assert.Equal(t, transaction.FetchResponse{}, response)
			mockRepo.AssertExpectations(t)
			mockForEx.AssertExpectations(t)
		})
	})
}

func setUp() {
	ctx = context.Background()
	mockForEx = MockForEx{}
	mockRepo = MockRepository{}
	service = transaction.NewRepositoryService(&mockRepo, &mockForEx)
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(txn transaction.Entity) (transaction.Entity, error) {
	args := m.Called(txn)
	return args.Get(0).(transaction.Entity), args.Error(1)
}

func (m *MockRepository) FindByID(id string) transaction.Entity {
	args := m.Called(id)
	return args.Get(0).(transaction.Entity)
}

type MockForEx struct {
	mock.Mock
}

func (m *MockForEx) Convert(ctx context.Context, country string, dateOfOldestExchangeRate time.Time, amountInCents int) (forex.ConversionResult, error) {
	args := m.Called(ctx, country, dateOfOldestExchangeRate, amountInCents)
	return args.Get(0).(forex.ConversionResult), args.Error(1)
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
