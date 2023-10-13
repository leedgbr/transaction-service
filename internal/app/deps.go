package app

import (
	"transaction-service/internal/forex"
	"transaction-service/internal/transaction"
)

// NewDependencies wires up the application's dependencies using the dependency injection pattern
func NewDependencies(txnIDGenerator transaction.IDGenerator, httpClient forex.HttpClient) Dependencies {
	txnRepository := transaction.NewInMemoryRepository(txnIDGenerator)
	forExRepository := forex.NewTreasuryRepository(httpClient)
	forExService := forex.NewRepositoryService(forExRepository)
	txnService := transaction.NewRepositoryService(txnRepository, forExService)
	return Dependencies{
		TxnService: txnService,
	}
}

// Dependencies holds the top level dependencies required for wiring to handlers.
type Dependencies struct {
	TxnService *transaction.RepositoryService
}
