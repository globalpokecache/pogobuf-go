package client

import (
	"fmt"
	"github.com/globalpokecache/pogobuf-go/auth"
	"github.com/globalpokecache/pogobuf-go/hash"
	"log"
	"testing"
)

func TestGetRequestId(t *testing.T) {
	fmt.Println(lMultiplier, lModulus, lMq, lMr)

	ptc, err := auth.NewProvider("ptc", "INSERT-YOUR-DUMMY-ACCOUNT", "INSERT-YOUR-DUMMY-ACCOUNT-PASSWORD")
	if err != nil {
		log.Fatal("Failed to create PTC provivder")
	}

	hp, err := hash.NewProvider("buddyauth", 5703)
	if err != nil {
		log.Fatal("Failed to create hash provivder")
	}

	c, _ := New(&Options{
		HashProvider: hp,
		AuthProvider: ptc,
	})
	c.rpcID = 0
	for i := 1; i <= 100; i++ {
		fmt.Println(i, c.getRequestId())
	}
}

func TestGetRandomDevice(t *testing.T) {
	ptc, _ := auth.NewProvider("ptc", "mytest", "")
	deviceInfo := NewDevice(ptc)
	fmt.Println(deviceInfo)
}
