package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/globalpokecache/POGOProtos-go"
	"github.com/golang/protobuf/proto"
)

func (c *Instance) DownloadRemoteConfigVersionRequest(platform protos.Platform, appVersion int) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.DownloadRemoteConfigVersionMessage{
		Platform:   platform,
		AppVersion: uint32(appVersion),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create DOWNLOAD_REMOTE_CONFIG_VERSION: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_DOWNLOAD_REMOTE_CONFIG_VERSION,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) DownloadRemoteConfigVersion(ctx context.Context, platform protos.Platform, appVersion int) (*protos.DownloadSettingsResponse, error) {
	request, err := c.DownloadRemoteConfigVersionRequest(platform, appVersion)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var downloadRemote protos.DownloadSettingsResponse
	err = proto.Unmarshal(response.Returns[0], &downloadRemote)
	if err != nil {
		return nil, fmt.Errorf("Failed to call DOWNLOAD_REMOTE_CONFIG_VERSION: %s", err)
	}

	return &downloadRemote, nil
}

func (c *Instance) GetAssetDigestRequest(platform protos.Platform, manufacturer, model, locale string, appVersion int) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.GetAssetDigestMessage{
		Platform:           protos.Platform_IOS,
		DeviceManufacturer: "",
		DeviceModel:        "",
		Locale:             "",
		AppVersion:         uint32(appVersion),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create GET_ASSET_DIGEST: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_GET_ASSET_DIGEST,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) GetAssetDigest(ctx context.Context, platform protos.Platform, manufacturer, model, locale string, appVersion int) (*protos.GetAssetDigestResponse, error) {
	request, err := c.GetAssetDigestRequest(platform, manufacturer, model, locale, appVersion)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var getAssetDigest protos.GetAssetDigestResponse
	err = proto.Unmarshal(response.Returns[0], &getAssetDigest)
	if err != nil {
		return nil, fmt.Errorf("Failed to call GET_ASSET_DIGEST: %s", err)
	}

	return &getAssetDigest, nil
}

func (c *Instance) DownloadItemTemplatesRequest(paginate bool, offset int32, ts uint64) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.DownloadItemTemplatesMessage{
		Paginate:      paginate,
		PageOffset:    offset,
		PageTimestamp: ts,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create DOWNLOAD_ITEM_TEMPLATES: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_DOWNLOAD_ITEM_TEMPLATES,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) DownloadItemTemplates(ctx context.Context, paginate bool, offset int32, ts uint64) (*protos.DownloadItemTemplatesResponse, error) {
	request, err := c.DownloadItemTemplatesRequest(paginate, offset, ts)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var downloadItemTemplates protos.DownloadItemTemplatesResponse
	err = proto.Unmarshal(response.Returns[0], &downloadItemTemplates)
	if err != nil {
		return nil, fmt.Errorf("Failed to call DOWNLOAD_ITEM_TEMPLATES: %s", err)
	}

	return &downloadItemTemplates, nil
}

