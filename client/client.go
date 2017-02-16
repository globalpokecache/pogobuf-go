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
	options            Options
	player             Player
	rpc                *RPC
	rpcID              int64
	hasTicket          bool
	authTicket         *protos.AuthTicket
	token2             int
	sessionHash        []byte
	inventoryTimestamp int64
	templateTimestamp  int64
	startedTime        time.Time
	serverURL          string
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

func (c *Instance) BuildCommon() []*protos.Request {
	checkChallenge, _ := c.CheckChallengeRequest()
	getHatchedEggs, _ := c.GetHatchedEggsRequest()
	getInventory, _ := c.GetInventoryRequest(c.inventoryTimestamp)
	checkAwarded, _ := c.CheckAwardedBadgesRequest()
	downloadSettings, _ := c.DownloadSettingsRequest()

	return []*protos.Request{
		checkChallenge,
		getHatchedEggs,
		getInventory,
		checkAwarded,
		downloadSettings,
	}
}

func (c *Instance) Init(ctx context.Context, account string) (*protos.GetPlayerResponse, error) {
	c.inventoryTimestamp = 0

	var response *protos.ResponseEnvelope
	c.Call(ctx)

	time.Sleep(1500 * time.Millisecond)

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

	c.options.MapObjectsMinDelay = time.Duration(downloadResponse.GetSettings().GetMapSettings().GetMapObjectsMinRefreshSeconds) * time.Second

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

	err = c.completeTutorial(ctx, getPlayer.PlayerData.TutorialState, account)
	if err != nil {
		return nil, err
	}

	levelUp, _ := c.LevelUpRewardsRequest(level)
	requests = append(c.BuildCommon(), levelUp)
	response, err = c.Call(ctx, requests...)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) < 6 {
		return nil, errors.New("Failed to initialize real player client")
	}

	return &getPlayer, nil
}

var tutorialRequirements = []int32{0, 1, 3, 4, 7}

func (c *Instance) completeTutorial(ctx context.Context, tutorialState []protos.TutorialState, account string) error {
	completed := 0
	tuto := map[int32]bool{}
	for _, t := range tutorialState {
		for _, req := range tutorialRequirements {
			if req == int32(t) {
				tuto[req] = true
				completed++
			}
		}
	}

	getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()

	if completed == 5 {
		getPlayerProfile, err := c.GetPlayerProfileRequest("")
		if err != nil {
			return err
		}
		var requests []*protos.Request
		requests = append(requests, getPlayerProfile)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		registerBackground, err := c.RegisterBackgroundDeviceRequest("", "apple_watch")
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, registerBackground)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		return nil
	}

	if _, ok := tuto[0]; !ok {
		time.Sleep(time.Duration(2+randInt(3)) * time.Second)
		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{0}, false, false)
		if err != nil {
			return err
		}
		requests := []*protos.Request{}
		requests = append(requests, markComplete)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	if _, ok := tuto[1]; !ok {
		time.Sleep(time.Duration(8+randInt(7)) * time.Second)
		setAvatar, err := c.SetAvatarRequest(
			randInt(3),
			randInt(5),
			randInt(3),
			randInt(2),
			randInt(4),
			randInt(6),
			0,
			randInt(4),
			randInt(5),
		)
		if err != nil {
			return err
		}
		requests := []*protos.Request{}
		requests = append(requests, setAvatar)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(1+randInt(1)) * time.Second)

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{1}, false, false)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, markComplete)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	getPlayerProfile, err := c.GetPlayerProfileRequest("")
	if err != nil {
		return err
	}
	var requests []*protos.Request
	requests = append(requests, getPlayerProfile)
	requests = append(requests, c.BuildCommon()...)
	requests = append(requests, getBuddyWalkedReq)
	_, err = c.Call(ctx, requests...)
	if err != nil {
		return err
	}

	registerBackground, err := c.RegisterBackgroundDeviceRequest("", "apple_watch")
	if err != nil {
		return err
	}
	requests = []*protos.Request{}
	requests = append(requests, registerBackground)
	requests = append(requests, c.BuildCommon()...)
	_, err = c.Call(ctx, requests...)
	if err != nil {
		return err
	}

	if _, ok := tuto[3]; !ok {
		getDownloadsURLs, err := c.GetDownloadURLsRequest([]string{
			"1a3c2816-65fa-4b97-90eb-0b301c064b7a/1477084786906000",
			"e89109b0-9a54-40fe-8431-12f7826c8194/1477084802881000",
		})
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, getDownloadsURLs)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(7+randInt(3)) * time.Second)
		crea := []int32{1, 4, 7}[randInt(3)]

		encounterRequest, err := c.EncounterTutorialCompleteRequest(crea)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, encounterRequest)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		getPlayerRequest, err := c.GetPlayerRequest("US", "en", "America/Chicago")
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, getPlayerRequest)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	if _, ok := tuto[4]; !ok {
		time.Sleep(time.Duration(5+randInt(7)) * time.Second)

		claimCodename, err := c.ClaimCodenameRequest(account)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, claimCodename)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{4}, false, false)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, markComplete)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	if _, ok := tuto[7]; !ok {
		time.Sleep(time.Duration(4+randInt(3)) * time.Second)

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{7}, false, false)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, markComplete)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	return nil
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

func (c *Instance) SetAuthToken(authToken string) {
	c.options.AuthToken = authToken
	c.authTicket = nil
	c.hasTicket = false
}
