package main

import (
	"fmt"
	"lis/command"
	"lis/location"
	"lis/point"
)

var p = fmt.Println

func main() {
	testSet()

	testPointQuery()
	testLocationQuery()

	summerize()
}

func testLocationQuery() {
	/**
	  type QueryObject struct {
	  	Lat    float64
	  	Lng    float64
	  	Radius uint32
	  	Role   uint8

	  	Limit uint32
	  	Order string, enum(distance/update)
	  }

	*/
	qr := location.QueryObject{Lat: 40.072811113, Lng: 116.318014, Radius: 4000, Role: 5, Limit: 3, Order: "distance"}
	ret := location.Query(qr)

	p("------location.Query query=>result -------", qr, ret)
}

func summerize() {
	point.Summerize()

	location.Summerize()
}

func testPointQuery() {

	qr := point.QueryObject{Id: 2, Role: 5}
	ret := command.PointQuery(qr)

	p("------point.Query query=>result -------", qr, ret)
}

func testSet() {
	pt := point.Point{Id: 1, Lat: 40.072811113, Lng: 116.318014, Role: 5, Ext: 222, Expire: 333333}
	command.PointSet(pt)

	pt = point.Point{Id: 2, Lat: 40.0728223332, Lng: 116.318014, Role: 5, Ext: 222, Expire: 333333}
	command.PointSet(pt)

	pt = point.Point{Id: 1, Lat: 41.072823123, Lng: 116.318014, Role: 5, Ext: 222, Expire: 333333}
	command.PointSet(pt)

	pt = point.Point{Id: 1, Lat: 43.0728, Lng: 116.318014, Role: 5, Ext: 222, Expire: 333333}
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
