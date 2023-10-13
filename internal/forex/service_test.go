package forex_test

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
)

var (
	ctx      context.Context
	mockRepo MockRepository
	service  *forex.RepositoryService
)

func TestService(t *testing.T) {
	tcs := []struct {
		name       string
		record     forex.Record
		err        error
		wantErr    error
		wantResult forex.ConversionResult
	}{
		{
			name: "should return currency conversion details when an exchange rate record is found",
			record: forex.Record{
				RecordDate: forex.RecordDate{
					Time: date.NewInUTC(2023, time.April, 4),
				},
				ExchangeRate: forex.ExchangeRate{
					Value: 0.745,
				},
			},
			err:     nil,
			wantErr: nil,
			wantResult: forex.ConversionResult{
				Amount:       9197,
				ExchangeRate: 0.745,
			},
		},
		{
			name:   "should return an error when no exchange rate record is found",
			record: forex.Record{},
			err:    nil,
			wantErr: &business.Error{
				Message: "UNABLE_TO_CONVERT_TO_TARGET_CURRENCY",
			},
			wantResult: forex.ConversionResult{},
		},
		{
			name:       "should return an error when there is system problem retrieving the exchange rate record",
			record:     forex.Record{},
			err:        errors.New("problem"),
			wantErr:    errors.New("problem"),
			wantResult: forex.ConversionResult{},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			setUpService()
			dateOfOldestRecord := date.NewInUTC(2023, time.February, 10)
			amountInCents := 12345
			mockRepo.On("FindByCountry", ctx, "*country*", dateOfOldestRecord).
				Return(tc.record, tc.err)

			result, err := service.Convert(context.Background(), "*country*", dateOfOldestRecord, amountInCents)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResult, result)
			mockRepo.AssertExpectations(t)
		})
	}
}

func setUpService() {
	ctx = context.Background()
	mockRepo = MockRepository{}
	service = forex.NewRepositoryService(&mockRepo)
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindByCountry(ctx context.Context, country string, oldest time.Time) (forex.Record, error) {
	args := m.Called(ctx, country, oldest)
	return args.Get(0).(forex.Record), args.Error(1)
}
