package point

import (
	"sync"
)

type Point struct {
	Id   uint64
	Lat  float64
	Lng  float64
	Hash string
	Role uint8

	Ext    uint64
	Update uint32
	Expire uint32
}

type PointQueryObject struct {
	Id   uint64
	Role uint8
}

type pointContainer struct {
	Hash string
	Data []Point
}

//var PointsCollector = []*Point

type QueryObject struct {
	role uint8
	id   uint64
}

func Set(pt PointSetObject) {
	//save to roleMap-pointHashContainer-point

	_, ok := roleMap.Data[pso.Role]
	if 1 {

	}

}

func Delete(pt PointQueryObject) {

}

type role struct {
	Id   int32
	pCon pointContainer
}

type roleMap struct {
	Data map[string]*string
}

var roleMap = make(map[uint8]*RoleContainer)

/**

roleMap.Data[string][]role

role[string]*pointContainer

PointContainer[]  point


role[role-n][id-hash-n][id]
*/
