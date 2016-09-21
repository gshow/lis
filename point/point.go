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

//type idHashContainer struct {
//	Lock sync.RWMutex
//	//ShellMap map[uint64]*PointShell
//	ShellMap *smap.SafeMap
//}
//type roleContainer struct {
//	Lock sync.RWMutex
//	//RoleMap map[int]roleObject
//	RoleMap *smap.SafeMap
//}
//type roleObject struct {
//	Lock sync.RWMutex
//	//IdHsashMap map[uint64]idHashContainer
//	IdHsashMap *smap.SafeMap
//}

//var PointsCollector = []*Point

var roleMap = smap.New()

func Query(qr QueryObject) Point {
	pt := Point{Id: qr.Id, Role: qr.Role}
	shellCon, ok := checkPointShellContainer(pt, false)
	if !ok {
		return Point{}
	}
	ishell, ok2 := shellCon.Get(qr.Id)
	if !ok2 {
		return Point{}
	}
	shell := ishell.(*PointShell)

	if !CheckNotExpire(shell) {
		return Point{}
	}
	return shell.Point

}

func SetPrepare(pt Point) (string, *PointShell, func(bool)) {
	//save to roleMap-pointHashContainer-point

	shellcon, _ := checkPointShellContainer(pt, true)

	oldHash := ""
	var shell *PointShell
	ishell, shellExist := shellcon.Get(pt.Id)
	if !shellExist {
		shell = new(PointShell)
	} else {
		shell = ishell.(*PointShell)

	}

	if shellExist {
		oldHash = shell.Point.Hash

	} else {
		shell = createPointShell(pt)

	}
	callback := shellcon.LockForSet(pt.Id, shell)

	//	if tool.Debug() {
	//		fmt.Println("-----point.Set()----", pt, roleMap)
	//		fmt.Println("geohash:", pt.Hash)
	//	}
	return oldHash, shell, callback

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

func checkPointShellContainer(pt Point, create bool) (*smap.SafeMap, bool) {

	var ok bool
	_, ok = roleMap.Get(pt.Role)

	if !ok {
		if !create {
			return smap.New(), false
		} else {
			roleMap.SetNotExist(pt.Role, smap.New())
		}
	}
	hashCon := roleMap.PositiveMapGet(pt.Role)
	mod := pt.Id % idHashMod
	_, ok = hashCon.Get(mod)
	if !ok {
		if !create {
			return smap.New(), false
		} else {
			hashCon.SetNotExist(mod, smap.New())
		}
	}
	ret := hashCon.PositiveGet(mod).(*smap.SafeMap)
	return ret, true
}

func Summerize() {
	totalPoint := 0
	fmt.Println("-----roleMap.size----", roleMap.Size())
	for /*roleid*/ son := range roleMap.Iterate() {
		rolecon := son.Value.(*smap.SafeMap)
		for /*idhash*/ idhashCon := range rolecon.Iterate() {
			totalPoint += idhashCon.Value.(*smap.SafeMap).Size()
		}

		//fmt.Println("-----roleMap.roleObject=>size----", role, len(son.IdHsashMap))

	}
	fmt.Println("-----total point size----", totalPoint)

}

func DeletePrepare(qr QueryObject) (bool, Point, func(bool)) {

	pt := Point{Id: qr.Id, Role: qr.Role}
	shellCon, ok := checkPointShellContainer(pt, false)
	if !ok {
		return false, Point{}, func(bool) {}
	}

	ishell, callback, exist := shellCon.LockForDelete(qr.Id)
	var shell *PointShell
	if !exist {
		shell = new(PointShell)

	} else {
		shell = ishell.(*PointShell)
	}

	return exist, shell.Point, callback

}
