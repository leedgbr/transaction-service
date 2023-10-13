package validation_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"transaction-service/internal/business"
	"transaction-service/internal/date"
	"transaction-service/internal/validation"
)

func TestIsRequiredString(t *testing.T) {
	tcs := []struct {
		name    string
		value   *string
		wantErr *business.FieldError
	}{
		{
			name:    "should not return validation error when the supplied value is a string that is not empty",
			value:   stringPtr("a"),
			wantErr: nil,
		},
		{
			name:    "should not return validation error when the supplied value is a long ish string",
			value:   stringPtr("abcdefghijklmnopqrstuvwxyz"),
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the supplied value is an empty string",
			value:   stringPtr(""),
			wantErr: nil,
		},
		{
			name:  "should return a validation error when the supplied value is nil",
			value: nil,
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("REQUIRED"),
			},
		},
	}
	for _, tc := range tcs {
		err := validation.IsRequiredString("*field-name*", tc.value)
		assert.Equal(t, tc.wantErr, err)
	}
}

func TestIsRequiredInt(t *testing.T) {
	tcs := []struct {
		name    string
		value   *int
		wantErr *business.FieldError
	}{
		{
			name:    "should not return a validation error when the value is zero",
			value:   intPtr(0),
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value is above zero",
			value:   intPtr(1),
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value is negative",
			value:   intPtr(-1),
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value is a large ish number",
			value:   intPtr(5000000),
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value is a large ish negative number",
			value:   intPtr(-5000000),
			wantErr: nil,
		},
		{
			name:  "should return a validation error when the value is nil",
			value: nil,
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("REQUIRED"),
			},
		},
	}
	for _, tc := range tcs {
		err := validation.IsRequiredInt("*field-name*", tc.value)
		assert.Equal(t, tc.wantErr, err)
	}
}

func TestIsMinLength(t *testing.T) {
	tcs := []struct {
		name    string
		value   *string
		min     int
		wantErr *business.FieldError
	}{
		{
			name:    "should not return a validation error when the value is nil",
			value:   nil,
			min:     10,
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value's length is at the minimum boundary",
			value:   stringPtr("1234567890"),
			min:     10,
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value's length is at the min boundary of zero",
			value:   stringPtr(""),
			min:     0,
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value's length is at the minimum boundary of one",
			value:   stringPtr("1"),
			min:     1,
			wantErr: nil,
		},
		{
			name:  "should return a validation error when the value's length is one less than the minimum boundary of one",
			value: stringPtr(""),
			min:   1,
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("MIN_LENGTH"),
			},
		},
		{
			name:  "should return a validation error when the value's length is one less than the minimum boundary of ten",
			value: stringPtr("123456789"),
			min:   10,
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("MIN_LENGTH"),
			},
		},
		{
			name:  "should return a validation error when the value's length is lots less than the minimum boundary",
			value: stringPtr("12"),
			min:   10,
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("MIN_LENGTH"),
			},
		},
	}
	for _, tc := range tcs {
		err := validation.IsMinLength("*field-name*", tc.value, tc.min)
		assert.Equal(t, tc.wantErr, err)
	}
}

func TestIsMaxLength(t *testing.T) {
	tcs := []struct {
		name    string
		value   *string
		max     int
		wantErr *business.FieldError
	}{
		{
			name:    "should not return a validation error when the value is nil",
			value:   nil,
			max:     20,
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value's length is at the maximum boundary",
			value:   stringPtr("12345678901234567890"),
			max:     20,
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value's length is at the max boundary of zero",
			value:   stringPtr(""),
			max:     0,
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value's length is at the maximum boundary of one",
			value:   stringPtr("1"),
			max:     1,
			wantErr: nil,
		},
		{
			name:  "should return a validation error when the value's length is one more than the maximum boundary of one",
			value: stringPtr("12"),
			max:   1,
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("MAX_LENGTH"),
			},
		},
		{
			name:  "should return a validation error when the value's length is one more than the maximum boundary of twenty",
			value: stringPtr("123456789012345678901"),
			max:   20,
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("MAX_LENGTH"),
			},
		},
		{
			name:  "should return a validation error when the value's length is lots more than the maximum boundary",
			value: stringPtr("12345678901234567890112345678"),
			max:   20,
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("MAX_LENGTH"),
			},
		},
	}
	for _, tc := range tcs {
		err := validation.IsMaxLength("*field-name*", tc.value, tc.max)
		assert.Equal(t, tc.wantErr, err)
	}
}

