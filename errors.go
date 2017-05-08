package pogobuf

import (
	"errors"
)

var (
	ErrServerDeniedRequest      = errors.New("Server denied your request. Seems your IP address is banned or something else really bad.")
	ErrServerUnexpectedResponse = errors.New("Unexpected response status code")
	ErrAuthExpired              = errors.New("Auth expired")
	ErrAccountBanned            = errors.New("Account is probably banned")
)
