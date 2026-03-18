package shared

import "fmt"

type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (e *DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("[%s] %s (field: %s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

func NewValidationError(code, message, field string) *DomainError {
	return &DomainError{Code: code, Message: message, Field: field}
}

// Common errors
var (
	ErrNotFound     = NewDomainError("NOT_FOUND", "resource not found")
	ErrUnauthorized = NewDomainError("UNAUTHORIZED", "unauthorized")
	ErrForbidden    = NewDomainError("FORBIDDEN", "forbidden")
)
