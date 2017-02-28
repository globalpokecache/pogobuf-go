package client

import (
	"fmt"
	"github.com/globalpokecache/pogobuf-go/auth"
	"testing"
)

func TestGetRandomDevice(t *testing.T) {
	ptc, _ := auth.NewProvider("ptc", "mytest", "")
	deviceInfo := NewDevice(ptc)
	fmt.Println(deviceInfo)
}
