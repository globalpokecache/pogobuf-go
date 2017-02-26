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
	"sync"
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
	options            Options
	player             Player
	rpc                *RPC
	request            int64
	lehmerSeed         int64
	rpcID              int64
	hasTicket          bool
	authTicket         *protos.AuthTicket
	token2             int
	sessionHash        []byte
	ptr8               string
	inventoryTimestamp int64
	templateTimestamp  int64
	startedTime        uint64
	serverURL          string
	firstGetMap        bool
	mapSettings        protos.MapSettings

	locationFixSync     sync.Mutex
	lastLocationCourse  float64
	lastLocationFixTime uint64
	lastLocationFix     *protos.Signature_LocationFix
	locationFixes       []*protos.Signature_LocationFix
	locationFixerStop   chan struct{}
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

	if opts.SignatureInfo.DeviceInfo == nil {
		opts.SignatureInfo.DeviceInfo = GetRandomDevice()
	}

	return &Instance{
		options:     *opts,
		token2:      token2,
		sessionHash: shash,
		startedTime: getTimestamp(time.Now().Add(time.Duration(-4500-randInt(1000)) * time.Millisecond)),
		rpc:         NewRPC(),
		firstGetMap: true,
		lehmerSeed:  1,
		ptr8:        "",
	}, nil
}

func (c *Instance) SetPosition(lat, lon, accu, alt float64) {
	c.player.Latitude = lat
	c.player.Longitude = lon
	if accu > 0 {
		c.player.Accuracy = accu
	}
	if alt > 0 {
		c.player.Altitude = alt
	}
}

func (c *Instance) BuildCommon() []*protos.Request {
	checkChallenge, _ := c.CheckChallengeRequest()
	getHatchedEggs, _ := c.GetHatchedEggsRequest()
	getInventory, _ := c.GetInventoryRequest(c.inventoryTimestamp)
	checkAwarded, _ := c.CheckAwardedBadgesRequest()
	downloadSettings, _ := c.DownloadSettingsRequest()
	// getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()

	return []*protos.Request{
		checkChallenge,
		getHatchedEggs,
		getInventory,
		checkAwarded,
		downloadSettings,
		// getBuddyWalkedReq,
	}
}

