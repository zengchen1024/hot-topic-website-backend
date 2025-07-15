package allerror

// overLimitedError
type overLimitedError struct {
	errorImpl
}

// OverLimit is a marker method for over limit rate error.
func (l overLimitedError) OverLimit() {}

// NewOverLimitError creates a new over limit error with the specified code and message.
func NewOverLimitError(msg string, err error) overLimitedError {
	return overLimitedError{errorImpl: New(errorCodeOverLimited, msg, err)}
}
