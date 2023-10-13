package transaction

import "time"

// Entity represents a transaction entity.
type Entity struct {
	ID              string    `json:"id"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transactionDate"`
	AmountInCents   int       `json:"amountInCents"`
}
