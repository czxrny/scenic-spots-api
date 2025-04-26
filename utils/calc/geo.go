package calc

import (
	"fmt"
	"math"
	"strconv"
)

type Coordinates struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

const EarthRadiusKm = 6371.0

func CoordinatesAfterRadius(latitude string, longitude string, radius string) (Coordinates, error) {
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return Coordinates{}, fmt.Errorf("invalid latitude parameter")
	}

	lon, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return Coordinates{}, fmt.Errorf("invalid longitude parameter")
	}

	radiusKm, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		return Coordinates{}, fmt.Errorf("invalid radius parameter")
	}

	latDistance := radiusKm / EarthRadiusKm * 180.0 / math.Pi
	lonDistance := latDistance / math.Cos(lat*math.Pi/180.0)

	minLat := lat - latDistance
	maxLat := lat + latDistance
	minLon := lon - lonDistance
	maxLon := lon + lonDistance

	if minLat < -90 {
		minLat = -90
	}
	if maxLat > 90 {
		maxLat = 90
	}

	if minLon < -180 {
		minLon = -180
	}
	if maxLon > 180 {
		maxLon = 180
	}

	return Coordinates{
		MinLat: minLat,
		MaxLat: maxLat,
		MinLon: minLon,
		MaxLon: maxLon,
	}, nil
}
