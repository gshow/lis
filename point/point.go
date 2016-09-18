package point

import (
	"sync"
	//"lis/role"
	"fmt"
	//"lis/tool"
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
	Update int
	Expire int
}

type PointShell struct {
	Point Point
	Lock  sync.RWMutex
}

type QueryObject struct {
	Id   uint64
	Role uint8
	Hash string
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

	pt.Update = time.Now().Second()
	if pt.Expire > 0 {
		pt.Expire += time.Now().Second()
	}
	oldHash := ""

	shell, ok := roleMap.Rdata[pt.Role].Sdata[pt.Id]
	if ok == true {
		oldHash = shell.Point.Hash
		shell.Lock.Lock()

		roleMap.Rdata[pt.Role].Sdata[pt.Id].Point = pt
	} else {
		shell := createPointShell(pt)
		shell.Lock.Lock()
		roleMap.Rdata[pt.Role].Sdata[pt.Id] = shell

	}

	//	if tool.Debug() {
	//		fmt.Println("-----point.Set()----", pt, roleMap)
	//		fmt.Println("geohash:", pt.Hash)
	//	}
	return true, oldHash, roleMap.Rdata[pt.Role].Sdata[pt.Id]

}

func createPointShell(pt Point) *PointShell {
	shell := new(PointShell)
	shell.Point = pt
	return shell
}

func CheckNotExpire(pshell *PointShell) bool {

	if pshell.Point.Expire > 0 && time.Now().Second() > pshell.Point.Expire {
		go ExpireQueue.Add(pshell)
		return false
	}
	return true

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
	totalPoint := 0
	fmt.Println("-----roleMap.size----", len(roleMap.Rdata))
	for role, son := range roleMap.Rdata {
		totalPoint += len(son.Sdata)
		fmt.Println("-----roleMap.role=>size----", role, len(son.Sdata))

	}
	fmt.Println("-----total point size----", totalPoint)

}

func Delete(qr QueryObject) (Point, bool) {

	ret := true
	pret := Point{}
	pt := Point{Id: qr.Id, Role: qr.Role}
	if !checkRoleContainer(pt, false) || !checkShellContainer(pt, false) {
		return pret, ret
	}

	shell, ok := roleMap.Rdata[pt.Role].Sdata[pt.Id]
	if ok == false {
		return pret, ret
	}

	pret = shell.Point
	//map object do not need mutex
	delete(roleMap.Rdata[pt.Role].Sdata, pt.Id)

	return pret, ret

}
