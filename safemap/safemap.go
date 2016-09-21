package safemap

/**

@author ricolau<ricolau@qq.com>
@version 201-09-21
@usage


	b := smap.New()
    for i:=0;i<10000;i++{
        go func(i int){
            b.Set(i,i+1000)
        }(i)
        go func (i int){
            b.Delete(i+10)
        }(i)

    }
    time.Sleep(time.Second * 2)
    p(b.Size())


*/
import (
	//"math"
	"sync"
)

var defaultMapSize = 10

type SafeMap struct {
	lock     sync.RWMutex
	size     int
	usedSize int
	mapdata  map[interface{}]interface{}
}

func New() *SafeMap {
	s := &SafeMap{size: defaultMapSize, usedSize: 0, mapdata: make(map[interface{}]interface{}, defaultMapSize)}

	return s
}

func (this *SafeMap) Lock() {
	this.lock.Lock()
}

func (this *SafeMap) Unlock() {
	this.lock.Unlock()
}
func (this *SafeMap) Set(key interface{}, value interface{}) bool {
	this.lock.Lock()
	//	if this.usedSize+6 <= this.size {

	//		newSize := int(math.Ceil(float64(this.size) * 1.5))

	//		newMap := make(map[interface{}]interface{}, newSize)

	//		for k, v := range this.mapdata {
	//			newMap[k] = v
	//		}
	//		this.mapdata = newMap
	//	}

	this.mapdata[key] = value
	this.usedSize += 1

	this.lock.Unlock()

	return true

}

func (this *SafeMap) Size() int {
	return this.size
}

func (this *SafeMap) Exist(key interface{}) bool {
	this.lock.Lock()
	_, ok := this.mapdata[key]
	this.lock.Unlock()
	return ok

}

func (this *SafeMap) PositiveGet(key interface{}) interface{} {
	this.lock.Lock()
	value, _ := this.mapdata[key]
	this.lock.Unlock()
	return value

}

func (this *SafeMap) PositiveLinkGet(key interface{}) *SafeMap {
	this.lock.Lock()
	value, _ := this.mapdata[key]
	this.lock.Unlock()
	return value.(*SafeMap)

}

func (this *SafeMap) Get(key interface{}) (interface{}, bool) {
	this.lock.Lock()
	value, ok := this.mapdata[key]
	this.lock.Unlock()
	return value, ok

}

func (this *SafeMap) Delete(key interface{}) bool {
	this.lock.Lock()
	_, ok := this.mapdata[key]
	if ok {
		delete(this.mapdata, key)
		this.usedSize -= 1
	}
	this.lock.Unlock()
	return true
}
