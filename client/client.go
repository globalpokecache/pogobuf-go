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
	Version              int
	SignatureInfo        SignatureInfo
	AuthProvider         auth.Provider
	HashProvider         hash.Provider
	SimulateApp          bool
	AutoCompleteTutorial bool
	GoogleMapsKey        string

	MaxTries             int
	MapObjectsMinDelay   time.Duration
	MinRequestInterval   time.Duration
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
		SimulateApp:          false,
		AutoCompleteTutorial: false,
		MinRequestInterval:   500 * time.Millisecond,
	}

	DefaultPtr8             = "90f6a704505bccac73cec99b07794993e6fd5a12"
	DefaultLehmerSeed int64 = 16807
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
	waitRequest        chan struct{}

	locationFixSync     sync.Mutex
	lastLocationCourse  float32
	lastLocationFixTime uint64
	lastLocationFix     *protos.Signature_LocationFix
	locationFixes       []*protos.Signature_LocationFix

	cancel func()
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
		options:     *opts,
		rpc:         NewRPC(),
		lehmerSeed:  DefaultLehmerSeed,
		ptr8:        DefaultPtr8,
		waitRequest: make(chan struct{}),
	}, nil
}

func (c *Instance) BuildCommon(init bool) []*protos.Request {
	checkChallenge, _ := c.CheckChallengeRequest()
	getHatchedEggs, _ := c.GetHatchedEggsRequest()
	getInventory, _ := c.GetInventoryRequest(c.inventoryTimestamp)
	checkAwarded, _ := c.CheckAwardedBadgesRequest()
	downloadSettings, _ := c.DownloadSettingsRequest(downloadSettingsHash)

	reqs := []*protos.Request{
		checkChallenge,
		getHatchedEggs,
		getInventory,
		checkAwarded,
	}

	if init {
		reqs = append(reqs, downloadSettings)
	}

	return reqs
}

func (c *Instance) simulateAppLogin(ctx context.Context) (*protos.GetPlayerResponse, error) {
	c.Call(ctx)

	randSleep(430, 970)

	getPlayerReq, _ := c.GetPlayerRequest("US", "en", "America/Chicago")
	response, err := c.Call(ctx, getPlayerReq)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Failed to initialize account")
	}

	var getPlayer protos.GetPlayerResponse
	err = proto.Unmarshal(response.Returns[0], &getPlayer)
	if err != nil {
		return nil, err
	}

	if getPlayer.Banned {
		return nil, pogobuf.ErrAccountBanned
	}

	randSleep(530, 1000)

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

	var level *int32
	for _, item := range getInventory.InventoryDelta.InventoryItems {
		if item.InventoryItemData.PlayerStats != nil {
			level = &item.InventoryItemData.PlayerStats.Level
		}
	}

	var challengeResponse protos.CheckChallengeResponse
	err = proto.Unmarshal(response.Returns[1], &challengeResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal CHECK_CHALLENGE: %s", err)
	}

	if challengeResponse.ShowChallenge {
		return nil, fmt.Errorf("CAPTCHA|%s", challengeResponse.ChallengeUrl)
	}

	var downloadResponse protos.DownloadSettingsResponse
	err = proto.Unmarshal(response.Returns[5], &downloadResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal DOWNLOAD_SETTINGS: %s", err)
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

	randSleep(870, 2000)

	getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()

	var alreadyComplete bool
	if c.options.AutoCompleteTutorial {
		var assets protos.GetAssetDigestResponse
		err = proto.Unmarshal(response.Returns[0], &assets)
		if err != nil {
			return nil, fmt.Errorf("Failed to call GE: %s", err)
		}

		var assetsIds []string
		for _, asset := range assets.Digest {
			if asset.BundleName == "pm0001" ||
				asset.BundleName == "pm0004" ||
				asset.BundleName == "pm0007" {
				assetsIds = append(assetsIds, asset.AssetId)
			}
		}
		alreadyComplete, err = c.completeTutorial(ctx, getPlayer.PlayerData.TutorialState, c.options.AuthProvider.GetUsername(), assetsIds)
		if err != nil {
			return nil, err
		}
	}

	if !c.options.AutoCompleteTutorial || alreadyComplete {
		getPlayerProfile, err := c.GetPlayerProfileRequest(c.options.AuthProvider.GetUsername())
		if err != nil {
			return nil, err
		}
		requests = []*protos.Request{getPlayerProfile}
		requests = append(requests, c.BuildCommon(true)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return nil, err
		}
		randSleep(200, 400)

		if level != nil {
			levelUpReq, _ := c.LevelUpRewardsRequest(*level)
			requests = append([]*protos.Request{}, levelUpReq)
			requests = append(requests, c.BuildCommon(true)...)
			requests = append(requests, getBuddyWalkedReq)
			_, err = c.Call(ctx, requests...)
			if err != nil {
				return nil, err
			}
			randSleep(450, 700)
		}

		regBg, _ := c.RegisterBackgroundDeviceRequest("", "apple_watch")
		requests = append([]*protos.Request{}, regBg)
		requests = append(requests, c.BuildCommon(true)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return nil, err
		}
		randSleep(500, 1300)
	}

	return &getPlayer, nil
}

func (c *Instance) minimalLogin(ctx context.Context) (*protos.GetPlayerResponse, error) {
	getPlayerReq, _ := c.GetPlayerRequest("US", "en", "America/Chicago")
	response, err := c.Call(ctx, getPlayerReq)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Failed to initialize account")
	}

	var getPlayer protos.GetPlayerResponse
	err = proto.Unmarshal(response.Returns[0], &getPlayer)
	if err != nil {
		return nil, err
	}

	if getPlayer.Banned {
		return nil, pogobuf.ErrAccountBanned
	}

	return &getPlayer, nil
}

func (c *Instance) newSessionHash() error {
	shash := make([]byte, 16)
	_, err := rand.Read(shash)
	if err != nil {
		return err
	}
	c.sessionHash = shash
	return nil
}

func (c *Instance) login(ctx context.Context) error {
	token, err := c.options.AuthProvider.Login(ctx)
	if err != nil {
		return err
	}
	c.SetAuthToken(token)
	return nil
}

func (c *Instance) Init(ctx context.Context) (*protos.GetPlayerResponse, error) {
	if c.cancel != nil {
		c.cancel()
	}

	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.login(ctx)

	err := c.newSessionHash()
	if err != nil {
		return nil, err
	}

	c.ptr8 = DefaultPtr8
	c.rpcID = 1
	c.lehmerSeed = DefaultLehmerSeed
	c.lastLocationFixTime = 0
	c.inventoryTimestamp = 0
	c.firstGetMap = true
	c.startedTime = getTimestamp(time.Now()) - uint64(5000+randInt(800))

	go c.locationFixer(ctx)
	go c.requestThrottle(ctx)

	if c.options.SimulateApp {
		return c.simulateAppLogin(ctx)
	}

	return c.minimalLogin(ctx)
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

	if DebugGMO {
		debugProto("MapObjects", &getMapObjects)
	}

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
