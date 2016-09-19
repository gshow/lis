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



RoleContainer{
	RoleMap:{
		roleid:ShellContainer{
			ShellMap:{
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
	Role int

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
	Role int
	Hash string
}

type Role struct {
	Lock     sync.RWMutex
	ShellMap map[uint64]*PointShell
}

//var PointsCollector = []*Point

type RoleContainer struct {
	Lock    sync.RWMutex
	RoleMap map[int]Role
}

var roleMap = RoleContainer{RoleMap: make(map[int]Role)}

func Query(qr QueryObject) Point {
	pt := Point{Id: qr.Id, Role: qr.Role}
	if !checkRoleContainer(pt, false) || !checkRole(pt, false) {
		return Point{}
	}

	shell, ok := roleMap.RoleMap[pt.Role].ShellMap[pt.Id]
	if ok == false {
		return Point{}
	}
	return shell.Point

}

func SetPrepare(pt Point) (string, *PointShell, func(bool)) {
	//save to roleMap-pointHashContainer-point

	checkRoleContainer(pt, true)
	checkRole(pt, true)

	oldHash := ""

	shell, shellExist := roleMap.RoleMap[pt.Role].ShellMap[pt.Id]

	if shellExist {
		oldHash = shell.Point.Hash

	} else {

		roleLock := roleMap.RoleMap[pt.Role].Lock
		roleLock.Lock()
		shell, shellExist = roleMap.RoleMap[pt.Role].ShellMap[pt.Id]
		if !shellExist {
			shell = createPointShell(pt)
			roleMap.RoleMap[pt.Role].ShellMap[pt.Id] = shell
		}
		roleLock.Unlock()
	}
	lock := shell.Lock
	lock.Lock()

	//	if tool.Debug() {
	//		fmt.Println("-----point.Set()----", pt, roleMap)
	//		fmt.Println("geohash:", pt.Hash)
	//	}
	return oldHash, shell, func(ret bool) {
		defer lock.Unlock()

		if !ret {
			return
		}
		shell.Point = pt

	}

}

func createPointShell(pt Point) *PointShell {
	shell := new(PointShell)
	shell.Point = pt
	return shell
}

func CheckNotExpire(pshell *PointShell) bool {

	if pshell.Point.Expire > 0 && time.Now().Second() > pshell.Point.Expire {
		go ExpireQueueAdd(pshell)
		return false
	}
	return true

}

func checkRole(pt Point, create bool) bool {
	_, ok := roleMap.RoleMap[pt.Role]
	if ok == false && create == false {
		return false
	}
	if ok == false && create == true {
		roleMap.Lock.Lock()
		defer roleMap.Lock.Unlock()

		_, ok := roleMap.RoleMap[pt.Role]
		if ok == false {
			pmap := Role{ShellMap: make(map[uint64]*PointShell)}
			roleMap.RoleMap[pt.Role] = pmap
		}
	}
	return true
}

func checkRoleContainer(pt Point, create bool) bool {
	_, ok := roleMap.RoleMap[pt.Role]
	if ok == false && create == false {
		return false
	}
	if ok == false && create == true {
		roleMap.Lock.Lock()
		defer roleMap.Lock.Unlock()

		_, ok := roleMap.RoleMap[pt.Role]
		if ok == false {
			pmap := Role{ShellMap: make(map[uint64]*PointShell)}
			roleMap.RoleMap[pt.Role] = pmap
		}
	}
	return true
}

func Summerize() {
	totalPoint := 0
	fmt.Println("-----roleMap.size----", len(roleMap.RoleMap))
	for role, son := range roleMap.RoleMap {
		totalPoint += len(son.ShellMap)
		fmt.Println("-----roleMap.role=>size----", role, len(son.ShellMap))

	}
	fmt.Println("-----total point size----", totalPoint)

}

func DeletePrepare(qr QueryObject) (bool, Point, func(bool)) {

	pt := Point{Id: qr.Id, Role: qr.Role}
	if !checkRoleContainer(pt, false) || !checkRole(pt, false) {
		return false, Point{}, func(bool) {}
	}

	shell, ok := roleMap.RoleMap[pt.Role].ShellMap[pt.Id]
	if ok == false {
		return false, Point{}, func(bool) {}
	}

	pret := shell.Point
	shellCon := roleMap.RoleMap[pt.Role]
	shellCon.Lock.Lock()

	return true, pret, func(result bool) {
		defer shellCon.Lock.Unlock()
		if result {
			delete(roleMap.RoleMap[pt.Role].ShellMap, pt.Id)
		}

	}

}
