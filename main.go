package main

import (
	"lis/command"
	"lis/point"
)

func main() {

	pt := point.Point{Id: 1, Lat: 40.0728, Lng: 116.318014, Role: 1, Ext: 222}
	command.PointSet(pt, 333)
}

/*
command.PointDelete()


query := location.QueryObject{}

ret := command.LocationQuery(query)


//query2 := point.QueryObject{}
command.PointQuery(query2)

*/
