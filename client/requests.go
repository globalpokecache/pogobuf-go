package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/globalpokecache/POGOProtos-go"
)

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

func (c *Instance) DownloadSettings(ctx context.Context) (downloadSettings *protos.DownloadSettingsResponse, err error) {
	request, err := c.DownloadSettingsRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return
	}

	if len(response.Returns) == 0 {
		err = errors.New("Server not accepted this request")
		return
	}

	err = proto.Unmarshal(response.Returns[0], downloadSettings)
	if err != nil {
		err = fmt.Errorf("Failed to call DOWNLOAD_SETTINGS: %s", err)
		return
	}

	return
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

func (c *Instance) PlayerUpdate(ctx context.Context) (playerUpdate *protos.PlayerUpdateResponse, err error) {
	request, err := c.PlayerUpdateRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return
	}

	if len(response.Returns) == 0 {
		err = errors.New("Server not accepted this request")
		return
	}

	err = proto.Unmarshal(response.Returns[0], playerUpdate)
	if err != nil {
		err = fmt.Errorf("Failed to call PLAYER_UPDATE: %s", err)
		return
	}

	return
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

func (c *Instance) GetPlayer(ctx context.Context, country, language, timezone string) (getPlayer *protos.GetPlayerResponse, err error) {
	request, err := c.GetPlayerRequest(country, language, timezone)
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return
	}

	if len(response.Returns) == 0 {
		err = errors.New("Server not accepted this request")
		return
	}

	err = proto.Unmarshal(response.Returns[0], getPlayer)
	if err != nil {
		err = fmt.Errorf("Failed to call GET_PLAYER: %s", err)
		return
	}

	return
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

func (c *Instance) GetHatchedEggs(ctx context.Context) (getHatchedEggs *protos.GetHatchedEggsResponse, err error) {
	request, err := c.GetHatchedEggsRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return
	}

	if len(response.Returns) == 0 {
		err = errors.New("Server not accepted this request")
		return
	}

	err = proto.Unmarshal(response.Returns[0], getHatchedEggs)
	if err != nil {
		err = fmt.Errorf("Failed to call GET_HATCHED_EGGS: %s", err)
		return
	}

	return
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

func (c *Instance) GetInventory(ctx context.Context) (getInventory *protos.GetInventoryResponse, err error) {
	request, err := c.GetInventoryRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return
	}

	if len(response.Returns) == 0 {
		err = errors.New("Server not accepted this request")
		return
	}

	err = proto.Unmarshal(response.Returns[0], getInventory)
	if err != nil {
		err = fmt.Errorf("Failed to call GET_INVENTORY: %s", err)
		return
	}

	return
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

func (c *Instance) CheckAwardedBadges(ctx context.Context) (checkAwardedBadges *protos.CheckAwardedBadgesResponse, err error) {
	request, err := c.CheckAwardedBadgesRequest()
	if err != nil {
		return nil, err
	}

	var response *protos.ResponseEnvelope
	response, err = c.Call(ctx, request)
	if err != nil {
		return
	}

	if len(response.Returns) == 0 {
		err = errors.New("Server not accepted this request")
		return
	}

	err = proto.Unmarshal(response.Returns[0], checkAwardedBadges)
	if err != nil {
		err = fmt.Errorf("Failed to call CHECK_AWARDED_BADGES: %s", err)
		return
	}

	return
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
