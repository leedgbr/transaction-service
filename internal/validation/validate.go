package validation

import (
	"time"

	"transaction-service/internal/business"
)

const (
	Required      business.Reason = "REQUIRED"
	MinLength     business.Reason = "MIN_LENGTH"
	MaxLength     business.Reason = "MAX_LENGTH"
	DateBadFormat business.Reason = "DATE_BAD_FORMAT"
	DateInFuture  business.Reason = "DATE_IN_FUTURE"
	ZeroValue     business.Reason = "ZERO_VALUE"

	DateFormat = "2006-01-02"
)

// IsRequiredString returns a business.FieldError if the supplied string field value is empty or nil.
func IsRequiredString(fieldName string, value *string) *business.FieldError {
	if value == nil {
		return business.NewFieldError(fieldName, Required)
	}
	return nil
}

// IsRequiredInt returns a business.FieldError if the supplied int value is nil.
func IsRequiredInt(fieldName string, value *int) *business.FieldError {
	if value == nil {
		return business.NewFieldError(fieldName, Required)
	}
	return nil
}

// IsMinLength returns a business.FieldError if the supplied string value has a length which is less than the supplied
// minimum
func IsMinLength(fieldName string, value *string, min int) *business.FieldError {
	if value != nil && len(*value) < min {
		return business.NewFieldError(fieldName, MinLength)
	}
	return nil
}

// IsMaxLength returns a business.FieldError if the supplied string value has a length which is more than the supplied
// maximum
func IsMaxLength(fieldName string, value *string, max int) *business.FieldError {
	if value != nil && len(*value) > max {
		return business.NewFieldError(fieldName, MaxLength)
	}
	return nil
}

// IsDate returns a business.FieldError if the supplied string value is not of the expected date format.
func IsDate(fieldName string, value *string) (time.Time, *business.FieldError) {
	if value == nil {
		return time.Time{}, nil
	}
	parsed, err := parseDate(fieldName, *value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

// IsDateNowOrEarlier returns a business.FieldError if the supplied string represents a date that is later than now.
func IsDateNowOrEarlier(fieldName string, value time.Time) *business.FieldError {
	if value.After(time.Now()) {
		return business.NewFieldError(fieldName, DateInFuture)
	}
	return nil
}

// IsNotZero returns a business.FieldError if the supplied int value is zero.
func IsNotZero(fieldName string, value *int) *business.FieldError {
	if value != nil && *value == 0 {
		return business.NewFieldError(fieldName, ZeroValue)
	}
	return nil
}

// parseDate returns a date parsed using the configured date format or a business.FieldError if it does not match the format.
func parseDate(fieldName string, value string) (time.Time, *business.FieldError) {
	parsed, err := time.Parse(DateFormat, value)
	if err != nil {
		return time.Time{}, business.NewFieldError(fieldName, DateBadFormat)
	}
	return parsed, nil
}
