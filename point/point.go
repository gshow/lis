package point

import (
	"sync"
	//"lis/role"
	"fmt"
	"time"
)

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

type PointQueryObject struct {
	Id   uint64
	Role uint8
}

type ShellContainer struct {
	Lock  sync.RWMutex
	Sdata map[uint64]*PointShell
}

//var PointsCollector = []*Point

type QueryObject struct {
	role uint8
	id   uint64
}

type RoleContainer struct {
	Lock  sync.RWMutex
	Rdata map[uint8]ShellContainer
}

var roleMap = RoleContainer{Rdata: make(map[uint8]ShellContainer)}

func Set(pt Point) (bool, string, *PointShell) {
	//save to roleMap-pointHashContainer-point

	pt.Update = time.Now().Unix()
	_, ok := roleMap.Rdata[pt.Role]

	//	fmt.Println(11111, pt)
	oldHash := ""
	if ok == false {
		roleMap.Lock.Lock()
		defer roleMap.Lock.Unlock()

		_, ok := roleMap.Rdata[pt.Role]
		if ok == false {
			pmap := ShellContainer{Sdata: make(map[uint64]*PointShell)}
			roleMap.Rdata[pt.Role] = pmap
		}
	}
	shell, ok := roleMap.Rdata[pt.Role].Sdata[pt.Id]
	if ok == true {
		oldHash = shell.Point.Hash
		shell.Lock.Lock()

		roleMap.Rdata[pt.Role].Sdata[pt.Id].Point = pt
	} else {
		shell := &PointShell{Point: pt}
		roleMap.Rdata[pt.Role].Sdata[pt.Id] = shell
		shell.Lock.Lock()
	}

	fmt.Println("-----point.Set()----", pt, roleMap)
	return true, oldHash, roleMap.Rdata[pt.Role].Sdata[pt.Id]

}

func Delete(pt PointQueryObject) {

}
