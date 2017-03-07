package client

import (
	"sync"

	"context"
	"googlemaps.github.io/maps"
)

type Player struct {
	sync.Mutex
	Latitude, Longitude, Accuracy, Altitude float64
}

func (c *Instance) SetPosition(ctx context.Context, lat, lon, accu, alt float64) {
	c.player.Lock()
	defer c.player.Unlock()
	moved := false
	if lat != c.player.Latitude || lon != c.player.Longitude {
		moved = true
	}
	c.player.Latitude = lat
	c.player.Longitude = lon
	if accu > 0 {
		c.player.Accuracy = accu
	}
	if alt > 0 {
		c.player.Altitude = alt
	} else if c.options.GoogleMapsKey != "" && moved {
		cli, err := maps.NewClient(maps.WithAPIKey(c.options.GoogleMapsKey))
		if err != nil {
			return
		}
		r := &maps.ElevationRequest{
			Locations: []maps.LatLng{
				{lat, lon},
			},
		}
		resp, err := cli.Elevation(ctx, r)
		if err != nil {
			return
		}
		if len(resp) > 0 {
			c.player.Altitude = resp[0].Elevation
		}
	}
}
