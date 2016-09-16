package point

import (
	"sync"
)

type Point struct {
	Id   uint64
	Lat  float64
	Lng  float64
	Role uint8
	Ext  uint64
}

type PointSetObject struct {
	Point
	Expire uint32
}

type PointQueryObject struct {
	Id   uint64
	Role uint8
}

type PointStorage struct {
	Point
	Update uint32
	//expire uint32
	Lock sync.RWMutex
}

type PointsHashContainer struct {
	pt     *PointStorage
	expire uint32
}

//var PointsCollector = []*Point

type QueryObject struct {
	role uint8
	id   uint64
}

/**

role[role-n][id-hash-n][id]
*/
