package clientgroups

import "errors"

var (
	// ErrGroupExists indicates that a group with the same ID already exists
	ErrGroupExists = errors.New("group already exists")

	// ErrGroupNotFound indicates that the requested group was not found
	ErrGroupNotFound = errors.New("group not found")

	// ErrClientExists indicates that a client with the same ID already exists
	ErrClientExists = errors.New("client already exists")

	// ErrClientNotFound indicates that the requested client was not found
	ErrClientNotFound = errors.New("client not found")

	// ErrInvalidCredentials indicates that the provided credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrIPNotAllowed indicates that the client's IP is not in the allowed range
	ErrIPNotAllowed = errors.New("IP not allowed")

	// ErrMaxBindsExceeded indicates that the client has exceeded their maximum bind count
	ErrMaxBindsExceeded = errors.New("maximum binds exceeded")
)
