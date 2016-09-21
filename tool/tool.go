package tool

import (
	"fmt"
	"math"
)

var P = fmt.Println
var debug = false

func Debug() bool {
	return debug == true
}

func SetDebug(dg bool) {
	debug = dg
}

func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6371000.0 // 6378137 meters
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))

	return dist * radius
}
