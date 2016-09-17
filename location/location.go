package location

import (
	"fmt"
	"lis/geohash"
	"lis/point"
	"sync"
)

/*

data structure:
locationMap{
	Ldata:{
		geohash:RoleContainer{
			Rdata:{
				roleid:ShellContainer{
					Sdata:{
						id:PointShell

					}

				}
			}


		}
	}

}

*/

const GeohashPrecision int = 6

type QueryObject struct {
	Lat    float64
	Lng    float64
	Radius uint32
	Role   uint8

	Limit uint32
	Order string
}

var RadisLoopMap map[int]uint32

func init() {
	RadisLoopMap = map[int]uint32{1: 1828, 2: 3047, 3: 4265, 4: 5484, 5: 6703}

}

func getLoopTimesByRadius(radius uint32) int {

	var ret, max int
	ret, max = 0, 0
	for times, distance := range RadisLoopMap {
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
	nei := geohash.LoopNeighbors(qr.Lat, qr.Lng, GeohashPrecision, loopTimes)
	var ret []string
	for _, circle := range nei {
		for _, hash := range circle {
			ret = append(ret, hash)
		}
	}
	fmt.Println("-------loop times:---", loopTimes)

	return ret

}

func Query(qr QueryObject) []point.Point {

	//qr := location.QueryObject{Lat: 40.072811113, Lng: 116.318014, Radius: 300, Role: 5, Limit: 3, Order: "distance/update"}

	hash, _ := geohash.Encode(qr.Lat, qr.Lng, GeohashPrecision)
	hash += ""

	neighbours := getNeighbours(qr)

	fmt.Println("-------loop neighbours:---", neighbours, len(neighbours))

	a := []point.Point{point.Point{}}
	return a
}

func Summerize() {
	fmt.Println("-----locationMap.size----", len(locationMap.Ldata))
	for hash, roleContainer := range locationMap.Ldata {
		fmt.Println("-----locationMap.hash=>size----", hash, len(roleContainer.Rdata))

		for roleid, shellContainer := range roleContainer.Rdata {
			fmt.Println("-----locationMap.hash,role=>size----", hash, roleid, len(shellContainer.Sdata))
		}
	}

}

type LocationContainer struct {
	Lock  sync.RWMutex
	Ldata map[string]point.RoleContainer
}

var locationMap = LocationContainer{Ldata: make(map[string]point.RoleContainer)}

func Set(shell *point.PointShell, oldGeohash string) bool {

	//save to location hash index

	checkHashContainer(shell)
	checkRoleContainer(shell)
	//checkShellContainer(shell)

	if oldGeohash != "" {
		fmt.Println(111, locationMap.Ldata[oldGeohash].Rdata[shell.Point.Role])
		shellContainer := locationMap.Ldata[oldGeohash].Rdata[shell.Point.Role]
		shellContainer.Lock.Lock()
		defer shellContainer.Lock.Unlock()

		_, ok := locationMap.Ldata[oldGeohash].Rdata[shell.Point.Role].Sdata[shell.Point.Id]

		if ok == true {
			delete(locationMap.Ldata[oldGeohash].Rdata[shell.Point.Role].Sdata, shell.Point.Id)
		}

	}
	locationMap.Ldata[shell.Point.Hash].Rdata[shell.Point.Role].Sdata[shell.Point.Id] = shell

	fmt.Println("-----location.Set()----", locationMap)
	return true

}

func checkRoleContainer(shell *point.PointShell) {
	_, ok := locationMap.Ldata[shell.Point.Hash].Rdata[shell.Point.Role]
	if ok == false {
		roleContainer := locationMap.Ldata[shell.Point.Hash]
		roleContainer.Lock.Lock()
		defer roleContainer.Lock.Unlock()

		_, ok := roleContainer.Rdata[shell.Point.Role]
		if ok == false {
			shellContainer := point.ShellContainer{Sdata: make(map[uint64]*point.PointShell)}
			roleContainer.Rdata[shell.Point.Role] = shellContainer
		}
	}
}

func checkHashContainer(shell *point.PointShell) {
	_, ok := locationMap.Ldata[shell.Point.Hash]

	if ok == false {
		locationMap.Lock.Lock()
		defer locationMap.Lock.Unlock()

		_, ok := locationMap.Ldata[shell.Point.Hash]
		if ok == false {
			roleContainer2 := point.RoleContainer{Rdata: make(map[uint8]point.ShellContainer)}
			locationMap.Ldata[shell.Point.Hash] = roleContainer2
		}
	}

}
