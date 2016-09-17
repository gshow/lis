package command

import (
	"fmt"

	"lis/geohash"
	"lis/location"
	"lis/point"
)

func PointSet(point2 point.Point) bool {
	//save to roleMap-pointHashContainer-point
	gh, _ := geohash.Encode(point2.Lat, point2.Lng, location.GeohashPrecision)
	point2.Hash = gh

	ok, oldGeohash, shell := point.Set(point2)
	defer shell.Lock.Unlock()

	//save to geohash
	location.Set(shell, oldGeohash)

	fmt.Println("-----finishi.result()----", ok, oldGeohash, shell)
	return true
}

func PointDelete(point.PointQueryObject) bool {

	return true
}

func PointQuery(point.PointQueryObject) *point.Point {
	return new(point.Point)
}
