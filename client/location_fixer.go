package client

import (
	"time"

	"github.com/globalpokecache/POGOProtos-go"
	"math"
)

func (c *Instance) locationFixer(quit chan struct{}) {
	lastpos := []float64{c.player.Latitude, c.player.Longitude}
	for {
		select {
		case <-quit:
			return
		default:
		}

		t := getTimestamp(time.Now())

		moving := (lastpos[0] != c.player.Latitude) && (lastpos[1] != c.player.Longitude)
		c.player.Lock()
		lastpos[0] = c.player.Latitude
		lastpos[1] = c.player.Longitude
		if c.lastLocationFix == nil || moving || randFloat() > 0.85 {
			c.player.Accuracy = []float64{5, 5, 5, 5, 10, 10, 10, 30, 30, 50, 65, math.Floor(randFloat()*(80-66)) + 66}[randInt(12)]

			junk := (randFloat() < 0.03)
			fix := &protos.Signature_LocationFix{
				Provider:       "fused",
				Latitude:       360.0,
				Longitude:      360.0,
				Altitude:       0.0,
				ProviderStatus: 3,
				LocationType:   1,
				Floor:          0,
				Course:         -1,
				Speed:          -1,
			}

			if !junk {
				fix.Latitude = float32(c.player.Latitude)
				fix.Longitude = float32(c.player.Longitude)
				fix.Altitude = float32(randTriang(300, 400, 350))
			}

			if randFloat() < 0.95 {
				fix.Course = float32(randTriang(0.0, 359.0, float64(c.lastLocationCourse)))
				fix.Speed = float32(randTriang(0.2, 4.25, 1))
				c.lastLocationCourse = fix.Course
			}

			if c.player.Accuracy >= 65 {
				fix.VerticalAccuracy = float32(randTriang(35, 100, 65))
				fix.HorizontalAccuracy = float32([]float64{c.player.Accuracy, 65, 65, 66 + (randFloat() * 14), 200}[randInt(5)])
			} else if c.player.Accuracy > 10 {
				fix.HorizontalAccuracy = float32(c.player.Accuracy)
				fix.VerticalAccuracy = float32([]float64{32, 48, 48, 64, 64, 96, 128}[randInt(7)])
			} else {
				fix.HorizontalAccuracy = float32(c.player.Accuracy)
				fix.VerticalAccuracy = float32([]float64{3, 4, 6, 6, 8, 12, 24}[randInt(7)])
			}

			fix.TimestampSnapshot = t - c.startedTime + uint64(-100+randInt(100))
			c.locationFixSync.Lock()
			c.locationFixes = append(c.locationFixes, fix)
			c.lastLocationFix = fix
			c.lastLocationFixTime = t
			c.locationFixSync.Unlock()
		}
		c.player.Unlock()

		time.Sleep(900 * time.Millisecond)
	}
}
