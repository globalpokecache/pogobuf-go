package client

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/globalpokecache/POGOProtos-go"
	"github.com/globalpokecache/pogobuf-go"
	"github.com/globalpokecache/pogobuf-go/hash"
	"github.com/globalpokecache/pogobuf-go/helpers"
	"github.com/golang/protobuf/proto"
	"time"
)

type Options struct {
	AuthToken     string
	AuthType      string
	Version       int
	SignatureInfo SignatureInfo
	HashProvider  hash.Provider

	MaxTries             int
	MapObjectsMinDelay   time.Duration
	MapObjectsThrottling bool
}

var (
	defaultOptions = Options{
		AuthToken:            "",
		AuthType:             "ptc",
		Version:              5500,
		SignatureInfo:        SignatureInfo{},
		HashProvider:         nil,
		MaxTries:             3,
		MapObjectsMinDelay:   5 * time.Second,
		MapObjectsThrottling: true,
	}
)

type Instance struct {
	options     Options
	player      Player
	rpc         *RPC
	rpcID       int64
	hasTicket   bool
	authTicket  *protos.AuthTicket
	token2      int
	sessionHash []byte
	startedTime time.Time
	serverURL   string
}

func New(opts *Options) (*Instance, error) {
	token2 := 1 + randInt(58)

	shash := make([]byte, 16)
	_, err := rand.Read(shash)
	if err != nil {
		return nil, err
	}

	if opts.HashProvider == nil {
		return nil, errors.New("Missing Hash Provider")
	}

	if opts.AuthType == "" {
		return nil, errors.New("Missing Auth Type")
	}

	if opts.Version != 0 {
		if _, err := getUnk25(opts.Version); err != nil {
			return nil, err
		}
	} else {
		opts.Version = defaultOptions.Version
	}

	if opts.MaxTries == 0 {
		opts.MaxTries = defaultOptions.MaxTries
	}

	return &Instance{
		options:     *opts,
		token2:      token2,
		sessionHash: shash,
		startedTime: time.Now(),
		rpc:         NewRPC(),
	}, nil
}

func (c *Instance) SetPosition(lat, lon, accu, alt float64) {
	c.player.Latitude = lat
	c.player.Longitude = lon
	c.player.Accuracy = accu
	c.player.Altitude = alt
}

func (c *Instance) buildCommon() []*protos.Request {
	getPlayer, _ := c.GetPlayerRequest("", "", "")
	getHatchedEggs, _ := c.GetHatchedEggsRequest()
	getInventory, _ := c.GetInventoryRequest()
	checkAwarded, _ := c.CheckAwardedBadgesRequest()
	checkChallenge, _ := c.CheckChallengeRequest()
	downloadSettings, _ := c.DownloadSettingsRequest()

	return []*protos.Request{
		getPlayer,
		getHatchedEggs,
		getInventory,
		checkAwarded,
		checkChallenge,
		downloadSettings,
	}
}

func (c *Instance) Init(ctx context.Context) error {
	var response *protos.ResponseEnvelope
	c.Call(ctx)

	response, err := c.Call(ctx,
		c.buildCommon()...,
	)
	if err != nil {
		return err
	}

	if err == pogobuf.ErrAuthExpired {
		return err
	}

	if len(response.Returns) < 6 {
		return errors.New("Failed to initialize real player client")
	}

	var challengeResponse protos.CheckChallengeResponse
	err = proto.Unmarshal(response.Returns[4], &challengeResponse)
	if err != nil {
		return fmt.Errorf("Failed to call DOWNLOAD_SETTINGS: %s", err)
	}

	if challengeResponse.ShowChallenge {
		return fmt.Errorf("CAPTCHA|%s", challengeResponse.ChallengeUrl)
	}

	var downloadResponse protos.DownloadSettingsResponse
	err = proto.Unmarshal(response.Returns[5], &downloadResponse)
	if err != nil {
		return fmt.Errorf("Failed to call DOWNLOAD_SETTINGS: %s", err)
	}

	c.options.MapObjectsMinDelay = time.Duration(downloadResponse.GetSettings().GetMapSettings().GetMapObjectsMinRefreshSeconds) * time.Second

	return nil
}

func (c *Instance) GetMap(ctx context.Context) (*protos.GetMapObjectsResponse, error) {
	cells := helpers.GetCellsFromRadius(c.player.Latitude, c.player.Longitude, 210, 17)

	request, err := c.GetMapObjectsRequest(cells, make([]int64, len(cells)))
	if err != nil {
		return nil, err
	}

	requests := c.buildCommon()
	requests = append(requests, request)

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) < len(requests) {
		return nil, errors.New("Server not accepted this request")
	}

	var challengeResponse protos.CheckChallengeResponse
	err = proto.Unmarshal(response.Returns[4], &challengeResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse CHECK_CHALLENGE: %s", err)
	}

	if challengeResponse.ShowChallenge {
		return nil, fmt.Errorf("CAPTCHA|%s", challengeResponse.ChallengeUrl)
	}

	var getMapObjects protos.GetMapObjectsResponse
	err = proto.Unmarshal(response.Returns[len(requests)-1], &getMapObjects)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse GET_MAP_OBJECTS: %s", err)
	}

	debugProto("MapObjects", &getMapObjects)

	return &getMapObjects, nil
}

func (c *Instance) SetAuthToken(authToken string) {
	c.options.AuthToken = authToken
}
