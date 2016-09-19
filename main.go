package main

import (
	"encoding/json"
	"fmt"
	"lis/command"
	"lis/location"
	"lis/point"
	"lis/tool"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	retOK        int = 0
	retArgsEmpty int = 101
	retArgsError int = 102

	retPointEmpty int = 151
	retError      int = 500
)

var responseMap map[int]string = make(map[int]string)

var p = fmt.Println

func responseMapDefine() {

	responseMap[retOK] = "ok"
	responseMap[retArgsEmpty] = "arguments empty"
	responseMap[retArgsError] = "arguments error:%s"

	responseMap[retPointEmpty] = "no point got"

	responseMap[retError] = "error"

}

//middle  where?   116.276329,40.056109

//pointTopLeft 上庄水库 116.210036,40.111421
//bottomright 清华大学  116.332556,40.009424

var pointTopLeft point.Point = point.Point{Lat: 40.111421, Lng: 116.210036}
var poinBottomRight point.Point = point.Point{Lat: 40.009424, Lng: 116.32556}

//var pointMiddle400m point.Point = point.Point{Lat: 40.057686, Lng: 116.291741}

var pointMiddle point.Point = point.Point{Lat: 40.056109, Lng: 116.276329}

var pointNum int = 288
var queryLimit int = 20

/***

dev route:
=================
point/set
point/query
location/query

point/delete

done precision = 5/6 的支持！





@todo http service
@snapshot
@snapshot recovery

*point/expire ?!
@done expire found when use,
@todo command.loopThoughtExpireCheck()


@master slave

*/

func main() {
	responseMapDefine()

	command.PointExpireCollect()

	//location.SetGeohashPrecision(5)
	testSet()

	//testPointQuery()
	testLocationQuery()

	//point garbage collect

	/* start to test delete**/
	testDelete()
	//time.Sleep(time.Second * 2)

	testLocationQuery()
	/* end test delete**/
	time.Sleep(time.Second * 1)

	testSummerize()
	//116.291741,40.057686

	http.HandleFunc("/point/set", pointSetHandler)
	http.HandleFunc("/point/delete", pointDeleteHandler)
	http.HandleFunc("/point/get", pointGetHandler)

	http.HandleFunc("/location/query", locationQueryHandler)

	http.ListenAndServe(":8000", nil)

	//fmt.Println("------test distance:----", tool.EarthDistance(pointMiddle.Lat, pointMiddle.Lng, pointMiddle400m.Lat, pointMiddle400m.Lng))
}

func pointSetHandler(response http.ResponseWriter, request *http.Request) {
	//queryForm, err := url.ParseQuery(request.URL.RawQuery)

}

func pointDeleteHandler(response http.ResponseWriter, request *http.Request) {
	//queryForm, err := url.ParseQuery(request.URL.RawQuery)

}

func locationQueryHandler(response http.ResponseWriter, request *http.Request) {
	//queryForm, err := url.ParseQuery(request.URL.RawQuery)

}

func pointRequestCommonArgsCheck(response http.ResponseWriter, request *http.Request) (point.QueryObject, bool) {
	args, _ := url.ParseQuery(request.URL.RawQuery)
	p("---request args---:", len(args), args)

	pqr := point.QueryObject{}
	if len(args) == 0 {
		response.Write(renderResponse(retArgsEmpty, "", nil))
		return pqr, false
	}
	oid, _ := args["id"]
	orole, _ := args["role"]

	if len(oid) < 1 || len(orole) < 1 {
		response.Write(renderResponse(retArgsEmpty, "", nil))

		return pqr, false
	}
	id, _ := strconv.Atoi(oid[0])
	role, _ := strconv.Atoi(orole[0])
	idu64 := uint64(id)
	roleint := int(role)
	if idu64 < uint64(1) || roleint < 1 {
		response.Write(renderResponse(retArgsError, "id/role should greater than 0", nil))
		return pqr, false
	}

	pqr.Id = idu64
	pqr.Role = roleint
	return pqr, true

}
func pointGetHandler(response http.ResponseWriter, request *http.Request) {

	pointQuery, check := pointRequestCommonArgsCheck(response, request)
	if !check {
		return
	}

	dt := command.PointQuery(pointQuery)
	p(pointQuery, dt)

	response.Write(renderResponse(retOK, "", formatPointForResponse(dt)))

}

func renderResponse(errno int, errmsg string, dt interface{}) []byte {

	ret := make(map[string]interface{})
	ret["errno"] = errno
	ret["errmsg"] = errmsg
	if errmsg == "" {
		msg, ok := responseMap[errno]
		if ok {
			ret["errmsg"] = msg
		}

	}

	ret["data"] = dt

	if dt == nil {
		ret["data"] = "{}"
	}

	bRet, _ := json.Marshal(ret)
	p(bRet)
	return bRet
}
func testInactiveRecycle() {
	for i := 160; i <= 169; i++ {
		qr := point.QueryObject{Id: uint64(i), Role: 5}
		command.PointDelete(qr)
	}
}

func formatPointForResponse(pt point.Point) map[string]interface{} {
	var ret map[string]interface{}
	ret = make(map[string]interface{})
	ret["id"] = pt.Id
	ret["role"] = pt.Role

	ret["lat"] = pt.Lat
	ret["lng"] = pt.Lng
	ret["update"] = pt.Update

	return ret
}

func testDelete() {
	for i := 160; i <= 169; i++ {
		qr := point.QueryObject{Id: uint64(i), Role: 5}
		command.PointDelete(qr)
	}

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
		for _, v := range ret {
			p(v.Pshell.Point.Id, ",")

		}
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

	role := 5
	for i := 0; i < pointNum; i++ {
		exp := int(i)
		//test expire
		if i >= 160 && i <= 169 {
			exp = 1

		}

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
