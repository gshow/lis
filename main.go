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
	"runtime"
	"strconv"
	"time"
)

const (
	retOK        int = 0
	retFailed    int = 1
	retArgsEmpty int = 101
	retArgsError int = 102

	retPointEmpty        int = 151
	retPointDeleteFailed int = 161
	retError             int = 500
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

@done
point/set
point/query
location/query

point/delete

precision = 5/6 的支持！




//性能测试

@todo http service
@snapshot
@snapshot recovery

*point/expire ?!
@done expire found when use,
@todo command.loopThoughtExpireCheck()


map concurrency thread safe resolve

@master slave



//config
geo precision,
http listen port,




*/

func main() {

	/*

				http://localhost:8000/point/get?id=188&role=5

				curl "http://localhost:8000/point/set" -d"id=188&role=5&lat=40.045839625&lng=116.2864475&ext=111111112&expire=2"

				curl "http://localhost:8000/point/delete" -d"id=199&role=5"

				http://localhost:8000/location/query?lat=40.056109&lng=116.276329&role=5&limit=10


		//ab

		ab -n1000 -c80 -p 'pointset.txt' -T 'application/x-www-form-urlencoded' 'http://t.hit.red:9001/point/set'





	*/
	runtime.GOMAXPROCS(runtime.NumCPU())
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

	error := http.ListenAndServe(":9001", nil)
	//fmt.Println("------test distance:----", tool.EarthDistance(pointMiddle.Lat, pointMiddle.Lng, pointMiddle400m.Lat, pointMiddle400m.Lng))
	fmt.Println(111111111, error)
}

func pointSetHandler(response http.ResponseWriter, request *http.Request) {
	//queryForm, err := url.ParseQuery(request.URL.RawQuery)

	method := request.Method
	if method != "POST" {
		methoderror := "request method eror"
		response.Write([]byte(methoderror))
		return
	}
	oid := request.PostFormValue("id")
	orole := request.PostFormValue("role")
	if len(oid) < 1 {
		response.Write(renderResponse(retArgsEmpty, "param id empty!", nil))
		return
	}
	if len(orole) < 1 {
		response.Write(renderResponse(retArgsEmpty, "param role empty!", nil))
		return

	}

	id, _ := strconv.Atoi(oid)
	role, _ := strconv.Atoi(orole)

	if id < 1 {
		response.Write(renderResponse(retArgsError, "id error!", nil))
		return
	}
	if role < 1 {
		response.Write(renderResponse(retArgsError, "role error!", nil))
		return
	}

	//lat/lng params check
	olat := request.PostFormValue("lat")
	olng := request.PostFormValue("lng")

	if len(olat) < 1 || len(olng) < 1 {
		response.Write(renderResponse(retArgsError, "request lat empty!", nil))
		return
	}
	if len(olng) < 1 || len(olng) < 1 {
		response.Write(renderResponse(retArgsError, "request lng empty!", nil))
		return
	}

	lat, _ := strconv.ParseFloat(olat, 64)
	lng, _ := strconv.ParseFloat(olng, 64)

	if lat < -180.0 || lat > 180.0 {
		response.Write(renderResponse(retArgsError, "lat range error", nil))
		return
	}
	if lng < -180.0 || lng > 180.0 {
		response.Write(renderResponse(retArgsError, "lng range error", nil))
		return
	}

	pt := point.Point{Id: uint64(id), Lat: lat, Lng: lng, Role: role}

	oext := request.PostFormValue("ext")
	if len(oext) > 0 {
		ext, _ := strconv.Atoi(oext)

		pt.Ext = int64(ext)
	}
	oexpire := request.PostFormValue("expire")
	if len(oexpire) > 0 {
		expire, _ := strconv.Atoi(oexpire)
		if expire > 86400*365 {
			response.Write(renderResponse(retArgsError, "expire time can be empty or within 86400*365", nil))
			return
		}

		pt.Expire = expire
	}

	//id,lat,lng,ext,expire
	//p("------set----point---", pt)
	set := command.PointSet(pt)
	if set {
		response.Write(renderResponse(retOK, "", nil))
	} else {
		response.Write(renderResponse(retFailed, "", nil))
	}

}

func pointDeleteHandler(response http.ResponseWriter, request *http.Request) {

	method := request.Method
	if method != "POST" {
		methoderror := "request method eror"
		response.Write([]byte(methoderror))
		return
	}
	oid := request.PostFormValue("id")
	orole := request.PostFormValue("role")
	if len(oid) < 1 {
		response.Write(renderResponse(retArgsEmpty, "param id empty!", nil))
		return
	}
	if len(orole) < 1 {
		response.Write(renderResponse(retArgsEmpty, "param role empty!", nil))
		return

	}

	id, _ := strconv.Atoi(oid)
	role, _ := strconv.Atoi(orole)

	if id < 1 {
		response.Write(renderResponse(retArgsError, "id error!", nil))
		return
	}
	if role < 1 {
		response.Write(renderResponse(retArgsError, "role error!", nil))
		return
	}

	pointQuery := point.QueryObject{Id: uint64(id), Role: role}

	command.PointDelete(pointQuery)
	//	ret := command.PointDelete(pointQuery)
	//	ret = true
	// 为了更友好，此处做了幂等 返回
	//	if !ret {
	//		response.Write(renderResponse(retPointDeleteFailed, "delete failed", nil))
	//		return
	//	}
	response.Write(renderResponse(retOK, "", nil))
}

func locationQueryHandler(response http.ResponseWriter, request *http.Request) {

	args, _ := url.ParseQuery(request.URL.RawQuery)
	//p("---request args---:", len(args), args)

	if len(args) == 0 {
		response.Write(renderResponse(retArgsEmpty, "", nil))
		return
	}

	//lat/lng params check
	olat, _ := args["lat"]
	olng, _ := args["lng"]

	if len(olat) < 1 || len(olng) < 1 {
		response.Write(renderResponse(retArgsError, "request lat or lng emtpy!", nil))

		return
	}

	lat, _ := strconv.ParseFloat(olat[0], 64)
	lng, _ := strconv.ParseFloat(olng[0], 64)

	if lat < -180.0 || lat > 180.0 {
		response.Write(renderResponse(retArgsError, "lat range error", nil))
		return
	}
	if lng < -180.0 || lng > 180.0 {
		response.Write(renderResponse(retArgsError, "lng range error", nil))
		return
	}

	//role check
	orole, _ := args["role"]
	role, _ := strconv.Atoi(orole[0])
	if role < 1 {
		response.Write(renderResponse(retArgsError, "role can not be 0", nil))
		return
	}

	//radius check
	oradius, radiusExist := args["radius"]
	var radius int
	if !radiusExist {
		radius = 2000
	} else {
		tradius, _ := strconv.Atoi(oradius[0])
		if tradius < 0 {
			response.Write(renderResponse(retArgsError, "radius must be bigger than 0", nil))
			return
		}
		radiusMax := location.GetRadiusMax()
		if radius > radiusMax {
			response.Write(renderResponse(retArgsError, "radius must be smaller than "+fmt.Sprintf("%.6f", radiusMax), nil))
			return
		}
		radius = tradius
	}

	//limit check
	olimit, limitExist := args["limit"]
	var limit int
	if !limitExist {
		limit = 20
	} else {
		tlimit, _ := strconv.Atoi(olimit[0])
		if tlimit < 1 {
			response.Write(renderResponse(retArgsError, "limit must be bigger than 0", nil))
			return
		}
		if tlimit > 1000 {
			response.Write(renderResponse(retArgsError, "limit must be smaller than 1000", nil))
			return
		}
		limit = tlimit

	}
	//order check
	oorder, orderExist := args["order"]
	var order string
	if !orderExist {
		order = "distance"
	} else {
		torder := oorder[0]
		if torder != "distance" && order != "update" {
			response.Write(renderResponse(retArgsError, "order type must be distance/update", nil))
			return
		}
		order = torder

	}

	qr := location.QueryObject{Lat: lat, Lng: lng, Radius: radius, Role: role, Limit: limit, Order: order}
	ret := command.LocationQuery(qr)

	response.Write(renderResponse(retOK, "", formatLocationQueryForResponse(qr, ret)))
	return

}

func formatLocationQueryForResponse(qr location.QueryObject, rs []location.QueryResult) map[string]interface{} {
	var ret map[string]interface{}
	ret = make(map[string]interface{})

	length := len(rs)
	if length < 1 {
		return ret
	}
	ret["count"] = length

	var collect []map[string]interface{}

	for _, po := range rs {
		item := make(map[string]interface{})
		pot := po.Pshell.Point
		item["id"] = pot.Id
		item["lat"] = pot.Lat
		item["lng"] = pot.Lng
		item["role"] = pot.Role
		item["ext"] = pot.Ext
		item["distance"] = strconv.FormatFloat(po.Distance, 'f', 3, 64) //fmt.Sprintf("%.3f", po.Distance)
		item["update"] = pot.Update

		collect = append(collect, item)

	}
	/**
	id,role,lat,lng,update,distance
	*/
	ret["role"] = qr.Role
	ret["requestLat"] = qr.Lat
	ret["requestLng"] = qr.Lng
	ret["requestRadius"] = qr.Radius
	ret["requestOrder"] = qr.Order
	ret["points"] = collect

	//p(ret)

	return ret
}

func pointRequestCommonArgsCheck(response http.ResponseWriter, request *http.Request) (point.QueryObject, bool) {
	args, _ := url.ParseQuery(request.URL.RawQuery)
	//p("---request args---:", len(args), args)

	pqr := point.QueryObject{}
	if len(args) == 0 {
		response.Write(renderResponse(retArgsEmpty, "no arguments!", nil))
		return pqr, false
	}
	oid, _ := args["id"]
	orole, _ := args["role"]

	if len(oid) < 1 || len(orole) < 1 {
		response.Write(renderResponse(retArgsEmpty, "id or role empty!", nil))

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
	//p(pointQuery, dt)
	if dt.Id < 1 {
		response.Write(renderResponse(retPointEmpty, "", nil))
		return

	}

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

	if dt != nil {
		ret["data"] = dt
	}

	bRet, _ := json.Marshal(ret)
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
	ret["ext"] = pt.Ext

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
