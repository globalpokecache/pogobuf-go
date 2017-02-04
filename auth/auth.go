package auth

import (
	"context"
	"errors"

	"github.com/globalpokecache/pogobuf-go/auth/ptc"
)

// Provider is a common interface for managing auth tokens with the different third party authenticators
type Provider interface {
	Login(ctx context.Context, username, password string) (authToken string, err error)
}

// NewProvider creates a new provider based on the provider identifier
func NewProvider(provider string) (Provider, error) {
	switch provider {
	case "ptc":
		return ptc.NewProvider(), nil
	default:
		return nil, errors.New("Auth provider not supported")
	}
}
