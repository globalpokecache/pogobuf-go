package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/globalpokecache/POGOProtos-go"
	"github.com/globalpokecache/pogobuf-go/pcrypt"
	"math"
	"time"
)

const defaultURL = "https://pgorelease.nianticlabs.com/plfe/rpc"
const downloadSettingsHash = "54b359c97e46900f87211ef6e6dd0b7f2a3ea1f5"

func (c *Instance) getServerURL() string {
	var url string
	if c.serverURL != "" {
		url = c.serverURL
	} else {
		url = defaultURL
	}
	return url
}

func (c *Instance) setTicket(ticket *protos.AuthTicket) {
	c.hasTicket = true
	c.authTicket = ticket
}

func (c *Instance) setURL(urlToken string) {
	c.serverURL = fmt.Sprintf("https://%s/rpc", urlToken)
}

func (c *Instance) getRequestId() uint64 {
	var r int64
	if c.rpcID == 0 {
		c.rpcID = 1
		if c.options.SignatureInfo.DeviceInfo != nil && c.options.SignatureInfo.DeviceInfo.DeviceBrand != "Apple" {
			r = 0x53B77E48
		} else {
			r = 0x000041A7
		}
	} else {
		r = int64(randInt(int(math.Pow(2.0, 31.0))))
	}
	c.rpcID++
	cnt := c.rpcID
	return uint64(((r | ((cnt & 0xFFFFFFFF) >> 31)) << 32) | cnt)
}

var randAccuSeed = []int{5, 5, 5, 5, 10, 10, 10, 30, 30, 50, 65}

func (c *Instance) Call(ctx context.Context, requests ...*protos.Request) (*protos.ResponseEnvelope, error) {
	var randAccu = c.player.Accuracy
	if randAccu == 0 {
		accuSeed := make([]int, len(randAccuSeed))
		copy(accuSeed, randAccuSeed)
		accuSeed = append(accuSeed, randInt(80-66)+66)
		randAccu = float64(accuSeed[randInt(len(accuSeed))])
	}

	requestEnvelope := &protos.RequestEnvelope{
		RequestId:  c.getRequestId(),
		StatusCode: int32(2),

		MsSinceLastLocationfix: int64(100 + randInt(900)),

		Longitude: c.player.Longitude,
		Latitude:  c.player.Latitude,
		Accuracy:  randAccu,

		Requests: requests,
	}

	if c.hasTicket {
		requestEnvelope.AuthTicket = c.authTicket
	} else {
		requestEnvelope.AuthInfo = &protos.RequestEnvelope_AuthInfo{
			Provider: c.options.AuthType,
			Token: &protos.RequestEnvelope_AuthInfo_JWT{
				Contents: c.options.AuthToken,
				Unknown2: int32(c.token2),
			},
		}
	}

	if c.hasTicket {
		t := getTimestamp(time.Now())

		ticket, err := proto.Marshal(c.authTicket)
		if err != nil {
			return nil, errors.New("Failed to marshal authTicket")
		}

		requestsBytes := make([][]byte, len(requests))
		for idx, request := range requests {
			req, err := proto.Marshal(request)
			if err != nil {
				return nil, err
			}
			requestsBytes[idx] = req
		}

		locHash1, locHash2, requestHash, err := c.options.HashProvider.Hash(
			ticket,
			c.sessionHash,
			requestEnvelope.Latitude,
			requestEnvelope.Longitude,
			requestEnvelope.Accuracy,
			t,
			requestsBytes,
		)
		if err != nil {
			return nil, fmt.Errorf("Hash provider failed to hash: %s", err)
		}

		uk25, err := getUnk25(c.options.Version)
		if err != nil {
			return nil, err
		}

		signature := &protos.Signature{
			RequestHash:         requestHash,
			LocationHash1:       int32(locHash1),
			LocationHash2:       int32(locHash2),
			SessionHash:         c.sessionHash,
			Timestamp:           t,
			TimestampSinceStart: (t - getTimestamp(c.startedTime)),
			Unknown25:           uk25,
			ActivityStatus: &protos.Signature_ActivityStatus{
				Stationary: true,
			},
			SensorInfo: []*protos.Signature_SensorInfo{
				{
					TimestampSnapshot:     getTimestamp(c.startedTime) + uint64(100+randInt(150)),
					LinearAccelerationX:   -0.7 + randFloat()*1.4,
					LinearAccelerationY:   -0.7 + randFloat()*1.4,
					LinearAccelerationZ:   -0.7 + randFloat()*1.4,
					RotationRateX:         0.7 * randFloat(),
					RotationRateY:         0.8 * randFloat(),
					RotationRateZ:         0.8 * randFloat(),
					AttitudePitch:         -1.0 + randFloat()*2.0,
					AttitudeRoll:          -1.0 + randFloat()*2.0,
					AttitudeYaw:           -1.0 + randFloat()*2.0,
					GravityX:              -1.0 + randFloat()*2.0,
					GravityY:              -1.0 + randFloat()*2.0,
					GravityZ:              -1.0 + randFloat()*2.0,
					MagneticFieldAccuracy: -1,
					Status:                3,
				},
			},
		}

		signature.LocationFix = buildLocationFixes(requestEnvelope, getTimestamp(c.startedTime))

		if signature.TimestampSinceStart < 5000 {
			signature.Timestamp = uint64(5000 + randInt(8000))
		}

		if c.options.SignatureInfo.DeviceInfo != nil {
			signature.DeviceInfo = c.options.SignatureInfo.DeviceInfo
		} else {
			signature.DeviceInfo = GetRandomDevice()
		}

		debugProto("Signature", signature)

		signatureProto, err := proto.Marshal(signature)
		if err != nil {
			return nil, errors.New("Failed to marshal the request signature")
		}

		requestMessage, err := proto.Marshal(&protos.SendEncryptedSignatureRequest{
			EncryptedSignature: pcrypt.Encrypt(signatureProto, uint32(signature.TimestampSinceStart)),
		})
		if err != nil {
			return nil, errors.New("Failed to marshal request message")
		}

		requestEnvelope.PlatformRequests = []*protos.RequestEnvelope_PlatformRequest{
			{
				Type:           protos.PlatformRequestType_SEND_ENCRYPTED_SIGNATURE,
				RequestMessage: requestMessage,
			},
		}
	}

	var responseEnvelope *protos.ResponseEnvelope
	var err error
	for i := 1; i <= c.options.MaxTries; i++ {
		responseEnvelope, err = c.rpc.Request(ctx, c.getServerURL(), requestEnvelope)

		if responseEnvelope.ApiUrl != "" {
			c.setURL(responseEnvelope.ApiUrl)
		}

		if responseEnvelope.GetAuthTicket() != nil {
			c.setTicket(responseEnvelope.GetAuthTicket())
		}

		if responseEnvelope.StatusCode == protos.ResponseEnvelope_REDIRECT {
			time.Sleep(time.Duration(i*300) * time.Millisecond)
			continue
		}

		if err == nil {
			break
		}

		time.Sleep(time.Duration(i*300) * time.Millisecond)
	}

	return responseEnvelope, err
}
