package errorhandler

import "fmt"

// DomainError represents a custom error type with a code, message, and optional underlying error.
type DomainError struct {
	Code    string
	Message string
	Err     error
}

// Error returns the error message of a DomainError instance. If Err is not nil, it appends the wrapped error message.
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error inside the DomainError, enabling error unwrapping for further inspection.
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new instance of DomainError with the provided code, message, and underlying error.
func NewDomainError(code string, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
