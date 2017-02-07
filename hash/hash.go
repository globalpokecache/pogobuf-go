package hash

import (
	"errors"
	"github.com/globalpokecache/pogobuf-go/hash/buddyauth"
)

type Provider interface {
	AddKey(string) error
	DelKey(string) error
	GetKeys() []interface{}
	Hash(authTicket, sessionData []byte, latitude, longitude, accuracy float64, timestamp uint64, requests [][]byte) (uint32, uint32, []uint64, error)
	SetDebug(bool)
}

func NewProvider(provider string, apiVersion int) (Provider, error) {
	switch provider {
	case "buddyauth":
		return buddyauth.NewProvider(apiVersion)
	default:
		return nil, errors.New("Hash provider not supported")
	}
}
