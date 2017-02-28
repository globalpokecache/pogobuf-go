package client

import (
	"sync"
)

type Player struct {
	sync.Mutex
	Latitude, Longitude, Accuracy, Altitude float64
}

func (c *Instance) SetPosition(lat, lon, accu, alt float64) {
	c.player.Lock()
	c.player.Latitude = lat
	c.player.Longitude = lon
	if accu > 0 {
		c.player.Accuracy = accu
	}
	if alt > 0 {
		c.player.Altitude = alt
	}
	c.player.Unlock()
}
