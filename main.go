package main

import (
	"fmt"
	"lis/command"
	"lis/location"
	"lis/point"
	"lis/tool"
	"time"
)

var p = fmt.Println

//middle  where?   116.276329,40.056109

//pointTopLeft 上庄水库 116.210036,40.111421
//bottomright 清华大学  116.332556,40.009424

var pointTopLeft point.Point = point.Point{Lat: 40.111421, Lng: 116.210036}
var poinBottomRight point.Point = point.Point{Lat: 40.009424, Lng: 116.32556}

//var pointMiddle400m point.Point = point.Point{Lat: 40.057686, Lng: 116.291741}

var pointMiddle point.Point = point.Point{Lat: 40.056109, Lng: 116.276329}

var pointNum int = 200
var queryLimit int = 20

func main() {

	testSet()

	testPointQuery()
	testLocationQuery()

	testSummerize()
	//116.291741,40.057686

	//fmt.Println("------test distance:----", tool.EarthDistance(pointMiddle.Lat, pointMiddle.Lng, pointMiddle400m.Lat, pointMiddle400m.Lng))
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
	//数字山谷，滴滴大厦 116.296769,40.04987
	qr := location.QueryObject{Lat: pointMiddle.Lat, Lng: pointMiddle.Lng, Radius: 4000, Role: 5, Limit: queryLimit, Order: "distance"}
	ret := command.LocationQuery(qr)

	if tool.Debug() {
		p("------location.Query query=>result -------", qr, ret)
	}
}

func testSummerize() {
	point.Summerize()

	location.Summerize()
}

func testPointQuery() {

	qr := point.QueryObject{Id: 2, Role: 5}
	ret := command.PointQuery(qr)

	if tool.Debug() {
		p("------point.Query query=>result -------", qr, ret)
	}
}

func testSet() {

	//	var pointTopLeft point.Point = point.Point{Lat: 40.111421, Lng: 116.210036}
	//	var poinBottomRight point.Point = point.Point{Lat: 40.009424, Lng: 116.332556}

	//	var pointMiddle point.Point = point.Point{Lat: 40.056109, Lng: 116.296329}

	latStep := (pointTopLeft.Lat - poinBottomRight.Lat) / float64(pointNum)
	lngStep := (poinBottomRight.Lng - pointTopLeft.Lng) / float64(pointNum)
	fmt.Println("-------steps---------", latStep, lngStep, pointNum)

	//	hashTopLeft, _ := geohash.Encode(pointTopLeft.Lat, pointTopLeft.Lng, 6)
	//	hashBottomRight, _ := geohash.Encode(poinBottomRight.Lat, poinBottomRight.Lng, 6)
	//	fmt.Print("---hashlimit---", hashTopLeft, "------", hashBottomRight)

	role := uint8(5)
	for i := 0; i < pointNum; i++ {
		exp := time.Now().Second() + int(i)
		pt := point.Point{Id: uint64(i), Lat: pointTopLeft.Lat - latStep*float64(i), Lng: pointTopLeft.Lng + lngStep*float64(i), Role: role, Ext: 0, Expire: exp}

		//fmt.Println("-------item---------", pt.Lat, pt.Lng)

		command.PointSet(pt)
	}

}

/*
command.PointDelete()


query := location.QueryObject{}

ret := command.LocationQuery(query)


//query2 := point.QueryObject{}
command.PointQuery(query2)

*/
