package transaction

// StoreRequest represents the user's request to store a transaction
type StoreRequest struct {
	Description     *string `json:"description"`
	TransactionDate *string `json:"transactionDate"`
	AmountInCents   *int    `json:"amountInCents"`
}
