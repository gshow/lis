package command

import (
	"fmt"

	"obis/geohash"
	"obis/location"
	"obis/point"
)

func PointSet(pso point.Point, expire int32) bool {
	//save to roleMap-pointHashContainer-point
	gh, _ := geohash.Encode(pso.Lat, pso.Lng, location.GeohashPrecision)

	pso.Hash = gh

	point.Set(pso, expire)

	//save to geohash

	fmt.Println(gh, pso)
	return true
}

func PointDelete(point.PointQueryObject) bool {

	return true
}

func PointQuery(point.PointQueryObject) *point.Point {
	return new(point.Point)
}
