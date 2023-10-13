package business

import (
	"fmt"
)

// Reason is intended to be a meaningful, readable text id representing the error
type Reason string

// NewFieldError creates a FieldError with the supplied fieldName and reason.
func NewFieldError(fieldName string, reason Reason) *FieldError {
	return &FieldError{
		FieldName: fieldName,
		Reason:    reason,
	}
}

// FieldError represents a problem with some input provided by the user
type FieldError struct {
	FieldName string `json:"fieldName"`
	Reason    Reason `json:"reason"`
}

// Error implements the error interface on business.FieldError
func (e *FieldError) Error() string {
	return fmt.Sprintf("[field '%s', Reason: %s]", e.FieldName, e.Reason)
}

// Error represents a top level business error, with a collection of field errors and a message.
type Error struct {
	Fields  []FieldError `json:"fields,omitempty"`
	Message string       `json:"message"`
}

// Error implements the error interface on business.Error
func (e *Error) Error() string {
	return fmt.Sprintf("[message '%s', fields '%+v']", e.Message, e.Fields)
}
