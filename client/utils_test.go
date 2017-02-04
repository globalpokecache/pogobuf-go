package client

import (
	"testing"
)

func TestGetRandomDevice(t *testing.T) {
	deviceInfo := GetRandomDevice()
	if len(deviceInfo.DeviceId) != 32 {
		t.Fatalf("Generated device id should have 32 characters, got %d", len(deviceInfo.DeviceId))
	}
}
