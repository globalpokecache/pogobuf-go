package helpers

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

const earthRadiusKm = 6371.01

func kmToAngle(km float64) s1.Angle {
	// The Earth's mean radius in kilometers (according to NASA).
	const earthRadiusKm = 6371.01
	return s1.Angle(km / earthRadiusKm)
}

func GetCellsFromRadius(lat, lng, radius float64, level int) []uint64 {

	// origin := uint64(s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng)).Parent(level))
	cap := s2.CapFromCenterAngle(s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng)), kmToAngle(radius/1000))

	rc := &s2.RegionCoverer{MaxLevel: level, MinLevel: level}
	r := s2.Region(cap)
	covering := rc.Covering(r)

	var cells []uint64
	for _, id := range covering {
		if cap.IntersectsCell(s2.CellFromCellID(id)) {
			cells = append(cells, uint64(id))
		}
	}

	return cells
}
