package main

import (
	"context"
	"log"

	"github.com/globalpokecache/pogobuf-go/auth"
	"github.com/globalpokecache/pogobuf-go/client"
	"github.com/globalpokecache/pogobuf-go/hash"
	"github.com/globalpokecache/pogobuf-go/helpers"
)

const (
	apiVersion = 5500
)

func main() {
	ctx := context.Background()
	ptc, err := auth.NewProvider("ptc")
	if err != nil {
		log.Fatal("Failed to create PTC provivder")
	}
	token, err := ptc.Login(ctx, "INSERT-YOUR-DUMMY-ACCOUNT", "INSERT-YOUR-DUMMY-ACCOUNT-PASSWORD")

	hp, err := hash.NewProvider("buddyauth", apiVersion)
	if err != nil {
		log.Fatal("Failed to create hash provivder")
	}
	hp.AddKey("INSERT-YOUR-BUDDY-AUTH-KEY")
	hp.SetDebug(true)

	client.Debug = true
	cli, err := client.New(&client.Options{
		AuthType:     "ptc",
		AuthToken:    token,
		HashProvider: hp,
		Version:      apiVersion,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	latitude := -1.234
	longitude := -1.234
	radius := 150.0

	cli.SetPosition(latitude, longitude, 0, 0)
	cli.Init(ctx)

	cells := helpers.GetCellsFromRadius(latitude, longitude, radius, 17)
	_, err = cli.GetMapObjects(ctx, cells, make([]int64, len(cells)))
	if err != nil {
		log.Fatalf("Failed to load map objects: %v\n", err)
	}
}
