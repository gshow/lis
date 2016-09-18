package command

import (
	"fmt"
	"lis/geohash"
	"lis/location"
	"lis/point"
	"time"
)

func PointSet(point2 point.Point) bool {
	//save to roleMap-pointHashContainer-point
	gh, _ := geohash.Encode(point2.Lat, point2.Lng, location.GetGeohashPrecision())
	point2.Hash = gh

	_, oldGeohash, shell := point.Set(point2)
	defer shell.Lock.Unlock()

	//save to geohash
	location.Set(shell, oldGeohash)

	return true
}

func PointDelete(qr point.QueryObject) bool {

	return _pointDelete(qr)
}

func _pointDelete(qr point.QueryObject) bool {
	pt, ok := point.Delete(qr)
	if ok == true && pt.Id > 0 {
		location.DeletePoint(pt)
	}
	return true
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

		ret := point.ExpireQueue.Read()
		if ret == nil {
			time.Sleep(time.Second * 1)
			continue
		}

		pshell := ret.Value.(*point.PointShell)

		_pointDelete(point.QueryObject{pshell.Point.Id, pshell.Point.Role, pshell.Point.Hash})

	}

}
