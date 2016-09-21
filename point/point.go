package point

import (
	"sync"
	//"lis/role"
	"fmt"
	//"lis/tool"
	smap "lis/safemap"
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

type QueryResultObject struct {
	Id  uint64
	Lat float64
	Lng float64
	//Hash string
	Role int

	Ext    int64
	Update int
	//Expire int
}

type idHashContainer struct {
	Lock sync.RWMutex
	//ShellMap map[uint64]*PointShell
	ShellMap *smap.SafeMap
}
type roleContainer struct {
	Lock sync.RWMutex
	//RoleMap map[int]roleObject
	RoleMap *smap.SafeMap
}
type roleObject struct {
	Lock sync.RWMutex
	//IdHsashMap map[uint64]idHashContainer
	IdHsashMap *smap.SafeMap
}

//var PointsCollector = []*Point

var roleMap = roleContainer{RoleMap: smap.New()}

func Query(qr QueryObject) Point {
	pt := Point{Id: qr.Id, Role: qr.Role}

	if !checkIdHashContainer(pt, false) {

		return Point{}
	}

	shell, ok := roleMap.RoleMap.PositiveLinkGet(pt.Role).PositiveLinkGet(qr.Id % idHashMod).Get(qr.Id)
	if ok == false {
		return Point{}
	}
	shellt := shell.(*PointShell)
	if !CheckNotExpire(shellt) {
		return Point{}
	}
	return shellt.Point

}

func SetPrepare(pt Point) (string, *PointShell, func(bool)) {
	//save to roleMap-pointHashContainer-point

	checkIdHashContainer(pt, true)
	mod := pt.Id % idHashMod

	oldHash := ""
	ishell, shellExist := roleMap.RoleMap.PositiveLinkGet(pt.Role).PositiveLinkGet(mod).Get(pt.Id)
	shell := ishell.(*PointShell)
	var roleLock sync.RWMutex
	if shellExist {
		oldHash = shell.Point.Hash

	} else {

		idhashCon := roleMap.RoleMap.PositiveLinkGet(pt.Role).PositiveGet(mod)
		tCon := idhashCon.(*idHashContainer)
		tCon.Lock.Lock()

		ishell, shellExist = roleMap.RoleMap.PositiveLinkGet(pt.Role).PositiveLinkGet(mod).Get(pt.Id)
		if !shellExist {
			shell = createPointShell(pt)
			roleMap.RoleMap.PositiveLinkGet(pt.Role).PositiveLinkGet(mod).Set(pt.Id, shell)
		}
		shell = ishell.(*PointShell)
		tCon.Lock.Unlock()

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

	if pshell.Point.Expire > 0 && int(time.Now().Unix()) > pshell.Point.Expire {
		go ExpireQueueAdd(pshell)
		return false
	}
	return true

}

func checkIdHashContainer(pt Point, create bool) bool {

	ok := checkRoleContainer(pt, create)
	if !create && !ok { //not exist and no create
		return ok
	}
	mod := pt.Id % idHashMod
	idhashCon, ok := roleMap.RoleMap.PositiveLinkGet(pt.Role).Get(mod)
	if ok == false {
		iroleCon := roleMap.RoleMap.PositiveGet(pt.Role)
		roleCon := iroleCon.(roleContainer)
		roleCon.Lock.Lock()

		_, ok2 := roleMap.RoleMap.PositiveLinkGet(pt.Role).Get(mod)
		if ok2 == false {
			idhashCon = idHashContainer{ShellMap: smap.New()}
			roleCon.Set(mod, idhashCon)

		}
		roleCon.Lock.Unlock()

	}

	//	roleCon := roleMap.RoleMap.Get(pt.Role)
	//	roleCon.Lock

	return true
}

func checkRoleContainer(pt Point, create bool) bool {
	ok := roleMap.RoleMap.Exist(pt.Role)
	if ok == false && create == false {
		return false
	}
	if ok == false && create == true {
		roleMap.Lock.Lock()
		defer roleMap.Lock.Unlock()

		ok2 := roleMap.RoleMap.Exist(pt.Role)
		if ok2 == false {
			pmap := roleObject{IdHsashMap: smap.New()}
			roleMap.RoleMap.Set(pt.Role, pmap)
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
