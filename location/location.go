package location

import (
	"fmt"
	"lis/geohash"
	"lis/point"
	"lis/tool"
	//"sync"
	smap "lis/safemap"
)

/*

data structure:

LocationContainer{
		geohash:RoleConainer{
				roleid:locationRole{
					ShellMap:{
						id:PointShell

					}

				}

		}

}

*/

var geohashPrecision int = 6

var locationMap = smap.New()

type QueryObject struct {
	Lat    float64
	Lng    float64
	Radius int
	Role   int

	Limit int
	Order string
}

type QueryResult struct {
	Pshell   *point.PointShell
	Distance float64
}

//var radiusLoopMap map[int]map[int]float64 = map[int]map[int]float64{5: {1: 14700.0, 2: 24500.0, 3: 34300.0, 4: 44100.0, 5: 53900.0}, 6: {1: 1828.0, 2: 3047.0, 3: 4265.0, 4: 5484.0, 5: 6703.0}}

var radiusLoopMap5 map[int]float64 = map[int]float64{1: 9800.0, 2: 19600.0, 3: 29400.0, 4: 39200.0, 5: 49000.0, 6: 58800.0, 7: 68600.0, 8: 78400.0, 9: 88200.0, 10: 98000.0, 11: 107800.0, 12: 117600.0}

var radiusLoopMap6 map[int]float64 = map[int]float64{1: 1219.2, 2: 2438.0, 3: 3656.8, 4: 4875.6, 5: 6094.4, 6: 7312.4, 7: 8526.0, 8: 9744.0, 9: 10943.0, 10: 12180.0, 11: 13392.0, 12: 14616.0}

var radiusLoopMap map[int]map[int]float64 = map[int]map[int]float64{5: radiusLoopMap5, 6: radiusLoopMap6}

func GetRadiusMax() int {
	return int(radiusLoopMap[GetGeohashPrecision()][5])
}

func SetGeohashPrecision(precision int) {
	if precision == 5 {
		geohashPrecision = 5
	} else {
		geohashPrecision = 6
	}
}

func GetGeohashPrecision() int {
	return geohashPrecision

}

func getLoopTimesByRadius(radius int) int {

	fradius := float64(radius)
	tmpMap := radiusLoopMap[GetGeohashPrecision()]
	var ret, max int
	ret, max = 0, 0
	for times, distance := range tmpMap {
		if fradius <= distance {
			ret = times
			break
		}
		max = times

	}
	if ret == 0 {
		ret = max
	}
	return ret

}

func getNeighbours(qr QueryObject) []string {
	loopTimes := getLoopTimesByRadius(qr.Radius)
	nei := geohash.LoopNeighbors(qr.Lat, qr.Lng, GetGeohashPrecision(), loopTimes)
	var ret []string
	for _, circle := range nei {
		for _, hash := range circle {
			ret = append(ret, hash)
		}
	}
	if tool.Debug() {
		fmt.Println("-------loop times:---", loopTimes)
	}

	return ret

}

func Query(qr QueryObject) []QueryResult {

	//qr := location.QueryObject{Lat: 40.072811113, Lng: 116.318014, Radius: 300, Role: 5, Limit: 3, Order: "distance/update"}

	hash, _ := geohash.Encode(qr.Lat, qr.Lng, GetGeohashPrecision())
	hash += ""

	neighbours := getNeighbours(qr)

	if tool.Debug() {
		fmt.Println("-------got neighbours:---", neighbours, len(neighbours))
	}
	if len(neighbours) <= 0 {
		return []QueryResult{}
	}
	tmp := []QueryResult{}
	for _, geohash := range neighbours {
		qret := queryHashArea(qr, geohash)
		if len(qret) > 0 {
			tmp = append(tmp, qret...)
		}
		if tool.Debug() {
			fmt.Println("-------got hash area query result:---", geohash, qret, len(qret))
		}

	}
	if len(tmp) <= 0 {
		return tmp
	}
	tmp = queryResultSort(tmp, qr.Order)
	if len(tmp) > qr.Limit {
		tmp = tmp[:qr.Limit]

	}

	if tool.Debug() {
		fmt.Println("-------got total for location query result:---", qr, tmp, len(tmp))
	}

	return tmp
}