func (c *Instance) DownloadSettingsRequest() (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.DownloadSettingsMessage{
		Hash: downloadSettingsHash,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create DOWNLOAD_SETTINGS: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_DOWNLOAD_SETTINGS,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) DownloadSettings(ctx context.Context) (*protos.DownloadSettingsResponse, error) {
	request, err := c.DownloadSettingsRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var downloadSettings protos.DownloadSettingsResponse
	err = proto.Unmarshal(response.Returns[0], &downloadSettings)
	if err != nil {
		return nil, fmt.Errorf("Failed to call DOWNLOAD_SETTINGS: %s", err)
	}

	return &downloadSettings, nil
}

func (c *Instance) PlayerUpdateRequest() (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.PlayerUpdateMessage{
		Latitude:  c.player.Latitude,
		Longitude: c.player.Longitude,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create PLAYER_UPDATE: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_PLAYER_UPDATE,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) PlayerUpdate(ctx context.Context) (*protos.PlayerUpdateResponse, error) {
	request, err := c.PlayerUpdateRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var playerUpdate protos.PlayerUpdateResponse
	err = proto.Unmarshal(response.Returns[0], &playerUpdate)
	if err != nil {
		return nil, fmt.Errorf("Failed to call PLAYER_UPDATE: %s", err)
	}

	return &playerUpdate, nil
}

func (c *Instance) GetPlayerRequest(country, language, timezone string) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.GetPlayerMessage{
		PlayerLocale: &protos.GetPlayerMessage_PlayerLocale{
			Country:  country,
			Language: language,
			Timezone: timezone,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create GET_PLAYER: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_GET_PLAYER,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) GetPlayer(ctx context.Context, country, language, timezone string) (*protos.GetPlayerResponse, error) {
	request, err := c.GetPlayerRequest(country, language, timezone)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var getPlayer protos.GetPlayerResponse
	err = proto.Unmarshal(response.Returns[0], &getPlayer)
	if err != nil {
		return nil, fmt.Errorf("Failed to call GET_PLAYER: %s", err)
	}

	return &getPlayer, nil
}

func (c *Instance) GetHatchedEggsRequest() (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.GetHatchedEggsResponse{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create GET_HATCHED_EGGS: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_GET_HATCHED_EGGS,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) GetHatchedEggs(ctx context.Context) (*protos.GetHatchedEggsResponse, error) {
	request, err := c.GetHatchedEggsRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var getHatchedEggs protos.GetHatchedEggsResponse
	err = proto.Unmarshal(response.Returns[0], &getHatchedEggs)
	if err != nil {
		return nil, fmt.Errorf("Failed to call GET_HATCHED_EGGS: %s", err)
	}

	return &getHatchedEggs, nil
}

func (c *Instance) GetInventoryRequest(last int64) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.GetInventoryMessage{
		LastTimestampMs: last,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create GET_INVENTORY: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_GET_INVENTORY,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) GetInventory(ctx context.Context, last int64) (*protos.GetInventoryResponse, error) {
	request, err := c.GetInventoryRequest(last)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var getInventory protos.GetInventoryResponse
	err = proto.Unmarshal(response.Returns[0], &getInventory)
	if err != nil {
		return nil, fmt.Errorf("Failed to call GET_INVENTORY: %s", err)
	}

	return &getInventory, nil
}

func (c *Instance) CheckAwardedBadgesRequest() (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.CheckAwardedBadgesMessage{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create CHECK_AWARDED_BADGES: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_CHECK_AWARDED_BADGES,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) CheckAwardedBadges(ctx context.Context) (*protos.CheckAwardedBadgesResponse, error) {
	request, err := c.CheckAwardedBadgesRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var checkAwardedBadges protos.CheckAwardedBadgesResponse
	err = proto.Unmarshal(response.Returns[0], &checkAwardedBadges)
	if err != nil {
		return nil, fmt.Errorf("Failed to call CHECK_AWARDED_BADGES: %s", err)
	}

	return &checkAwardedBadges, nil
}

func (c *Instance) GetMapObjectsRequest(cellIDs []uint64, sinceTimestampMs []int64) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.GetMapObjectsMessage{
		CellId:           cellIDs,
		SinceTimestampMs: sinceTimestampMs,
		Latitude:         c.player.Latitude,
		Longitude:        c.player.Longitude,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create GET_MAP_OBJECTS: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_GET_MAP_OBJECTS,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) GetMapObjects(ctx context.Context, cellIDs []uint64, sinceTimestampMs []int64) (*protos.GetMapObjectsResponse, error) {
	request, err := c.GetMapObjectsRequest(cellIDs, sinceTimestampMs)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var getMapObjects protos.GetMapObjectsResponse
	err = proto.Unmarshal(response.Returns[0], &getMapObjects)
	if err != nil {
		return nil, fmt.Errorf("Failed to call GET_MAP_OBJECTS: %s", err)
	}

	debugProto("MapObjects", &getMapObjects)

	return &getMapObjects, nil
}

func (c *Instance) CheckChallengeRequest() (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.CheckChallengeMessage{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create CHECK_CHALLENGE: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_CHECK_CHALLENGE,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) CheckChallenge(ctx context.Context) (*protos.CheckChallengeResponse, error) {
	request, err := c.CheckChallengeRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var checkChallange protos.CheckChallengeResponse
	err = proto.Unmarshal(response.Returns[0], &checkChallange)
	if err != nil {
		return nil, fmt.Errorf("Failed to call CHECK_CHALLENGE: %s", err)
	}

	return &checkChallange, nil
}

func (c *Instance) VerifyChallengeRequest(token string) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.VerifyChallengeMessage{
		Token: token,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create VERIFY_CHALLENGE: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_VERIFY_CHALLENGE,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) VerifyChallenge(ctx context.Context, token string) (*protos.VerifyChallengeResponse, error) {
	request, err := c.VerifyChallengeRequest(token)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var verifyChallange protos.VerifyChallengeResponse
	err = proto.Unmarshal(response.Returns[0], &verifyChallange)
	if err != nil {
		return nil, fmt.Errorf("Failed to call VERIFY_CHALLENGE: %s", err)
	}

	return &verifyChallange, nil
}

func (c *Instance) GetBuddyWalkedRequest() (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.GetBuddyWalkedMessage{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create GET_BUDDY_WALKED: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_GET_BUDDY_WALKED,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) GetBuddyWalked(ctx context.Context) (*protos.GetBuddyWalkedResponse, error) {
	request, err := c.GetBuddyWalkedRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var getBuddyWalked protos.GetBuddyWalkedResponse
	err = proto.Unmarshal(response.Returns[0], &getBuddyWalked)
	if err != nil {
		return nil, fmt.Errorf("Failed to call GET_BUDDY_WALKED: %s", err)
	}

	return &getBuddyWalked, nil
}

func (c *Instance) EncounterRequest(eid uint64, spawnPoint string) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.EncounterMessage{
		EncounterId:     eid,
		SpawnPointId:    spawnPoint,
		PlayerLatitude:  c.player.Latitude,
		PlayerLongitude: c.player.Longitude,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create ENCOUNTER: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_ENCOUNTER,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) Encounter(ctx context.Context, eid uint64, spawnPoint string) (*protos.EncounterResponse, error) {
	request, err := c.EncounterRequest(eid, spawnPoint)
	if err != nil {
		return nil, err
	}

	var encounter protos.EncounterResponse

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		err = errors.New("Server not accepted this request")
		return nil, err
	}

	err = proto.Unmarshal(response.Returns[0], &encounter)
	if err != nil {
		err = fmt.Errorf("Failed to call ENCOUNTER: %s", err)
		return nil, err
	}

	return &encounter, nil
}

func (c *Instance) CatchPokemonRequest(eid uint64, spawnPoint string, iid protos.ItemId, nrs float64, nhp float64, hit bool, spin float64) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.CatchPokemonMessage{
		EncounterId:           eid,
		SpawnPointId:          spawnPoint,
		Pokeball:              iid,
		NormalizedReticleSize: nrs,
		NormalizedHitPosition: nhp,
		HitPokemon:            hit,
		SpinModifier:          spin,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create CATCH_POKEMON: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_CATCH_POKEMON,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) CatchPokemon(ctx context.Context, eid uint64, spawnPoint string, iid protos.ItemId, nrs float64, nhp float64, hit bool, spin float64) (*protos.CatchPokemonResponse, error) {
	request, err := c.CatchPokemonRequest(eid, spawnPoint, iid, nrs, nhp, hit, spin)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var catchResult protos.CatchPokemonResponse
	err = proto.Unmarshal(response.Returns[0], &catchResult)
	if err != nil {
		return nil, fmt.Errorf("Failed to call CATCH_POKEMON: %s", err)
	}

	return &catchResult, nil
}

func (c *Instance) ReleasePokemonRequest(pid uint64) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.ReleasePokemonMessage{
		PokemonId: pid,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create RELEASE_POKEMON: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_RELEASE_POKEMON,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) ReleasePokemon(ctx context.Context, pid uint64) (*protos.ReleasePokemonResponse, error) {
	request, err := c.ReleasePokemonRequest(pid)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var release protos.ReleasePokemonResponse
	err = proto.Unmarshal(response.Returns[0], &release)
	if err != nil {
		return nil, fmt.Errorf("Failed to call RELEASE_POKEMON: %s", err)
	}

	return &release, nil
}

func (c *Instance) ReleaseMultiPokemonRequest(pids []uint64) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.ReleasePokemonMessage{
		PokemonIds: pids,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create RELEASE_POKEMON: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_RELEASE_POKEMON,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) ReleaseMultiPokemon(ctx context.Context, pids []uint64) (*protos.ReleasePokemonResponse, error) {
	request, err := c.ReleaseMultiPokemonRequest(pids)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var release protos.ReleasePokemonResponse
	err = proto.Unmarshal(response.Returns[0], &release)
	if err != nil {
		return nil, fmt.Errorf("Failed to call RELEASE_POKEMON: %s", err)
	}

	return &release, nil
}

func (c *Instance) FortSearchRequest(fortid string, lat, lon float64) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.FortSearchMessage{
		FortId:          fortid,
		PlayerLatitude:  c.player.Latitude,
		PlayerLongitude: c.player.Longitude,
		FortLatitude:    lat,
		FortLongitude:   lon,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create FORT_SEARCH: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_FORT_SEARCH,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) FortSearch(ctx context.Context, fortid string, lat, lon float64) (*protos.FortSearchResponse, error) {
	request, err := c.FortSearchRequest(fortid, lat, lon)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var search protos.FortSearchResponse
	err = proto.Unmarshal(response.Returns[0], &search)
	if err != nil {
		return nil, fmt.Errorf("Failed to call FORT_SEARCH: %s", err)
	}

	return &search, nil
}

func (c *Instance) LevelUpRewardsRequest(level int32) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.LevelUpRewardsMessage{
		Level: level,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create LEVEL_UP_REWARDS: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_LEVEL_UP_REWARDS,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) LevelUpRewards(ctx context.Context, level int32) (*protos.LevelUpRewardsResponse, error) {
	request, err := c.LevelUpRewardsRequest(level)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var levelup protos.LevelUpRewardsResponse
	err = proto.Unmarshal(response.Returns[0], &levelup)
	if err != nil {
		return nil, fmt.Errorf("Failed to call GET_PLAYER_PROFILE: %s", err)
	}

	return &levelup, nil
}

func (c *Instance) GetPlayerProfileRequest(name string) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.GetPlayerProfileMessage{
		PlayerName: name,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create GET_PLAYERP_PROFILE: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_GET_PLAYER_PROFILE,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) GetPlayerProfile(ctx context.Context, name string) (*protos.GetPlayerProfileResponse, error) {
	request, err := c.GetPlayerProfileRequest(name)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var playerProfile protos.GetPlayerProfileResponse
	err = proto.Unmarshal(response.Returns[0], &playerProfile)
	if err != nil {
		return nil, fmt.Errorf("Failed to call GET_PLAYER_PROFILE: %s", err)
	}

	return &playerProfile, nil
}

func (c *Instance) RegisterBackgroundDeviceRequest(device string, devicetype string) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.RegisterBackgroundDeviceMessage{
		DeviceId:   device,
		DeviceType: devicetype,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create REGISTER_BACKGROUND_DEVICE: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_REGISTER_BACKGROUND_DEVICE,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) RegisterBackgroundDevice(ctx context.Context, device string, devicetype string) (*protos.RegisterBackgroundDeviceResponse, error) {
	request, err := c.RegisterBackgroundDeviceRequest(device, devicetype)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var registerDevice protos.RegisterBackgroundDeviceResponse
	err = proto.Unmarshal(response.Returns[0], &registerDevice)
	if err != nil {
		return nil, fmt.Errorf("Failed to call REGISTER_BACKGROUND_DEVICE: %s", err)
	}

	return &registerDevice, nil
}

func (c *Instance) MarkTutorialCompleteRequest(ids []protos.TutorialState, sendMail, sendNotif bool) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.MarkTutorialCompleteMessage{
		TutorialsCompleted:    ids,
		SendMarketingEmails:   sendMail,
		SendPushNotifications: sendNotif,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create MARK_TUTORIAL_COMPLETE: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_MARK_TUTORIAL_COMPLETE,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) MarkTutorialComplete(ctx context.Context, ids []protos.TutorialState, sendMail, sendNotif bool) (*protos.MarkTutorialCompleteResponse, error) {
	request, err := c.MarkTutorialCompleteRequest(ids, sendMail, sendNotif)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var markTutorial protos.MarkTutorialCompleteResponse
	err = proto.Unmarshal(response.Returns[0], &markTutorial)
	if err != nil {
		return nil, fmt.Errorf("Failed to call MARK_TUTORIAL_COMPLETE: %s", err)
	}

	return &markTutorial, nil
}

func (c *Instance) SetAvatarRequest(skin, hair, shirt, pants, hat, shoes, avatar, eyes, backpack int) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.SetAvatarMessage{
		PlayerAvatar: &protos.PlayerAvatar{
			Skin:     int32(skin),
			Hair:     int32(hair),
			Shirt:    int32(shirt),
			Pants:    int32(pants),
			Hat:      int32(hat),
			Shoes:    int32(shoes),
			Avatar:   int32(avatar),
			Eyes:     int32(eyes),
			Backpack: int32(backpack),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create SET_AVATAR: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_SET_AVATAR,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) SetAvatar(ctx context.Context, skin, hair, shirt, pants, hat, shoes, avatar, eyes, backpack int) (*protos.SetAvatarResponse, error) {
	request, err := c.SetAvatarRequest(skin, hair, shirt, pants, hat, shoes, avatar, eyes, backpack)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var setAvatar protos.SetAvatarResponse
	err = proto.Unmarshal(response.Returns[0], &setAvatar)
	if err != nil {
		return nil, fmt.Errorf("Failed to call SET_AVATAR: %s", err)
	}

	return &setAvatar, nil
}

func (c *Instance) GetDownloadURLsRequest(ids []string) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.GetDownloadUrlsMessage{
		AssetId: ids,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create GET_DOWNLOAD_URLS: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_GET_DOWNLOAD_URLS,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) GetDownloadURLs(ctx context.Context, ids []string) (*protos.GetDownloadUrlsResponse, error) {
	request, err := c.GetDownloadURLsRequest(ids)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var getDownloadURLs protos.GetDownloadUrlsResponse
	err = proto.Unmarshal(response.Returns[0], &getDownloadURLs)
	if err != nil {
		return nil, fmt.Errorf("Failed to call GET_DOWNLOAD_URLS: %s", err)
	}

	return &getDownloadURLs, nil
}

func (c *Instance) EncounterTutorialCompleteRequest(id int32) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.EncounterTutorialCompleteMessage{
		PokemonId: protos.PokemonId(id),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create ENCOUNTER_TUTORIAL_COMPLETE: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_ENCOUNTER_TUTORIAL_COMPLETE,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) EncounterTutorialComplete(ctx context.Context, id int32) (*protos.EncounterTutorialCompleteResponse, error) {
	request, err := c.EncounterTutorialCompleteRequest(id)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var encounter protos.EncounterTutorialCompleteResponse
	err = proto.Unmarshal(response.Returns[0], &encounter)
	if err != nil {
		return nil, fmt.Errorf("Failed to call ENCOUNTER_TUTORIAL_COMPLETE: %s", err)
	}

	return &encounter, nil
}

func (c *Instance) ClaimCodenameRequest(codename string) (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.ClaimCodenameMessage{
		Codename: codename,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create CLAIM_CODENAME: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_CLAIM_CODENAME,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) ClaimCodename(ctx context.Context, codename string) (*protos.ClaimCodenameResponse, error) {
	request, err := c.ClaimCodenameRequest(codename)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Returns) == 0 {
		return nil, errors.New("Server not accepted this request")
	}

	var claim protos.ClaimCodenameResponse
	err = proto.Unmarshal(response.Returns[0], &claim)
	if err != nil {
		return nil, fmt.Errorf("Failed to call CLAIM_CODENAME: %s", err)
	}

	return &claim, nil
}
