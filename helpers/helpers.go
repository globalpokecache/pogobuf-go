package helpers

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"sort"
)

const earthRadiusKm = 6371.01

func kmToAngle(km float64) s1.Angle {
	// The Earth's mean radius in kilometers (according to NASA).
	const earthRadiusKm = 6371.01
	return s1.Angle(km / earthRadiusKm)
}

type AscID []uint64

func (a AscID) Len() int           { return len(a) }
func (a AscID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AscID) Less(i, j int) bool { return a[i] < a[j] }

func GetCellsFromRadius(lat, lng, radius float64, level int) []uint64 {
	if radius > 1500 {
		radius = 1500
	}

	cap := s2.CapFromCenterAngle(s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng)), kmToAngle(radius/1000))

	rc := &s2.RegionCoverer{MaxLevel: level, MinLevel: level}
	r := s2.Region(cap)
	covering := rc.Covering(r)

	var cells []uint64
	for _, id := range covering {
		cells = append(cells, uint64(id))
	}

	if len(cells) > 100 {
		cells = cells[:100]
	}

	sort.Sort(AscID(cells))

	return cells
}