func queryResultSort(rs []QueryResult, orderby string) []QueryResult {

	rsLen := len(rs)

	if orderby == "update" { //update time

		for i := 0; i < rsLen; i++ {
			for j := 0; j < rsLen-i-1; j++ {
				if rs[j].Pshell.Point.Update < rs[j+1].Pshell.Point.Update {
					rs[j], rs[j+1] = rs[j+1], rs[j]
				}
			}
		}

	} else { //distance

		for i := 0; i < rsLen; i++ {
			for j := 0; j < rsLen-i-1; j++ {
				if rs[j].Distance > rs[j+1].Distance {
					rs[j], rs[j+1] = rs[j+1], rs[j]
				}
			}
		}
	}
	return rs
}

func queryHashArea(qr QueryObject, geohash string) []QueryResult {
	pt := point.Point{Lat: qr.Lat, Lng: qr.Lng, Role: qr.Role, Hash: geohash}
	//qr := location.QueryObject{Lat: 40.072811113, Lng: 116.318014, Radius: 300, Role: 5, Limit: 3, Order: "distance/update"}

	ret := []QueryResult{}
	ishellCon, ok := checkPointShellContainer(pt, false)
	if !ok {
		return ret
	}

	ptNum := ishellCon.Size()
	if ptNum <= 0 {
		return ret
	}

	rangeCall := func(key interface{}, value interface{}) bool {
		pshell := value.(*point.PointShell)
		if point.CheckNotExpire(pshell) == false {
			return true
		}
		distance := tool.EarthDistance(qr.Lat, qr.Lng, pshell.Point.Lat, pshell.Point.Lng)
		if int(distance) > qr.Radius {
			return true
		}
		ret = append(ret, QueryResult{Pshell: pshell, Distance: distance})
		return true
	}

	ishellCon.Range(rangeCall)

	if len(ret) > qr.Limit {
		ret = queryResultSort(ret, qr.Order)
		ret = ret[:qr.Limit]
	}

	return ret

}

func Summerize() {

	rangeCallSon := func(k interface{}, value interface{}) bool {

		fmt.Println("-----locationMap,role=>size----", k, value.(*smap.SafeMap).Size())

		return true

	}
	rangeCall := func(k interface{}, value interface{}) bool {
		rolecon := value.(*smap.SafeMap)
		fmt.Println("-----locationMap.hash=>size----", k, rolecon.Size())

		rolecon.Range(rangeCallSon)
		return true

	}
	locationMap.Range(rangeCall)

}

func Set(shell *point.PointShell) bool {

	//save to location hash index

	ishellCon, _ := checkPointShellContainer(shell.Point, true)

	ishellCon.Set(shell.Point.Id, shell)

	//	if tool.Debug() {
	//		fmt.Println("-----location.Set()----", locationMap)
	//	}

	return true

}

func DeletePoint(pt point.Point) bool {

	ishellCon, exist := checkPointShellContainer(pt, false)
	if exist {
		return ishellCon.Delete(pt.Id)
	}

	return true

}

func checkPointShellContainer(pt point.Point, create bool) (*smap.SafeMap, bool) {

	//location, geohash,role,shell,point

	var ok bool
	_, ok = locationMap.Get(pt.Hash)

	if !ok {
		if !create {
			return smap.New(), false
		} else {
			locationMap.SetNotExist(pt.Hash, smap.New())
		}
	}
	roleCon := locationMap.PositiveMapGet(pt.Hash)
	//mod := pt.Id % idHashMod

	_, ok = roleCon.Get(pt.Role)
	if !ok {
		if !create {
			return smap.New(), false
		} else {
			roleCon.SetNotExist(pt.Role, smap.New())
		}
	}
	shellCon := roleCon.PositiveGet(pt.Role).(*smap.SafeMap)
	return shellCon, true
}
