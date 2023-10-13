package forex

import (
	"context"
	"time"

	"transaction-service/internal/business"
)

const (
	unableToConvertToTargetCurrency = "UNABLE_TO_CONVERT_TO_TARGET_CURRENCY"
)

// ConversionResult represents the output of a currency conversion operation
type ConversionResult struct {
	Amount       int
	ExchangeRate float64
}

// NewRepositoryService creates a RepositoryService that uses the supplied repository and default Converter for
// performing exchange rate calculations.
func NewRepositoryService(repository Repository) *RepositoryService {
	return &RepositoryService{
		repository: repository,
		converter:  Converter{},
	}
}

// Repository defines the interface expected of the Repository for finding exchange rate records.
type Repository interface {
	FindByCountry(ctx context.Context, country string, oldest time.Time) (Record, error)
}

// RepositoryService is the business service for performing foreign exchange currency conversion calculations.
type RepositoryService struct {
	repository Repository
	converter  Converter
}

// Convert will convert the provided amount (in cents) to the currency of the specified country, using an exchange
// rate sourced from the configured data source which is not older than the provided dateOfOldestExchangeRate.  If no
// suitable exchange rate can be found, an error will be returned.
func (s *RepositoryService) Convert(ctx context.Context,
	country string,
	dateOfOldestExchangeRate time.Time,
	amountInCents int) (ConversionResult, error) {

	record, err := s.repository.FindByCountry(ctx, country, dateOfOldestExchangeRate)
	if err != nil {
		return ConversionResult{}, err
	}
	if record == (Record{}) {
		return ConversionResult{}, &business.Error{Message: unableToConvertToTargetCurrency}
	}
	exchangeRate := record.ExchangeRate.Value
	return ConversionResult{
		Amount:       s.converter.Convert(amountInCents, exchangeRate),
		ExchangeRate: exchangeRate,
	}, nil
}
