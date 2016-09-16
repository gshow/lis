package location

import (
	"sync"

	"github.com/gshow/obis/point"
)

type Container struct {
	//precision 6
	Hash_value string
	Points     []*point.Point
}

type ContainerMap struct {
	data map[string]*Container
	lock sync.RWMutex
}

var ContainerMapAll = ContainerMap{}
