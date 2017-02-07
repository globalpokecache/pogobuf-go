package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"github.com/globalpokecache/POGOProtos-go"
	"github.com/globalpokecache/pogobuf-go"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context/ctxhttp"
)

const rpcUserAgent = "Niantic App"

func raise(message string) error {
	return fmt.Errorf("rpc/client: %s", message)
}

// RPC is used to communicate with the Pokémon Go API
type RPC struct {
	http *http.Client
}

// NewRPC constructs a Pokémon Go RPC API client
func NewRPC() *RPC {
	options := &cookiejar.Options{}
	jar, _ := cookiejar.New(options)
	httpClient := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return raise("Did not follow redirect")
		},
	}

	return &RPC{
		http: httpClient,
	}
}

// Request queries the Pokémon Go API will all pending requests
func (c *RPC) Request(ctx context.Context, endpoint string, requestEnvelope *protos.RequestEnvelope) (responseEnvelope *protos.ResponseEnvelope, err error) {
	responseEnvelope = &protos.ResponseEnvelope{}

	debugProto("RequestEnvelope", requestEnvelope)

	// Build request
	requestBytes, err := proto.Marshal(requestEnvelope)
	if err != nil {
		return responseEnvelope, raise("Could not encode request body")
	}
	requestReader := bytes.NewReader(requestBytes)
	request, err := http.NewRequest("POST", endpoint, requestReader)
	if err != nil {
		return responseEnvelope, raise("Unable to create the request")
	}
	request.Header.Add("User-Agent", rpcUserAgent)

	// Perform call to API
	response, err := ctxhttp.Do(ctx, c.http, request)
	if err != nil {
		return responseEnvelope, raise(fmt.Sprintf("There was an error requesting the API: %s", err))
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusForbidden {
		return nil, pogobuf.ErrServerDeniedRequest
	}

	if response.StatusCode != 200 {
		return responseEnvelope, pogobuf.ErrServerUnexpectedResponse
	}

	// Read the response
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return responseEnvelope, raise("Could not decode response body")
	}

	proto.Unmarshal(responseBytes, responseEnvelope)

	debugProto("ResponseEnvelope", responseEnvelope)

	if responseEnvelope.StatusCode == protos.ResponseEnvelope_BAD_REQUEST {
		return responseEnvelope, raise("Bad request")
	}

	if responseEnvelope.StatusCode == protos.ResponseEnvelope_INVALID_AUTH_TOKEN {
		return responseEnvelope, pogobuf.ErrAuthExpired
	}

	if responseEnvelope.StatusCode == protos.ResponseEnvelope_INVALID_PLATFORM_REQUEST {
		return responseEnvelope, raise("Invalid Platform Request")
	}

	return responseEnvelope, nil
}
