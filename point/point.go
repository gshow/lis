package point

import (
	"sync"
	//"lis/role"
	"fmt"
	"time"
)

/*

data structure:
roleMap{
	Rdata:{
		roleid:ShellContainer{
			Sdata:{
				id:PointShell

			}

		}
	}

}

*/

type Point struct {
	Id   uint64
	Lat  float64
	Lng  float64
	Hash string
	Role uint8

	Ext    int64
	Update int64
	Expire int64
}

type PointShell struct {
	Point Point
	Lock  sync.RWMutex
}

type QueryObject struct {
	Id   uint64
	Role uint8
}

type ShellContainer struct {
	Lock  sync.RWMutex
	Sdata map[uint64]*PointShell
}

//var PointsCollector = []*Point

type RoleContainer struct {
	Lock  sync.RWMutex
	Rdata map[uint8]ShellContainer
}

var roleMap = RoleContainer{Rdata: make(map[uint8]ShellContainer)}

func Query(qr QueryObject) Point {
	pt := Point{Id: qr.Id, Role: qr.Role}
	if !checkRoleContainer(pt, false) || !checkShellContainer(pt, false) {
		return Point{}
	}

	shell, ok := roleMap.Rdata[pt.Role].Sdata[pt.Id]
	if ok == false {
		return Point{}
	}
	return shell.Point

}

func Set(pt Point) (bool, string, *PointShell) {
	//save to roleMap-pointHashContainer-point

	checkRoleContainer(pt, true)
	checkShellContainer(pt, true)

	pt.Update = time.Now().Unix()
	oldHash := ""

	shell, ok := roleMap.Rdata[pt.Role].Sdata[pt.Id]
	if ok == true {
		oldHash = shell.Point.Hash
		shell.Lock.Lock()

		roleMap.Rdata[pt.Role].Sdata[pt.Id].Point = pt
	} else {
		shell := &PointShell{Point: pt}
		shell.Lock.Lock()
		roleMap.Rdata[pt.Role].Sdata[pt.Id] = shell

	}

	fmt.Println("-----point.Set()----", pt, roleMap)
	return true, oldHash, roleMap.Rdata[pt.Role].Sdata[pt.Id]

}

func checkShellContainer(pt Point, create bool) bool {
	_, ok := roleMap.Rdata[pt.Role]
	if ok == false && create == false {
		return false
	}
	if ok == false && create == true {
		roleMap.Lock.Lock()
		defer roleMap.Lock.Unlock()

		_, ok := roleMap.Rdata[pt.Role]
		if ok == false {
			pmap := ShellContainer{Sdata: make(map[uint64]*PointShell)}
			roleMap.Rdata[pt.Role] = pmap
		}
	}
	return true
}

func checkRoleContainer(pt Point, create bool) bool {
	_, ok := roleMap.Rdata[pt.Role]
	if ok == false && create == false {
		return false
	}
	if ok == false && create == true {
		roleMap.Lock.Lock()
		defer roleMap.Lock.Unlock()

		_, ok := roleMap.Rdata[pt.Role]
		if ok == false {
			pmap := ShellContainer{Sdata: make(map[uint64]*PointShell)}
			roleMap.Rdata[pt.Role] = pmap
		}
	}
	return true
}

func Summerize() {
	fmt.Println("-----roleMap.size----", len(roleMap.Rdata))
	for role, son := range roleMap.Rdata {
		fmt.Println("-----roleMap.role=>size----", role, len(son.Sdata))

	}

}

func Delete(pt QueryObject) {

}
