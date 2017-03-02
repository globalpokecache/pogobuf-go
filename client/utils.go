package client

import (
	"math"
	"math/rand"
	"time"

	"github.com/globalpokecache/POGOProtos-go"
	"github.com/globalpokecache/pogobuf-go/auth"
	"github.com/satori/go.uuid"
	"strings"
)

func randTriang(lower, upper, mode float64) float64 {
	var c = (mode - lower) / (upper - lower)
	var u = randFloat()

	if u <= c {
		return lower + math.Sqrt(u*(upper-lower)*(mode-lower))
	}

	return upper - math.Sqrt((1-u)*(upper-lower)*(upper-mode))
}
func randInt(l int) int {
	var s1 = rand.NewSource(time.Now().UnixNano())
	var r1 = rand.New(s1)
	return r1.Intn(l)
}

func randInt64(l int64) int64 {
	var s1 = rand.NewSource(time.Now().UnixNano())
	var r1 = rand.New(s1)
	return r1.Int63n(l)
}

func randFloat() float64 {
	var s1 = rand.NewSource(time.Now().UnixNano())
	var r1 = rand.New(s1)
	return r1.Float64()
}

func getTimestamp(t time.Time) uint64 {
	return uint64(t.UnixNano() / int64(time.Millisecond))
}

var (
	Devices = [][]string{
		{"iPhone5,1", "iPhone", "N41AP"},
		{"iPhone5,2", "iPhone", "N42AP"},
		{"iPhone5,3", "iPhone", "N48AP"},
		{"iPhone5,4", "iPhone", "N49AP"},
		{"iPhone6,1", "iPhone", "N51AP"},
		{"iPhone6,2", "iPhone", "N53AP"},
		{"iPhone7,1", "iPhone", "N56AP"},
		{"iPhone7,2", "iPhone", "N61AP"},
		{"iPhone8,1", "iPhone", "N71AP"},
		{"iPhone8,2", "iPhone", "N66AP"},
		{"iPhone8,4", "iPhone", "N69AP"},
		{"iPhone9,1", "iPhone", "D10AP"},
		{"iPhone9,2", "iPhone", "D11AP"},
		{"iPhone9,3", "iPhone", "D101AP"},
		{"iPhone9,4", "iPhone", "D111AP"},
	}

	OsVersions = []string{
		// 8
		"8.1.1", "8.1.2", "8.1.3", "8.2", "8.3",
		"8.4", "8.4.1",

		// 9
		"9.0", "9.0.1", "9.0.2",
		"9.1", "9.2", "9.2.1", "9.3", "9.3.1",
		"9.3.2", "9.3.3", "9.3.4", "9.3.5",

		// 10
		"10.0", "10.0.1", "10.0.2", "10.0.3", "10.1", "10.1.1",
	}
)

func NewDevice(p auth.Provider) *protos.Signature_DeviceInfo {
	uuid := uuid.NewV5(uuid.Nil, p.GetUsername())
	deviceID := strings.Replace(uuid.String(), "-", "", -1)
	device := Devices[randInt(len(Devices))]
	firmware := OsVersions[randInt(len(OsVersions))]

	return &protos.Signature_DeviceInfo{
		DeviceId:             deviceID,
		DeviceBrand:          "Apple",
		DeviceModel:          "iPhone",
		FirmwareBrand:        "iPhone OS",
		HardwareManufacturer: "Apple",
		DeviceModelBoot:      device[0],
		HardwareModel:        device[2],
		FirmwareType:         firmware,
	}
}
