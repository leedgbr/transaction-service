package transaction

import (
	"encoding/json"
	"time"
)

// FetchResponse represents the response for a 'fetch transaction' operation, containing details of the transaction.
type FetchResponse struct {
	Transaction Response `json:"transaction"`
}

// Response represents the fetched transaction details.
type Response struct {
	// ID is the generated id for the transaction.
	ID string `json:"id"`

	// Description is the supplied text description of the transaction.
	Description string `json:"description"`

	// TransactionDate is the date on which the transaction occurred.
	TransactionDate *FormattedDate `json:"transactionDate"`

	// Amount contains details concerning the transaction amount.
	Amount Amount `json:"amount"`
}

// Amount contains the various details relating to the amount of the transaction
type Amount struct {
	// USDAmountInCents is the original transaction amount
	USDAmountInCents int `json:"usdAmountInCents"`

	// ConvertedAmountInCents is the original transaction amount converted to the currency of the requested country
	ConvertedAmountInCents int `json:"convertedAmountInCents"`

	// ExchangeRate is the rate used to convert the USDAmountInCents to ConvertedAmountInCents
	ExchangeRate float64 `json:"exchangeRate"`
}

// FormattedDate enables custom serialization of the transactionDate field to the response.
type FormattedDate struct {
	time.Time
}

// MarshalJSON provides a custom json serialization implementation for the transactionDate.
func (d *FormattedDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}

// StoreResponse represents the response for a 'store transaction' operation, and contains the generated id for
// transaction.
type StoreResponse struct {

	// ID is the transaction's generated id
	ID string `json:"id"`
}