func TestIsDate(t *testing.T) {
	tcs := []struct {
		name     string
		value    *string
		wantErr  *business.FieldError
		wantDate time.Time
	}{
		{
			name:    "should not return validation error when the supplied value is nil",
			value:   nil,
			wantErr: nil,
		},
		{
			name:     "should not return validation error when the supplied value is a date in the correct format",
			value:    stringPtr("2023-10-25"),
			wantErr:  nil,
			wantDate: date.NewInUTC(2023, time.October, 25),
		},
		{
			name:  "should return a validation error when the supplied value is an empty string",
			value: stringPtr(""),
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("DATE_BAD_FORMAT"),
			},
		},
		{
			name:  "should return a validation error when the supplied value is a date that has too many days",
			value: stringPtr("2023-10-32"),
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("DATE_BAD_FORMAT"),
			},
		},
		{
			name:  "should return validation error when the supplied value is not a date",
			value: stringPtr("123456"),
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("DATE_BAD_FORMAT"),
			},
		},
	}
	for _, tc := range tcs {
		parsed, err := validation.IsDate("*field-name*", tc.value)
		assert.Equal(t, tc.wantErr, err)
		assert.Equal(t, tc.wantDate, parsed)
	}
}

func TestIsDateNowOrEarlier(t *testing.T) {
	tcs := []struct {
		name    string
		value   time.Time
		wantErr *business.FieldError
	}{
		{
			name:    "should not return validation error when the supplied value is empty",
			value:   time.Time{},
			wantErr: nil,
		},
		{
			name:    "should not return validation error when the supplied value is the current date",
			value:   time.Now(),
			wantErr: nil,
		},
		{
			name:    "should not return validation error when the supplied value is yesterday",
			value:   time.Now().AddDate(0, 0, -1),
			wantErr: nil,
		},
		{
			name:    "should not return validation error when the supplied value is last month",
			value:   time.Now().AddDate(0, -1, 0),
			wantErr: nil,
		},
		{
			name:    "should not return validation error when the supplied value is last year",
			value:   time.Now().AddDate(-1, 0, 0),
			wantErr: nil,
		},
		{
			name:  "should return validation error when the supplied value is tomorrow",
			value: time.Now().AddDate(0, 0, 1),
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("DATE_IN_FUTURE"),
			},
		},
		{
			name:  "should return validation error when the supplied value is a month in the future",
			value: time.Now().AddDate(0, 1, 0),
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("DATE_IN_FUTURE"),
			},
		},
		{
			name:  "should return validation error when the supplied value is a year in the future",
			value: time.Now().AddDate(1, 0, 0),
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("DATE_IN_FUTURE"),
			},
		},
	}
	for _, tc := range tcs {
		err := validation.IsDateNowOrEarlier("*field-name*", tc.value)
		assert.Equal(t, tc.wantErr, err)
	}
}

func TestIsNotZero(t *testing.T) {
	tcs := []struct {
		name    string
		value   *int
		wantErr *business.FieldError
	}{
		{
			name:    "should not return a validation error when the value is nil",
			value:   nil,
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value is just less than zero",
			value:   intPtr(-1),
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value is just more than zero",
			value:   intPtr(1),
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value is a lot less than zero",
			value:   intPtr(-12345678),
			wantErr: nil,
		},
		{
			name:    "should not return a validation error when the value is a lot more than zero",
			value:   intPtr(12345678),
			wantErr: nil,
		},
		{
			name:  "should return a validation error when the value is zero",
			value: intPtr(0),
			wantErr: &business.FieldError{
				FieldName: "*field-name*",
				Reason:    business.Reason("ZERO_VALUE"),
			},
		},
	}
	for _, tc := range tcs {
		err := validation.IsNotZero("*field-name*", tc.value)
		assert.Equal(t, tc.wantErr, err)
	}
}

func stringPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}
