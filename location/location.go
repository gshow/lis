package location

import (
	"fmt"
	"lis/point"
	"sync"
)

var GeohashPrecision int = 6

type QueryObject struct {
	lat    float64
	lng    float64
	radius uint32
	limit  uint32
	role   uint8
	order  string
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
		//		fmt.Println(111, locationMap.Ldata[oldGeohash].Rdata[shell.Point.Role])
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
