package command

import (
	"lis/geohash"
	"lis/location"
	"lis/point"
)

func PointSet(point2 point.Point) bool {
	//save to roleMap-pointHashContainer-point
	gh, _ := geohash.Encode(point2.Lat, point2.Lng, location.GeohashPrecision)
	point2.Hash = gh

	_, oldGeohash, shell := point.Set(point2)
	defer shell.Lock.Unlock()

	//save to geohash
	location.Set(shell, oldGeohash)

	return true
}

func PointDelete(point.QueryObject) bool {

	return true
}

func PointQuery(qr point.QueryObject) point.Point {
	return point.Query(qr)
}

func LocationQuery(qr location.QueryObject) []location.QueryResult {
	return location.Query(qr)
}
