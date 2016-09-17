package location

import (
	"sync"

	"obis/point"
)

type Container struct {
	//precision 6
	Hash_value string
	Points     []*point.Point
}

type ContainerMap struct {
	Data map[uint8]*Container
	Lock sync.RWMutex
}

var ContainerMapAll = ContainerMap{}
