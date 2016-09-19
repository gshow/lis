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



roleContainer{
	RoleMap:{
		roleid:Role{
			IdHashMap:{
				idhash:ShellContainer{
					ShellMap:{
						id:PointShell
					}

				}
			}

		}
	}

}

*/

var idHashMod = uint64(102400)

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

type idHashContainer struct {
	Lock     sync.RWMutex
	ShellMap map[uint64]*PointShell
}
type roleContainer struct {
	Lock    sync.RWMutex
	RoleMap map[int]roleObject
}
type roleObject struct {
	Lock       sync.RWMutex
	IdHsashMap map[uint64]idHashContainer
}

//var PointsCollector = []*Point

var roleMap = roleContainer{RoleMap: make(map[int]roleObject)}

func Query(qr QueryObject) Point {
	pt := Point{Id: qr.Id, Role: qr.Role}
	if checkIdHashContainer(pt, false) {
		return Point{}
	}

	shell, ok := roleMap.RoleMap[pt.Role].IdHsashMap[qr.Id%idHashMod].ShellMap[qr.Id]
	if ok == false {
		return Point{}
	}
	return shell.Point

}

func SetPrepare(pt Point) (string, *PointShell, func(bool)) {
	//save to roleMap-pointHashContainer-point

	checkIdHashContainer(pt, true)
	mod := pt.Id % idHashMod

	oldHash := ""
	shell, shellExist := roleMap.RoleMap[pt.Role].IdHsashMap[mod].ShellMap[pt.Id]

	if shellExist {
		oldHash = shell.Point.Hash

	} else {

		roleLock := roleMap.RoleMap[pt.Role].Lock
		roleLock.Lock()
		shell, shellExist = roleMap.RoleMap[pt.Role].IdHsashMap[mod].ShellMap[pt.Id]
		if !shellExist {
			shell = createPointShell(pt)
			roleMap.RoleMap[pt.Role].IdHsashMap[mod].ShellMap[pt.Id] = shell
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

func checkIdHashContainer(pt Point, create bool) bool {

	checkRoleContainer(pt, true)
	mod := pt.Id % idHashMod
	_, ok := roleMap.RoleMap[pt.Role].IdHsashMap[mod]
	if ok == false && create == false {
		return false
	}

	roleCon := roleMap.RoleMap[pt.Role]
	roleCon.Lock.Lock()
	defer roleCon.Lock.Unlock()

	idhashCon, ok := roleMap.RoleMap[pt.Role].IdHsashMap[mod]
	if ok == false {
		idhashCon = idHashContainer{ShellMap: make(map[uint64]*PointShell)}
		roleMap.RoleMap[pt.Role].IdHsashMap[mod] = idhashCon
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
			pmap := roleObject{IdHsashMap: make(map[uint64]idHashContainer)}
			roleMap.RoleMap[pt.Role] = pmap
		}
	}
	return true
}

func Summerize() {
	totalPoint := 0
	fmt.Println("-----roleMap.size----", len(roleMap.RoleMap))
	for /*roleid*/ _, son := range roleMap.RoleMap {
		for /*idhash*/ _, idhashCon := range son.IdHsashMap {
			totalPoint += len(idhashCon.ShellMap)
		}

		//fmt.Println("-----roleMap.roleObject=>size----", role, len(son.IdHsashMap))

	}
	fmt.Println("-----total point size----", totalPoint)

}

func DeletePrepare(qr QueryObject) (bool, Point, func(bool)) {

	pt := Point{Id: qr.Id, Role: qr.Role}
	if !checkIdHashContainer(pt, false) {
		return false, Point{}, func(bool) {}
	}

	mod := qr.Id % idHashMod
	shell, ok := roleMap.RoleMap[pt.Role].IdHsashMap[mod].ShellMap[qr.Id]
	if ok == false {
		return false, Point{}, func(bool) {}
	}

	pret := shell.Point
	shellCon := roleMap.RoleMap[pt.Role]
	shellCon.Lock.Lock()

	return true, pret, func(result bool) {
		defer shellCon.Lock.Unlock()
		if result {
			delete(roleMap.RoleMap[pt.Role].IdHsashMap[mod].ShellMap, pt.Id)
		}

	}

}
