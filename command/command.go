package command

import (
	"fmt"
	"lis/geohash"
	"lis/location"
	"lis/point"
	"lis/tool"
	"time"
)

func PointSet(point2 point.Point) bool {
	//save to roleMap-pointHashContainer-point
	gh, _ := geohash.Encode(point2.Lat, point2.Lng, location.GetGeohashPrecision())
	point2.Hash = gh

	point2.Update = int(time.Now().Unix())
	if point2.Expire > 0 {
		point2.Expire += int(time.Now().Unix())
	}
	if tool.Debug() {
		fmt.Println("----got set---", point2)
	}

	//	//save to geohash

	aggrement := func(id interface{}, value interface{}, oldvalue interface{}) bool {
		//rid := id.(int)
		if oldvalue != nil {
			oldshell := oldvalue.(*point.PointShell)
			//oldhash := oldshell.Point.Hash
			return location.DeletePoint(oldshell.Point)

		}
		shell := value.(*point.PointShell)
		return location.Set(shell)
	}

	ret := point.SetWithAggrement(point2, aggrement)

	return ret
}

func PointDelete(qr point.QueryObject) bool {
	expireCheck := false
	return _pointDelete(qr, expireCheck)
}

func _pointDelete(qr point.QueryObject, expireDelete bool) bool {

	aggrement := func(id interface{}, oldvalue interface{}) bool {
		//rid := id.(int)

		if oldvalue != nil {
			oldshell := oldvalue.(*point.PointShell)
			if expireDelete && oldshell.Point.Expire >= int(time.Now().Unix()) {

				return false
			}
			return location.DeletePoint(oldshell.Point)
		}
		return true

	}
	ok := point.DeleteWithAggrement(qr, aggrement)

	return ok

}

func PointQuery(qr point.QueryObject) point.Point {
	return point.Query(qr)
}

func LocationQuery(qr location.QueryObject) []location.QueryResult {
	return location.Query(qr)
}

func PointExpireCollect() {
	go pointQueueExpireCheck()
	go loopThoughtExpireCheck()

}

func loopThoughtExpireCheck() {
}
func pointQueueExpireCheck() {

	for true {
		//tpl := &point.PointShell{}
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println("error::PointQueueExpireCheck", err)
			}

		}()

		ret := point.ExpireQueueRead()
		if ret == nil {
			time.Sleep(time.Second * 1)
			continue
		}

		pshell := ret.Value.(*point.PointShell)

		expireCheck := true
		_pointDelete(point.QueryObject{pshell.Point.Id, pshell.Point.Role, pshell.Point.Hash}, expireCheck)

	}

}