func (c *Instance) Init(ctx context.Context, nickname string) (*protos.GetPlayerResponse, error) {
	c.ptr8 = "90f6a704505bccac73cec99b07794993e6fd5a12"
	c.request = 1
	c.lehmerSeed = 1

	c.lastLocationFixTime = getTimestamp(time.Now())

	if c.locationFixerStop != nil {
		c.locationFixerStop <- struct{}{}
	}
	locationFixerStop := make(chan struct{})
	go c.locationFixer(locationFixerStop)
	c.locationFixerStop = locationFixerStop

	time.Sleep(time.Duration(1000+randInt(1000)) * time.Millisecond)

	c.Call(ctx)

	time.Sleep(1500 * time.Millisecond)

	var response *protos.ResponseEnvelope

	getPlayerReq, _ := c.GetPlayerRequest("US", "en", "America/Chicago")
	response, err := c.Call(ctx, getPlayerReq)
	if err != nil {
		return nil, err
	}

	var getPlayer protos.GetPlayerResponse
	err = proto.Unmarshal(response.Returns[0], &getPlayer)
	if err != nil {
		return nil, err
	}

	if getPlayer.Banned {
		return nil, pogobuf.ErrAccountBanned
	}

	time.Sleep(1500 * time.Millisecond)

	downloadRemoteConfigReq, _ := c.DownloadRemoteConfigVersionRequest(protos.Platform_IOS, c.options.Version)
	var requests = []*protos.Request{}
	requests = append(requests, downloadRemoteConfigReq)
	requests = append(requests, c.BuildCommon()...)
	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) < 6 {
		return nil, errors.New("Failed to initialize real player client")
	}

	var downloadRemoteConfig protos.DownloadRemoteConfigVersionResponse
	err = proto.Unmarshal(response.Returns[0], &downloadRemoteConfig)
	if err != nil {
		return nil, err
	}
	c.templateTimestamp = int64(downloadRemoteConfig.ItemTemplatesTimestampMs)

	var getInventory protos.GetInventoryResponse
	err = proto.Unmarshal(response.Returns[3], &getInventory)
	if err != nil {
		return nil, err
	}
	c.inventoryTimestamp = getInventory.InventoryDelta.NewTimestampMs

	var level int32
	for _, item := range getInventory.InventoryDelta.InventoryItems {
		if item.InventoryItemData.PlayerStats != nil {
			level = item.InventoryItemData.PlayerStats.Level
		}
	}

	var challengeResponse protos.CheckChallengeResponse
	err = proto.Unmarshal(response.Returns[1], &challengeResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to call CHECK_CHALLENGE: %s", err)
	}

	if challengeResponse.ShowChallenge {
		return nil, fmt.Errorf("CAPTCHA|%s", challengeResponse.ChallengeUrl)
	}

	var downloadResponse protos.DownloadSettingsResponse
	err = proto.Unmarshal(response.Returns[5], &downloadResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to call DOWNLOAD_SETTINGS: %s", err)
	}

	mapSettings := downloadResponse.GetSettings().GetMapSettings()
	if mapSettings != nil {
		c.mapSettings = *mapSettings
	}

	getAssetDigest, _ := c.GetAssetDigestRequest(protos.Platform_IOS, "", "", "", c.options.Version)
	requests = append(c.BuildCommon(), getAssetDigest)
	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) < 6 {
		return nil, errors.New("Failed to initialize real player client")
	}

	downloadItemTemplates, _ := c.DownloadItemTemplatesRequest(false, 0, 0)
	requests = append(c.BuildCommon(), downloadItemTemplates)
	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) < 6 {
		return nil, errors.New("Failed to initialize real player client")
	}

	if nickname != "" {
		err = c.completeTutorial(ctx, getPlayer.PlayerData.TutorialState, nickname)
		if err != nil {
			return nil, err
		}
	}

	getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()
	levelUpReq, _ := c.LevelUpRewardsRequest(level)
	requests = append(requests, levelUpReq)
	requests = append(requests, c.BuildCommon()...)
	requests = append(requests, getBuddyWalkedReq)
	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, err
	}

	return &getPlayer, nil
}

func (c *Instance) GetMap(ctx context.Context) (*protos.GetMapObjectsResponse, *protos.ResponseEnvelope, error) {
	cells := helpers.GetCellsFromRadius(c.player.Latitude, c.player.Longitude, 500, 15)
	var response *protos.ResponseEnvelope

	getMapReq, err := c.GetMapObjectsRequest(cells, make([]int64, len(cells)))
	if err != nil {
		return nil, response, err
	}

	getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()

	var requests []*protos.Request
	requests = append(requests, getMapReq)
	requests = append(requests, c.BuildCommon()...)
	requests = append(requests, getBuddyWalkedReq)

	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, response, err
	}

	if len(response.Returns) < len(requests) {
		return nil, response, errors.New("Server not accepted this request")
	}

	var getInventory protos.GetInventoryResponse
	err = proto.Unmarshal(response.Returns[3], &getInventory)
	if err != nil {
		return nil, response, err
	}
	c.inventoryTimestamp = getInventory.InventoryDelta.NewTimestampMs

	var challengeResponse protos.CheckChallengeResponse
	err = proto.Unmarshal(response.Returns[1], &challengeResponse)
	if err != nil {
		return nil, response, fmt.Errorf("Failed to call DOWNLOAD_SETTINGS: %s", err)
	}

	if challengeResponse.ShowChallenge {
		return nil, response, fmt.Errorf("CAPTCHA|%s", challengeResponse.ChallengeUrl)
	}

	var getMapObjects protos.GetMapObjectsResponse
	err = proto.Unmarshal(response.Returns[0], &getMapObjects)
	if err != nil {
		return nil, response, fmt.Errorf("Failed to parse GET_MAP_OBJECTS: %s", err)
	}

	debugProto("MapObjects", &getMapObjects)

	return &getMapObjects, response, nil
}

func (c Instance) MapSettings() protos.MapSettings {
	return c.mapSettings
}

func (c *Instance) SetAuthToken(authToken string) {
	c.options.AuthToken = authToken
	c.authTicket = nil
	c.hasTicket = false
}
