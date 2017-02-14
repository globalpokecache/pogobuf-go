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
		return nil, fmt.Errorf("Failed to create PLAYER_UPDATE: %s", err)
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

func (c *Instance) GetInventoryRequest() (*protos.Request, error) {
	msg, err := proto.Marshal(&protos.GetInventoryMessage{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create GET_INVENTORY: %s", err)
	}

	return &protos.Request{
		RequestType:    protos.RequestType_GET_INVENTORY,
		RequestMessage: msg,
	}, nil
}

func (c *Instance) GetInventory(ctx context.Context) (*protos.GetInventoryResponse, error) {
	request, err := c.GetInventoryRequest()
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
