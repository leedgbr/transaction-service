package transaction

import (
	"context"
	"time"

	"transaction-service/internal/business"
	"transaction-service/internal/forex"
	"transaction-service/internal/validation"
)

const (
	transactionNotFound = "TRANSACTION_NOT_FOUND"
)

// ForExService is the expected interface for the service used to determine the exchange rate and perform the
// exchange rate calculation
type ForExService interface {
	Convert(ctx context.Context, country string, dateOfOldestExchangeRate time.Time, amountInCents int) (forex.ConversionResult, error)
}

// Repository is the expected interface for the repository of transactions.
type Repository interface {
	Save(transaction Entity) (Entity, error)
	FindByID(id string) Entity
}

// NewRepositoryService creates a RepositoryService that uses the supplied transaction repository and foreign exchange
// service.
func NewRepositoryService(txnRepository Repository, forExService ForExService) *RepositoryService {
	return &RepositoryService{
		txnRepository:  txnRepository,
		forExService:   forExService,
		fetchValidator: fetchValidator{},
		storeValidator: storeValidator{},
	}
}

// RepositoryService is responsible for orchestrating the processes to store transactions and fetch transactions with
// the amount converted to the currency of the requested country.
type RepositoryService struct {
	txnRepository  Repository
	forExService   ForExService
	storeValidator storeValidator
	fetchValidator fetchValidator
}

// Store first ensures the request is validated, then stores the transaction in the repository and returns the new id
// generated for the transaction.
func (s *RepositoryService) Store(txn StoreRequest) (StoreResponse, error) {
	if err := s.storeValidator.validate(txn); err != nil {
		return StoreResponse{}, err
	}
	entity, err := mapToEntity(txn)
	if err != nil {
		return StoreResponse{}, err
	}
	updated, err := s.txnRepository.Save(entity)
	if err != nil {
		return StoreResponse{}, err
	}
	return StoreResponse{
		ID: updated.ID,
	}, nil
}

// Fetch first ensures the country is validated, then fetches the transaction from the repository, has its amount
// converted to the currency of the requested country and returns the transaction details, including the exchange rate
// used and the converted currency amount.
//
// transactionID cannot be invalid since the path parameter used in the route makes this impossible.  We could add
// validation for transactionID here, but since it will never be executed in the current configuration I have left it
// out for now.
func (s *RepositoryService) Fetch(ctx context.Context, transactionID, country string) (FetchResponse, error) {
	if err := s.fetchValidator.validate(country); err != nil {
		return FetchResponse{}, err
	}
	entity := s.txnRepository.FindByID(transactionID)
	if entity == (Entity{}) {
		return FetchResponse{}, &business.Error{Message: transactionNotFound}
	}
	dateOfOldestExchangeRate := monthsOlderThan(entity.TransactionDate, 6)
	result, err := s.forExService.Convert(ctx, country, dateOfOldestExchangeRate, entity.AmountInCents)
	if err != nil {
		return FetchResponse{}, err
	}
	return FetchResponse{
		Transaction: Response{
			ID:          entity.ID,
			Description: entity.Description,
			TransactionDate: &FormattedDate{
				Time: entity.TransactionDate,
			},
			Amount: Amount{
				USDAmountInCents:       entity.AmountInCents,
				ConvertedAmountInCents: result.Amount,
				ExchangeRate:           result.ExchangeRate,
			},
		},
	}, nil
}

// monthsOlderThan returns a time.Time representing a date that is numberOfMonths earlier than the date provided.
func monthsOlderThan(date time.Time, numberOfMonths int) time.Time {
	return date.AddDate(0, numberOfMonths*-1, 0)
}

// mapToEntity maps the provided transaction StoreRequest into a transaction Entity.
func mapToEntity(txn StoreRequest) (Entity, error) {
	txnDate, err := mapDate(txn.TransactionDate)
	if err != nil {
		return Entity{}, err
	}
	entity := Entity{
		Description:     *txn.Description,
		TransactionDate: txnDate,
		AmountInCents:   *txn.AmountInCents,
	}
	return entity, nil
}

// mapDate parses the supplied date into a time.Time according to the configured date format.
func mapDate(date *string) (time.Time, error) {
	if date == nil {
		return time.Time{}, nil
	}
	mappedDate, err := time.Parse(validation.DateFormat, *date)
	if err != nil {
		return time.Time{}, err
	}
	return mappedDate, nil
}
