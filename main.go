package main

import (
	"lis/command"
	"lis/point"
)

func main() {

	pt := point.Point{Id: 1, Lat: 40.0728, Lng: 116.318014, Role: 5, Ext: 222, Expire: 333333}
	command.PointSet(pt)

	pt = point.Point{Id: 2, Lat: 40.0728, Lng: 116.318014, Role: 5, Ext: 222, Expire: 333333}
	command.PointSet(pt)

	pt = point.Point{Id: 1, Lat: 41.0728, Lng: 116.318014, Role: 5, Ext: 222, Expire: 333333}
	command.PointSet(pt)

	pt = point.Point{Id: 3, Lat: 40.0728, Lng: 116.318014, Role: 9, Ext: 222, Expire: 333333}
	command.PointSet(pt)

}

/*
command.PointDelete()


query := location.QueryObject{}

ret := command.LocationQuery(query)


//query2 := point.QueryObject{}
command.PointQuery(query2)

*/
