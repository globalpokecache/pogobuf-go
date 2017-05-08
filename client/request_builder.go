package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/globalpokecache/POGOProtos-go"
	"github.com/globalpokecache/pogobuf-go/pcrypt"
	"github.com/golang/protobuf/proto"
)

const defaultURL = "https://pgorelease.nianticlabs.com/plfe/rpc"
const downloadSettingsHash = "f43e9e403f233d4541feda9816e9d6085bccb087"

type SignatureInfo struct {
	DeviceInfo *protos.Signature_DeviceInfo
}

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

var (
	lMultiplier int64 = 16807
	lModulus    int64 = 0x7fffffff
	lMq               = lModulus / lMultiplier
	lMr               = lModulus % lMultiplier
)

func (c *Instance) getRequestId() int64 {
	var temp = lMultiplier*(c.lehmerSeed%lMq) - (lMr * (c.lehmerSeed / lMq))
	if temp > 0 {
		c.lehmerSeed = temp
	} else {
		c.lehmerSeed = temp + lModulus
	}
	c.rpcID++
	return (c.lehmerSeed << 32) | int64(c.rpcID)
}

var randAccuSeed = []int{5, 5, 5, 5, 10, 10, 10, 30, 30, 50, 65}

func (c *Instance) requestThrottle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		c.waitRequest <- struct{}{}
		time.Sleep(c.options.MinRequestInterval)
	}
}

