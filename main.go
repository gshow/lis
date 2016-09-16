package main

import (
	"github.com/gshow/obis/command"
	"github.com/gshow/obis/point"
)

func main() {
	pt := point.Point{Id: 1, Lat: 40.0728, Lng: 116.318014, Role: 1, Ext: 222}
	obj := point.PointSetObject{Point: pt, Expire: 33333}

	command.PointSet(&obj)
}

/*
command.PointDelete()


query := location.QueryObject{}

ret := command.LocationQuery(query)


//query2 := point.QueryObject{}
command.PointQuery(query2)

*/
