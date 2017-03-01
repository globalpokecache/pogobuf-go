package client

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/globalpokecache/POGOProtos-go"
	"github.com/globalpokecache/pogobuf-go"
	"github.com/globalpokecache/pogobuf-go/auth"
	"github.com/globalpokecache/pogobuf-go/hash"
	"github.com/globalpokecache/pogobuf-go/helpers"
	"github.com/golang/protobuf/proto"
	"sync"
	"time"
)

type Options struct {
	Version       int
	SignatureInfo SignatureInfo
	AuthProvider  auth.Provider
	HashProvider  hash.Provider

	MaxTries             int
	MapObjectsMinDelay   time.Duration
	MapObjectsThrottling bool
}

var (
	defaultOptions = Options{
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
	lehmerSeed         int64
	rpcID              int64
	hasTicket          bool
	authToken          string
	authTicket         *protos.AuthTicket
	sessionHash        []byte
	ptr8               string
	inventoryTimestamp int64
	templateTimestamp  int64
	startedTime        uint64
	serverURL          string
	firstGetMap        bool
	mapSettings        protos.MapSettings

	locationFixSync     sync.Mutex
	lastLocationCourse  float32
	lastLocationFixTime uint64
	lastLocationFix     *protos.Signature_LocationFix
	locationFixes       []*protos.Signature_LocationFix
	locationFixerStop   chan struct{}
}

func New(opts *Options) (*Instance, error) {
	if opts.HashProvider == nil {
		return nil, errors.New("Missing Hash Provider")
	}

	if opts.AuthProvider == nil {
		return nil, errors.New("Missing Auth Provider")
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
		opts.SignatureInfo.DeviceInfo = NewDevice(opts.AuthProvider)
	}

	return &Instance{
		options: *opts,
		rpc:     NewRPC(),
	}, nil
}

func (c *Instance) BuildCommon(init bool) []*protos.Request {
	checkChallenge, _ := c.CheckChallengeRequest()
	getHatchedEggs, _ := c.GetHatchedEggsRequest()
	getInventory, _ := c.GetInventoryRequest(c.inventoryTimestamp)
	checkAwarded, _ := c.CheckAwardedBadgesRequest()
	downloadSettings, _ := c.DownloadSettingsRequest(downloadSettingsHash)
	getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()

	reqs := []*protos.Request{
		checkChallenge,
		getHatchedEggs,
		getInventory,
		checkAwarded,
	}

	if init {
		reqs = append(reqs, downloadSettings)
	} else {
		reqs = append(reqs, getBuddyWalkedReq)
	}

	return reqs
}

func (c *Instance) Init(ctx context.Context) (*protos.GetPlayerResponse, error) {
	shash := make([]byte, 16)
	_, err := rand.Read(shash)
	if err != nil {
		return nil, err
	}

	token, err := c.options.AuthProvider.Login(ctx)
	if err != nil {
		return nil, err
	}
	c.SetAuthToken(token)

	c.sessionHash = shash
	c.ptr8 = "90f6a704505bccac73cec99b07794993e6fd5a12"
	c.rpcID = 1
	c.lehmerSeed = 16807
	c.lastLocationFixTime = 0
	c.inventoryTimestamp = 0
	c.firstGetMap = true
	c.startedTime = getTimestamp(time.Now()) - uint64(5000+randInt(800))

	if c.locationFixerStop != nil {
		c.locationFixerStop <- struct{}{}
	}

	locationFixerStop := make(chan struct{})
	go c.locationFixer(locationFixerStop)
	c.locationFixerStop = locationFixerStop

	c.Call(ctx)

	var response *protos.ResponseEnvelope

	getPlayerReq, _ := c.GetPlayerRequest("US", "en", "America/Chicago")
	response, err = c.Call(ctx, getPlayerReq)
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
	downloadSettings, _ := c.DownloadSettingsRequest("")
	var requests = []*protos.Request{downloadRemoteConfigReq}
	requests = append(requests, c.BuildCommon(true)...)
	requests[5] = downloadSettings
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
	requests = []*protos.Request{getAssetDigest}
	requests = append(requests, c.BuildCommon(true)...)
	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) < 6 {
		return nil, errors.New("Failed to initialize real player client")
	}

	downloadItemTemplates, _ := c.DownloadItemTemplatesRequest(false, 0, 0)
	requests = []*protos.Request{downloadItemTemplates}
	requests = append(requests, c.BuildCommon(true)...)
	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) < 6 {
		return nil, errors.New("Failed to initialize real player client")
	}

	err = c.completeTutorial(ctx, getPlayer.PlayerData.TutorialState, c.options.AuthProvider.GetUsername())
	if err != nil {
		return nil, err
	}

	getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()
	levelUpReq, _ := c.LevelUpRewardsRequest(level)
	requests = append([]*protos.Request{}, levelUpReq)
	requests = append(requests, c.BuildCommon(true)...)
	requests = append(requests, getBuddyWalkedReq)
	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, err
	}

	c.CallWithPlatformRequests(ctx, nil, []*protos.RequestEnvelope_PlatformRequest{
		{
			Type: protos.PlatformRequestType_GET_STORE_ITEMS,
		},
	})

	return &getPlayer, nil
}

func (c *Instance) GetMap(ctx context.Context) (*protos.GetMapObjectsResponse, *protos.ResponseEnvelope, error) {
	cells := helpers.GetCellsFromRadius(c.player.Latitude, c.player.Longitude, 640, 15)
	var response *protos.ResponseEnvelope

	getMapReq, err := c.GetMapObjectsRequest(cells, make([]int64, len(cells)))
	if err != nil {
		return nil, response, err
	}

	requests := []*protos.Request{getMapReq}
	requests = append(requests, c.BuildCommon(false)...)

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

	// debugProto("MapObjects", &getMapObjects)

	return &getMapObjects, response, nil
}

func (c Instance) MapSettings() protos.MapSettings {
	return c.mapSettings
}

func (c *Instance) SetAuthToken(authToken string) {
	c.authToken = authToken
	c.authTicket = nil
	c.hasTicket = false
}
