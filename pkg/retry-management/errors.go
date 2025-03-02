package retrymanagement

import "errors"

var (
	// ErrProfileNotFound indicates that the requested retry profile was not found
	ErrProfileNotFound = errors.New("retry profile not found")

	// ErrMaxAttemptsReached indicates that the maximum number of retry attempts has been reached
	ErrMaxAttemptsReached = errors.New("maximum retry attempts reached")
)
