package client

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/globalpokecache/POGOProtos-go"
	"github.com/globalpokecache/pogobuf-go/hash"
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

func (c *Instance) Init(ctx context.Context) error {
	var response *protos.ResponseEnvelope
	c.Call(ctx)

	getPlayer, _ := c.GetPlayerRequest("", "", "")
	getHatchedEggs, _ := c.GetHatchedEggsRequest()
	getInventory, _ := c.GetInventoryRequest()
	checkAwarded, _ := c.CheckAwardedBadgesRequest()
	downloadSettings, _ := c.DownloadSettingsRequest()

	response, err := c.Call(ctx,
		getPlayer,
		getHatchedEggs,
		getInventory,
		checkAwarded,
		downloadSettings,
	)
	if err != nil {
		return errors.New("Failed to initialize real player client")
	}

	if len(response.Returns) < 5 {
		return errors.New("Failed to initialize real player client")
	}

	var downloadResponse protos.DownloadSettingsResponse
	err = proto.Unmarshal(response.Returns[4], &downloadResponse)
	if err != nil {
		return fmt.Errorf("Failed to call DOWNLOAD_SETTINGS: %s", err)
	}

	c.options.MapObjectsMinDelay = time.Duration(downloadResponse.GetSettings().GetMapSettings().GetMapObjectsMinRefreshSeconds) * time.Second

	return nil
}
