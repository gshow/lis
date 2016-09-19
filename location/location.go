package location

import (
	"fmt"
	"lis/geohash"
	"lis/point"
	"lis/tool"
	"sync"
)

/*

data structure:

LocationContainer{
	Ldata:{
		geohash:RoleContainer{
			RoleMap:{
				roleid:ShellContainer{
					ShellMap:{
						id:PointShell

					}

				}
			}


		}
	}

}

*/

type LocationContainer struct {
	Lock  sync.RWMutex
	Ldata map[string]point.RoleContainer
}

var geohashPrecision int = 6

type QueryObject struct {
	Lat    float64
	Lng    float64
	Radius float64
	Role   int

	Limit int
	Order string
}

type QueryResult struct {
	Pshell   *point.PointShell
	Distance float64
}

var radiusLoopMap5 map[int]float64 = map[int]float64{1: 14700.0, 2: 24500.0, 3: 34300.0, 4: 44100.0, 5: 53900.0}

var radiusLoopMap6 map[int]float64 = map[int]float64{1: 1828.0, 2: 3047.0, 3: 4265.0, 4: 5484.0, 5: 6703.0}

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

func getLoopTimesByRadius(radius float64) int {

	radiusLoopMap := radiusLoopMap6
	if GetGeohashPrecision() == 5 {
		radiusLoopMap = radiusLoopMap5
	}
	var ret, max int
	ret, max = 0, 0
	for times, distance := range radiusLoopMap {
		if radius <= distance {
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
		fmt.Println("-------got total for location query result:---", tmp, len(tmp))
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
	if false == checkHashContainer(pt, false) {
		return ret
	}
	if false == checkRoleContainer(pt, false) {

		return ret
	}

	ptNum := len(locationMap.Ldata[pt.Hash].RoleMap[pt.Role].ShellMap)
	if ptNum <= 0 {
		return ret
	}
	for _, pshell := range locationMap.Ldata[pt.Hash].RoleMap[pt.Role].ShellMap {
		if point.CheckNotExpire(pshell) == false {
			continue
		}
		distance := tool.EarthDistance(qr.Lat, qr.Lng, pshell.Point.Lat, pshell.Point.Lng)
		if distance > qr.Radius {
			continue
		}
		ret = append(ret, QueryResult{Pshell: pshell, Distance: distance})
	}

	if len(ret) > qr.Limit {
		ret = queryResultSort(ret, qr.Order)
		ret = ret[:qr.Limit]
	}

	return ret

}

func Summerize() {
	fmt.Println("-----locationMap.size----", len(locationMap.Ldata))
	for hash, roleContainer := range locationMap.Ldata {
		fmt.Println("-----locationMap.hash=>size----", hash, len(roleContainer.RoleMap))

		for roleid, shellContainer := range roleContainer.RoleMap {
			fmt.Println("-----locationMap.hash,role=>size----", hash, roleid, len(shellContainer.ShellMap))
		}
	}

}

var locationMap = LocationContainer{Ldata: make(map[string]point.RoleContainer)}

func Set(shell *point.PointShell, oldGeohash string, callback func(bool)) bool {

	//save to location hash index

	checkHashContainer(shell.Point, true)
	checkRoleContainer(shell.Point, true)

	if oldGeohash != "" && shell.Point.Hash != oldGeohash {

		delete(locationMap.Ldata[oldGeohash].RoleMap[shell.Point.Role].ShellMap, shell.Point.Id)

	}
	//	lock := locationMap.Ldata[shell.Point.Hash].RoleMap[shell.Point.Role].Lock
	//	lock.Lock()
	//	defer lock.Unlock()
	//此处，并发锁，会由 point.SetPrepare/point.DeletePrepare 控制，所以此处不使用锁了，不会出现这里的并发写问题

	locationMap.Ldata[shell.Point.Hash].RoleMap[shell.Point.Role].ShellMap[shell.Point.Id] = shell

	//	if tool.Debug() {
	//		fmt.Println("-----location.Set()----", locationMap)
	//	}

	callback(true)
	return true

}

func DeletePoint(pt point.Point, callback func(bool)) bool {

	//此处，并发锁，会由 point.SetPrepare/point.DeletePrepare 控制，所以此处不使用锁了，不会出现这里的并发写问题

	//map delete, do not need mutex
	delete(locationMap.Ldata[pt.Hash].RoleMap[pt.Role].ShellMap, pt.Id)

	callback(true)

	return true

}

func checkRoleContainer(pt point.Point, create bool) bool {
	_, ok := locationMap.Ldata[pt.Hash].RoleMap[pt.Role]

	if ok == false && create == false {
		return false
	}
	if ok == false {
		roleContainer := locationMap.Ldata[pt.Hash]
		roleContainer.Lock.Lock()
		defer roleContainer.Lock.Unlock()

		_, ok := roleContainer.RoleMap[pt.Role]
		if ok == false {
			shellContainer := point.Role{ShellMap: make(map[uint64]*point.PointShell)}
			roleContainer.RoleMap[pt.Role] = shellContainer
		}
	}
	return true
}

func checkHashContainer(pt point.Point, create bool) bool {
	_, ok := locationMap.Ldata[pt.Hash]
	if ok == false && create == false {
		return false
	}

	if ok == false {
		locationMap.Lock.Lock()
		defer locationMap.Lock.Unlock()

		_, ok := locationMap.Ldata[pt.Hash]
		if ok == false {
			roleContainer2 := point.RoleContainer{RoleMap: make(map[int]point.Role)}
			locationMap.Ldata[pt.Hash] = roleContainer2
		}
	}
	return true

}
