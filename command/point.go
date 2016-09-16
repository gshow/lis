package command

import (
	"fmt"

	"github.com/gshow/obis/location"
	"github.com/gshow/obis/point"
)

func PointSet(pso *point.PointSetObject) bool {
	//save to roleMap-pointHashContainer-point
	//save to geohash
	a := &location.ContainerMapAll
	b := &point.Point{}
	fmt.Println(a, b, pso)
	return true
}

func PointDelete() bool {

	return true
}

func PointQuery() *point.Point {
	return new(point.Point)
}
