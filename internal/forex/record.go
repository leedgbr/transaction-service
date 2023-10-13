package forex

import (
	"strconv"
	"time"
)

// APIResponse represents the response received from the Treasury Exchange Rate API.
type APIResponse struct {
	Data []Record `json:"data"`
}

// RecordDate is used to help parse the record_date field from the APIResponse.
type RecordDate struct {
	time.Time
}

// ExchangeRate is used to help parse the exchange_rate field from the APIResponse.
type ExchangeRate struct {
	Value float64
}

// UnmarshalJSON is a custom json deserialization implementation to read a float from the exchange_rate field.
func (r *ExchangeRate) UnmarshalJSON(bytes []byte) error {
	unquotedValue, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}
	value, err := strconv.ParseFloat(unquotedValue, 64)
	if err != nil {
		return err
	}
	r.Value = value
	return nil
}

// Record represents an exchange rate record received from the Treasury API
type Record struct {
	RecordDate   RecordDate   `json:"record_date"`
	ExchangeRate ExchangeRate `json:"exchange_rate"`
}

// UnmarshalJSON is a custom json deserialization implementation to read a date from the record_date field.
func (d *RecordDate) UnmarshalJSON(bytes []byte) error {
	unquotedValue, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}
	date, err := time.Parse(dateFormat, unquotedValue)
	if err != nil {
		return err
	}
	d.Time = date
	return nil
}
