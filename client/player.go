package client

import (
	"sync"

	"context"
	"googlemaps.github.io/maps"
)

type Player struct {
	sync.RWMutex
	lat, lng, accu, alt float64
}

func (p *Player) SetLatitude(l float64) {
	p.Lock()
	p.lat = l
	p.Unlock()
}

func (p *Player) SetLongitude(l float64) {
	p.Lock()
	p.lng = l
	p.Unlock()
}

func (p *Player) SetAccuracy(l float64) {
	p.Lock()
	p.accu = l
	p.Unlock()
}

func (p *Player) SetAltitude(l float64) {
	p.Lock()
	p.alt = l
	p.Unlock()
}

func (p *Player) Latitude() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.lat
}

func (p *Player) Longitude() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.lng
}

func (p *Player) Accuracy() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.accu
}

func (p *Player) Altitude() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.alt
}

func (c *Instance) SetPosition(ctx context.Context, lat, lon, accu, alt float64) error {
	moved := false
	if lat != c.player.Latitude() || lon != c.player.Longitude() {
		moved = true
	}
	if alt > 0 {
		c.player.SetAltitude(alt)
	} else if c.options.GoogleMapsKey != "" && moved {
		cli, err := maps.NewClient(maps.WithAPIKey(c.options.GoogleMapsKey))
		if err == nil {
			r := &maps.ElevationRequest{
				Locations: []maps.LatLng{
					{lat, lon},
				},
			}
			resp, err := cli.Elevation(ctx, r)
			if err == nil && len(resp) > 0 {
				c.player.SetAltitude(resp[0].Elevation)
			}
		}
	}

	c.player.SetLatitude(lat)
	c.player.SetLongitude(lon)
	if accu > 0 {
		c.player.SetAccuracy(accu)
	}

	return nil
}
