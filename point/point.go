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

func SetWithAggrement(pt Point, aggrement func(key interface{}, value interface{}, oldvalue interface{}) bool) bool {
	//save to roleMap-pointHashContainer-point

	shellcon, _ := checkPointShellContainer(pt, true)

	//oldHash := ""
	var shell *PointShell
	ishell, shellExist := shellcon.Get(pt.Id)
	if !shellExist {
		shell = new(PointShell)
	} else {
		shell = ishell.(*PointShell)

	}

	if !shellExist {
		shell = createPointShell(pt)

	}
	setResult := shellcon.SetWithAggrement(pt.Id, shell, aggrement)

	//	if tool.Debug() {
	//		fmt.Println("-----point.Set()----", pt, roleMap)
	//		fmt.Println("geohash:", pt.Hash)
	//	}
	return setResult

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
	var rolecon *smap.SafeMap

	rangeCallSon := func(k interface{}, value interface{}) bool {
		rolecon = value.(*smap.SafeMap)
		totalPoint += value.(*smap.SafeMap).Size()
		return true

	}
	rangeCall := func(k interface{}, value interface{}) bool {
		rolecon := value.(*smap.SafeMap)
		rolecon.Range(rangeCallSon)
		return true

	}
	roleMap.Range(rangeCall)

	fmt.Println("-----roleMap.size----", roleMap.Size())

	fmt.Println("-----total point size----", totalPoint)

}

func DeleteWithAggrement(qr QueryObject, aggrement func(key interface{}, oldvalue interface{}) bool) bool {
	//save to roleMap-pointHashContainer-point

	pt := Point{Id: qr.Id, Role: qr.Role}

	shellCon, ok := checkPointShellContainer(pt, false)
	if !ok {
		return false
	}

	result := shellCon.DeleteWithAggrement(pt.Id, aggrement)

	return result

}
