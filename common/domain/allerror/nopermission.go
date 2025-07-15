package allerror

import "errors"

// noPermissionError
type noPermissionError struct {
	errorImpl
}

// NoPermission is a marker method for a "no permission" error.
func (e noPermissionError) NoPermission() {}

// NewNoPermission creates a new "no permission" error with the specified message.
func NewNoPermission(msg string, err error) noPermissionError {
	return noPermissionError{errorImpl: New(errorCodeNoPermission, msg, err)}
}

// IsNoPermission checks if the given error is a "no permission" error.
func IsNoPermission(err error) bool {
	return errors.As(err, &noPermissionError{})
}
