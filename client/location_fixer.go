package client

import (
	"context"
	"math"
	"time"

	"github.com/globalpokecache/POGOProtos-go"
)

func (c *Instance) locationFixer(ctx context.Context) {
	lastpos := []float64{c.player.Latitude(), c.player.Longitude()}
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		t := getTimestamp(time.Now())

		moving := (lastpos[0] != c.player.Latitude()) || (lastpos[1] != c.player.Longitude())
		lastpos[0] = c.player.Latitude()
		lastpos[1] = c.player.Longitude()
		if c.lastLocationFix == nil || moving || randFloat() > 0.85 {
			c.player.SetAccuracy([]float64{5, 5, 5, 5, 10, 10, 10, 30, 30, 50, 65, math.Floor(randFloat()*(80-66)) + 66}[randInt(12)])

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
				fix.Latitude = float32(c.player.Latitude())
				fix.Longitude = float32(c.player.Longitude())
				if c.player.Altitude() > 0 {
					fix.Altitude = float32(c.player.Altitude())
				} else {
					fix.Altitude = float32(randTriang(300, 400, 350))
				}
			}

			if randFloat() < 0.95 {
				fix.Course = float32(randTriang(0.0, 359.0, float64(c.lastLocationCourse)))
				fix.Speed = float32(randTriang(0.2, 4.25, 1))
				c.lastLocationCourse = fix.Course
			}

			if c.player.Accuracy() >= 65 {
				fix.VerticalAccuracy = float32(randTriang(35, 100, 65))
				fix.HorizontalAccuracy = float32([]float64{c.player.Accuracy(), 65, 65, 66 + (randFloat() * 14), 200}[randInt(5)])
			} else if c.player.Accuracy() > 10 {
				fix.HorizontalAccuracy = float32(c.player.Accuracy())
				fix.VerticalAccuracy = float32([]float64{32, 48, 48, 64, 64, 96, 128}[randInt(7)])
			} else {
				fix.HorizontalAccuracy = float32(c.player.Accuracy())
				fix.VerticalAccuracy = float32([]float64{3, 4, 6, 6, 8, 12, 24}[randInt(7)])
			}

			fix.TimestampSnapshot = t - c.startedTime + uint64(-100+randInt(100))

			for done := false; !done; {
				select {
				case c.locationFixes <- fix:
					done = true
				default:
					<-c.locationFixes
				}
			}
		}

		randSleep(900, 950)
	}
}
