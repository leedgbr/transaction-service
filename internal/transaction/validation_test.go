package transaction

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"transaction-service/internal/business"
)

func TestStoreValidation(t *testing.T) {
	validator := storeValidator{}

	t.Run("valid", func(t *testing.T) {
		err := validator.validate(StoreRequest{
			Description:     stringPtr("*description"),
			TransactionDate: stringPtr("2023-01-25"),
			AmountInCents:   intPtr(1),
		})
		assert.Nil(t, err)
	})

	t.Run("mandatory", func(t *testing.T) {
		tcs := []struct {
			name    string
			request StoreRequest
			wantErr []business.FieldError
		}{
			{
				name: "should return an error when one mandatory field is not provided",
				request: StoreRequest{
					Description:     nil,
					TransactionDate: stringPtr("2023-05-20"),
					AmountInCents:   intPtr(123),
				},
				wantErr: []business.FieldError{
					{
						FieldName: "description",
						Reason:    "REQUIRED",
					},
				},
			},
			{
				name:    "should return errors when description, transactionDate and amountInCents are not provided",
				request: StoreRequest{},
				wantErr: []business.FieldError{
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
			},
		}
		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				err := validator.validate(tc.request)
				validationError, ok := err.(*business.Error)
				assert.True(t, ok)
				assert.Equal(t, "VALIDATION_ERROR", validationError.Message)
				assert.ElementsMatch(t, tc.wantErr, validationError.Fields)
			})
		}
	})

	t.Run("correctness", func(t *testing.T) {
		t.Run("amount", func(t *testing.T) {
			t.Run("valid", func(t *testing.T) {
				tcs := []struct {
					name    string
					amount  *int
					wantErr []business.FieldError
				}{
					{
						name:   "should not return an error when amount is one",
						amount: intPtr(1),
					},
					{
						name:   "should not return an error when amount is less than one",
						amount: intPtr(-1),
					},
					{
						name:   "should not return an error when amount large",
						amount: intPtr(1000),
					},
					{
						name:   "should not return an error when amount large negative",
						amount: intPtr(-1000),
					},
				}
				for _, tc := range tcs {
					t.Run(tc.name, func(t *testing.T) {
						request := validRequest()
						request.AmountInCents = tc.amount

						err := validator.validate(request)
						assert.Nil(t, err)
					})
				}
			})
			t.Run("invalid - should return an error when amount is zero", func(t *testing.T) {
				request := validRequest()
				request.AmountInCents = intPtr(0)

				err := validator.validate(request)
				validationError, ok := err.(*business.Error)
				assert.True(t, ok)
				assert.Equal(t, "VALIDATION_ERROR", validationError.Message)
				expectedErr := business.FieldError{
					FieldName: "amountInCents",
					Reason:    "ZERO_VALUE",
				}
				assert.Contains(t, validationError.Fields, expectedErr)
			})
		})
	})
	t.Run("description min / max length", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			tcs := []struct {
				name        string
				description *string
				wantErr     []business.FieldError
			}{
				{
					name:        "should not return an error when description is at minimum length",
					description: stringPtr("1"),
				},
				{
					name:        "should not return an error when description is at maximum length",
					description: stringPtr("12345678901234567890123456789012345678901234567890"),
				},
			}
			for _, tc := range tcs {
				t.Run(tc.name, func(t *testing.T) {
					request := validRequest()
					request.Description = tc.description

					err := validator.validate(request)
					assert.Nil(t, err)
				})
			}
		})
		t.Run("invalid", func(t *testing.T) {
			tcs := []struct {
				name        string
				description *string
				wantErr     []business.FieldError
				wantReason  business.Reason
			}{
				{
					name:        "should return an error when description is one below minimum length",
					description: stringPtr(""),
					wantReason:  "MIN_LENGTH",
				},
				{
					name:        "should return an error when description is one above maximum length",
					description: stringPtr("123456789012345678901234567890123456789012345678901"),
					wantReason:  "MAX_LENGTH",
				},
			}
			for _, tc := range tcs {
				t.Run(tc.name, func(t *testing.T) {
					request := validRequest()
					request.Description = tc.description
					err := validator.validate(request)
					validationError, ok := err.(*business.Error)
					assert.True(t, ok)
					assert.Equal(t, "VALIDATION_ERROR", validationError.Message)
					expectedErr := business.FieldError{
						FieldName: "description",
						Reason:    tc.wantReason,
					}
					assert.Contains(t, validationError.Fields, expectedErr)
				})
			}
		})
	})
	t.Run("transaction date", func(t *testing.T) {
		t.Run("valid - should not return an error when transaction date is today's date correctly formatted", func(t *testing.T) {
			today := time.Now().Format("2006-01-02")
			request := validRequest()
			request.TransactionDate = &today

			err := validator.validate(request)
			assert.Nil(t, err)
		})
		t.Run("invalid", func(t *testing.T) {
			tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
			tcs := []struct {
				name            string
				transactionDate *string
				wantErr         []business.FieldError
				wantReason      business.Reason
			}{
				{
					name:            "should return an error when the transaction date is tomorrow's date correctly formatted",
					transactionDate: &tomorrow,
					wantReason:      "DATE_IN_FUTURE",
				},
				{
					name:            "should return an error when the transaction date is not correctly formatted",
					transactionDate: stringPtr("abcd"),
					wantReason:      "DATE_BAD_FORMAT",
				},
			}
			for _, tc := range tcs {
				t.Run(tc.name, func(t *testing.T) {
					request := validRequest()
					request.TransactionDate = tc.transactionDate

					err := validator.validate(request)
					validationError, ok := err.(*business.Error)
					assert.True(t, ok)
					assert.Equal(t, "VALIDATION_ERROR", validationError.Message)
					expectedErr := business.FieldError{
						FieldName: "transactionDate",
						Reason:    tc.wantReason,
					}
					assert.Contains(t, validationError.Fields, expectedErr)
				})
			}
		})
	})
}

func TestFetchValidation(t *testing.T) {
	validator := fetchValidator{}

	t.Run("valid", func(t *testing.T) {
		tcs := []struct {
			name    string
			country string
		}{
			{
				name:    "country is at minimum length",
				country: "ab",
			},
			{
				name:    "country is one above minimum length",
				country: "abc",
			},
			{
				name: "country is at maximum length",
				country: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890" +
					"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890" +
					"123456789012345678901234567890123456789012345678901234567890123456789012345",
			},
		}
		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				err := validator.validate("*country*")
				assert.Nil(t, err)
			})
		}
	})

	t.Run("invalid", func(t *testing.T) {
		tcs := []struct {
			name    string
			country string
			reason  business.Reason
		}{
			{
				name:    "country is one below minimum length",
				country: "a",
				reason:  "MIN_LENGTH",
			},
			{
				name: "country is one above maximum length",
				country: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890" +
					"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890" +
					"1234567890123456789012345678901234567890123456789012345678901234567890123456",
				reason: "MIN_LENGTH",
			},
		}
		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				err := validator.validate("a")
				wantErr := &business.Error{
					Message: "VALIDATION_ERROR",
					Fields: []business.FieldError{
						{
							FieldName: "country",
							Reason:    "MIN_LENGTH",
						},
					},
				}
				assert.Equal(t, wantErr, err)
			})
		}
	})
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func validRequest() StoreRequest {
	return StoreRequest{
		Description:     stringPtr("*description"),
		TransactionDate: stringPtr("2023-01-25"),
		AmountInCents:   intPtr(1),
	}
}