func (c *Instance) call(ctx context.Context, requests []*protos.Request, prs []*protos.RequestEnvelope_PlatformRequest) (*protos.ResponseEnvelope, error) {
	// interval between requests
	<-c.waitRequest

	var respErr error
	var responseEnvelope *protos.ResponseEnvelope

	lat, long := c.player.Latitude(), c.player.Longitude()

	var randAccu = c.player.Accuracy()
	if randAccu == 0 {
		accuSeed := make([]int, len(randAccuSeed))
		copy(accuSeed, randAccuSeed)
		accuSeed = append(accuSeed, randInt(80-66)+66)
		randAccu = float64(accuSeed[randInt(len(accuSeed))])
	}

	requestEnvelope := &protos.RequestEnvelope{
		RequestId:  uint64(c.getRequestId()),
		StatusCode: int32(2),
	}

	if c.hasTicket {
		requestEnvelope.AuthTicket = c.authTicket
	} else {
		var unk2 int32
		if c.options.AuthProvider.Type() == "ptc" {
			unk2 = []int32{2, 8, 21, 21, 21, 28, 37, 56, 59, 59, 59}[randInt(11)]
		}

		requestEnvelope.AuthInfo = &protos.RequestEnvelope_AuthInfo{
			Provider: c.options.AuthProvider.Type(),
			Token: &protos.RequestEnvelope_AuthInfo_JWT{
				Contents: c.authToken,
				Unknown2: unk2,
			},
		}
	}

	var locFix = []*protos.Signature_LocationFix{}
	var lastLocFixTime uint64

	var ticket []byte
	var err error
	if c.hasTicket {
		ticket, err = proto.Marshal(requestEnvelope.AuthTicket)
		if err != nil {
			return nil, errors.New("Failed to marshal authTicket")
		}
	} else {
		ticket, err = proto.Marshal(requestEnvelope.AuthInfo)
		if err != nil {
			return nil, errors.New("Failed to marshal authTicket")
		}
	}

	var requestsBytes [][]byte
	if requests != nil && len(requests) > 0 {
		requestEnvelope.Requests = requests
		requestsBytes = make([][]byte, len(requests))
		for idx, request := range requests {
			debugProto(fmt.Sprintf("Request(%d)", idx), request)
			req, err := proto.Marshal(request)
			if err != nil {
				return nil, err
			}
			requestsBytes[idx] = req
		}
	}

	for done := false; !done; {
		select {
		case l := <-c.locationFixes:
			locFix = append(locFix, l)
			c.locationFixSync.Lock()
			c.lastLocationFix = l
			c.locationFixSync.Unlock()
			lastLocFixTime = l.TimestampSnapshot
		default:
			done = true
		}
	}

	if len(locFix) == 0 && c.lastLocationFix != nil {
		locFix = append(locFix, c.lastLocationFix)
		lastLocFixTime = c.lastLocationFix.TimestampSnapshot
	}

	for i := 0; i <= c.options.MaxTries; i++ {
		time.Sleep(time.Duration(i*300) * time.Millisecond)

		t := getTimestamp(time.Now())

		sinceStart := (t - c.startedTime)

		var sensorTS uint64
		if c.lastLocationFix != nil {
			requestEnvelope.MsSinceLastLocationfix = int64(t - (c.startedTime + lastLocFixTime))
			sensorTS = lastLocFixTime - uint64(-800+randInt(800))
		} else {
			requestEnvelope.MsSinceLastLocationfix = -1
			sensorTS = sinceStart - uint64(100+randInt(100))
		}

		requestEnvelope.Longitude = long
		requestEnvelope.Latitude = lat
		requestEnvelope.Accuracy = randAccu

		locHash1, locHash2, requestHash, err := c.options.HashProvider.Hash(
			ticket,
			c.sessionHash,
			lat,
			long,
			randAccu,
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
			LocationHash1:       int32(locHash1),
			LocationHash2:       int32(locHash2),
			SessionHash:         c.sessionHash,
			Timestamp:           t,
			TimestampSinceStart: sinceStart,
			Unknown25:           uk25,
			ActivityStatus: &protos.Signature_ActivityStatus{
				Stationary: true,
			},
		}

		signature.SensorInfo = []*protos.Signature_SensorInfo{
			{
				TimestampSnapshot:   sensorTS,
				LinearAccelerationX: randTriang(-1.7, 1.2, 0),
				LinearAccelerationY: randTriang(-1.4, 1.9, 0),
				LinearAccelerationZ: randTriang(-1.4, .9, 0),
				AttitudePitch:       randTriang(-1.5, 1.5, 0.4),
				AttitudeYaw:         randTriang(-3.1, 3.1, .198),
				AttitudeRoll:        randTriang(-2.8, 3.04, 0),
				RotationRateX:       randTriang(-4.7, 3.9, 0),
				RotationRateY:       randTriang(-4.7, 4.3, 0),
				RotationRateZ:       randTriang(-4.7, 6.5, 0),
				GravityX:            randTriang(-1, 1, 0),
				GravityY:            randTriang(-1, 1, -.2),
				GravityZ:            randTriang(-1, .7, -0.7),
				Status:              3,
			},
		}

		if len(signature.SensorInfo) > 0 {
			if c.rpcID == 2 {
				signature.SensorInfo[0].MagneticFieldAccuracy = -1
			} else {
				signature.SensorInfo[0].MagneticFieldX = randTriang(-54, 50, 0)
				signature.SensorInfo[0].MagneticFieldY = randTriang(-51, 57, -4.8)
				signature.SensorInfo[0].MagneticFieldZ = randTriang(-56, 43, -30)
				signature.SensorInfo[0].MagneticFieldAccuracy = []int32{-1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}[randInt(8)]
			}
		}

		if requests != nil && len(requests) > 0 {
			signature.RequestHash = requestHash
		}

		signature.LocationFix = locFix
		signature.DeviceInfo = c.options.SignatureInfo.DeviceInfo

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

		requestEnvelope.PlatformRequests = []*protos.RequestEnvelope_PlatformRequest{}

		if prs != nil {
			requestEnvelope.PlatformRequests = append(requestEnvelope.PlatformRequests, prs...)
		}

		if c.shouldAddPtr8(requests) {
			ptr8byte, err := proto.Marshal(&protos.UnknownPtr8Request{
				Message: c.ptr8,
			})
			if err == nil {
				requestEnvelope.PlatformRequests = append(requestEnvelope.PlatformRequests, &protos.RequestEnvelope_PlatformRequest{
					Type:           protos.PlatformRequestType_UNKNOWN_PTR_8,
					RequestMessage: ptr8byte,
				})
			}
		}

		requestEnvelope.PlatformRequests = append(requestEnvelope.PlatformRequests, &protos.RequestEnvelope_PlatformRequest{
			Type:           protos.PlatformRequestType_SEND_ENCRYPTED_SIGNATURE,
			RequestMessage: requestMessage,
		})

		responseEnvelope, respErr = c.rpc.Request(ctx, c.getServerURL(), requestEnvelope)
		if responseEnvelope == nil || respErr != nil {
			continue
		}

		for _, pr := range responseEnvelope.PlatformReturns {
			if pr.Type == protos.PlatformRequestType_UNKNOWN_PTR_8 {
				var ptr8 protos.UnknownPtr8Response
				err := proto.Unmarshal(pr.Response, &ptr8)
				if err == nil {
					if ptr8.Message != "" {
						c.ptr8 = ptr8.Message
					}
				}
			}
		}

		if responseEnvelope.ApiUrl != "" {
			c.setURL(responseEnvelope.ApiUrl)
		}

		if responseEnvelope.GetAuthTicket() != nil {
			c.setTicket(responseEnvelope.GetAuthTicket())
		}

		if responseEnvelope.StatusCode == protos.ResponseEnvelope_REDIRECT {
			continue
		}

		if responseEnvelope.StatusCode == protos.ResponseEnvelope_INVALID_AUTH_TOKEN {
			c.authTicket = nil
			c.login(ctx)
			continue
		}

		if respErr == nil {
			break
		}
	}

	return responseEnvelope, respErr
}

func (c *Instance) Call(ctx context.Context, requests ...*protos.Request) (*protos.ResponseEnvelope, error) {
	return c.call(ctx, requests, nil)
}

func (c *Instance) CallWithPlatformRequests(ctx context.Context, requests []*protos.Request, prs []*protos.RequestEnvelope_PlatformRequest) (*protos.ResponseEnvelope, error) {
	return c.call(ctx, requests, prs)
}

func (c *Instance) shouldAddPtr8(requests []*protos.Request) bool {
	r := randFloat()
	if len(requests) == 0 {
		return false
	}
	rtype := requests[0].RequestType
	if c.ptr8 != "" &&
		((rtype == protos.RequestType_GET_PLAYER || rtype == protos.RequestType_GET_MAP_OBJECTS) && r > 0.5) &&
		(rtype == protos.RequestType_ENCOUNTER && r > 0.9) ||
		r > 0.97 {
		return true
	}

	// hasMap := false
	// for _, req := range requests {
	// 	if req.RequestType == protos.RequestType_GET_MAP_OBJECTS {
	// 		return true
	// 	}
	// }

	// if hasMap {
	// 	if !c.firstGetMap {
	// 		return true
	// 	}
	// 	c.firstGetMap = false
	// }

	return false
}
