package transaction

import (
	"transaction-service/internal/business"
	"transaction-service/internal/validation"
)

const (
	message = "VALIDATION_ERROR"

	countryFieldName = "country"
	countryMinLength = 2

	descriptionFieldName = "description"
	descriptionMinLength = 1
	descriptionMaxLength = 50

	transactionDateFieldName = "transactionDate"

	amountInCentsFieldName = "amountInCents"
)

// storeValidator is responsible for validating input of the 'store transaction' operation.
type storeValidator struct{}

// validate performs business validation on the supplied StoreRequest.
func (v storeValidator) validate(transaction StoreRequest) error {
	if err := mandatory(transaction); err != nil {
		return err
	}
	if err := correctness(transaction); err != nil {
		return err
	}
	return nil
}

// mandatory performs business validation relating to mandatory fields on the StoreRequest.
func mandatory(transaction StoreRequest) error {
	var fieldErrors []business.FieldError
	if err := validation.IsRequiredString(descriptionFieldName, transaction.Description); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if err := validation.IsRequiredString(transactionDateFieldName, transaction.TransactionDate); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if err := validation.IsRequiredInt(amountInCentsFieldName, transaction.AmountInCents); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	return checkForErrors(fieldErrors)
}

// correctness performs business validation relating to correctness of field values on the StoreRequest.  It would be
// prudent to apply a sensible minimum and maximum amount constraint to help safeguard against overflow during
// conversion, although I have not done this here.
func correctness(transaction StoreRequest) error {
	var fieldErrors []business.FieldError
	if err := validation.IsMinLength(descriptionFieldName, transaction.Description, descriptionMinLength); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if err := validation.IsMaxLength(descriptionFieldName, transaction.Description, descriptionMaxLength); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	transactionDate, err := validation.IsDate(transactionDateFieldName, transaction.TransactionDate)
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if err := validation.IsDateNowOrEarlier(transactionDateFieldName, transactionDate); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if err := validation.IsNotZero(amountInCentsFieldName, transaction.AmountInCents); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	return checkForErrors(fieldErrors)
}

// checkForErrors returns a business error containing the fieldErrors if any fieldErrors are provided.
func checkForErrors(fieldErrors []business.FieldError) error {
	if len(fieldErrors) > 0 {
		return &business.Error{
			Message: message,
			Fields:  fieldErrors,
		}
	}
	return nil
}

// fetchValidator is responsible for validating input of the 'fetch transaction' operation.
type fetchValidator struct{}

// validate performs business validation on the supplied country field value.
func (v *fetchValidator) validate(country string) error {
	var fieldErrors []business.FieldError
	if err := validation.IsRequiredString(countryFieldName, &country); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	if err := validation.IsMinLength(countryFieldName, &country, countryMinLength); err != nil {
		fieldErrors = append(fieldErrors, *err)
	}
	return checkForErrors(fieldErrors)
}
