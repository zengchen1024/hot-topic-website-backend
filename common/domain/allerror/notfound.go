package allerror

import "errors"

// NewNotFoundError creates a new not found error with the specified code and message.
func NewNotFoundError(msg string, err error) notfoudError {
	return notfoudError{errorImpl: New(errorCodeNotFound, msg, err)}
}

// IsNotFoundError checks if the given error is a not found error.
func IsNotFoundError(err error) bool {
	return errors.As(err, &notfoudError{})
}

// notfoudError
type notfoudError struct {
	errorImpl
}

// NotFound is a marker method for a not found error.
func (e notfoudError) NotFound() {}
