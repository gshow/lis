package command

import (
	"fmt"

	"github.com/gshow/obis/geohash"
	"github.com/gshow/obis/location"
	"github.com/gshow/obis/point"
)




func PointSet(pso point.PointSetObject) bool {
	//save to roleMap-pointHashContainer-point
	
	//save to geohash
	
	
	
	
	
	
	a := &location.ContainerMapAll

	gh, _ := geohash.Encode(pso.Lat, pso.Lng, location.GeohashPrecision)

	_, ok := location.ContainerMapAll.Data[pso.Role]
	if(){
		
	}

	fmt.Println(a, pso, gh, ok)
	return true
}

func PointDelete(point.PointQueryObject) bool {

	return true
}

func PointQuery(point.PointQueryObject) *point.Point {
	return new(point.Point)
}
