package safemap

/**

@author ricolau<ricolau@qq.com>
@version 2016-09-21
@usage

@warning, presently, there's not validation for the "key" of map,
for users,please make sure to use a right key type yourself

///=====code=====///
package main
import(
  smap "github.com/ricolau/safemap"
 "fmt"
 "time"
)
func main(){
    b := smap.New()
    for i:=0;i<1000;i++{
        go func(i int){
            b.Set(i,i+1000)
        }(i)

        go func(i int){
            b.Exist(i)
        }(i)
        go func (i int){
            b.Get(i+10)
        }(i)
        go func (i int){
            b.Delete(i+10)
        }(i)
    }
    time.Sleep(time.Second * 2)
    fmt.Println(b.Size())
}
///=====code=====///


*/
import (
	"sync"
)

type SafeMap struct {
	lock     sync.RWMutex
	size     int
	usedSize int
	mapdata  map[interface{}]interface{}
}

func New() *SafeMap {
	s := &SafeMap{usedSize: 0, mapdata: make(map[interface{}]interface{})}

	return s
}

func (smap *SafeMap) SetWithAggrement(key interface{}, value interface{}, aggrement func(key interface{}, value interface{}, oldvalue interface{}) bool) bool {
	smap.lock.Lock()
	defer smap.lock.Unlock()

	oldvalue, _ := smap.mapdata[key]

	aggred := aggrement(key, value, oldvalue)
	if aggred {
		smap.mapdata[key] = value
	}
	return aggred
}

func (smap *SafeMap) DeleteWithAggrement(key interface{}, aggrement func(key interface{}, oldvalue interface{}) bool) bool {
	smap.lock.Lock()
	defer smap.lock.Unlock()

	oldvalue, ok := smap.mapdata[key]

	aggred := aggrement(key, oldvalue)
	ifDelete := aggred && ok
	if ifDelete {
		delete(smap.mapdata, key)
	}
	return ifDelete
}

func (smap *SafeMap) Range(callback func(key interface{}, value interface{}) bool) {
	smap.lock.Lock()
	defer smap.lock.Unlock()

	for k, v := range smap.mapdata {
		continued := callback(k, v)
		if !continued {
			break
		}

	}

}

func (smap *SafeMap) Set(key interface{}, value interface{}) bool {
	smap.lock.Lock()

	smap.mapdata[key] = value
	smap.usedSize += 1

	smap.lock.Unlock()

	return true

}

func (smap *SafeMap) SetNotExist(key interface{}, value interface{}) bool {
	smap.lock.Lock()

	_, ok := smap.mapdata[key]
	if ok {
		smap.lock.Unlock()
		return false

	}
	smap.mapdata[key] = value

	smap.lock.Unlock()
	smap.usedSize += 1

	return true
}

func (smap *SafeMap) Size() int {
	return smap.usedSize
}

func (smap *SafeMap) Len() int {
	return smap.usedSize
}

func (smap *SafeMap) Exist(key interface{}) bool {
	smap.lock.Lock()
	_, ok := smap.mapdata[key]
	smap.lock.Unlock()
	return ok
}

func (smap *SafeMap) PositiveGet(key interface{}) interface{} {
	smap.lock.Lock()
	value, ok := smap.mapdata[key]
	smap.lock.Unlock()
	if !ok {
		panic("*SafeMap.PositiveGet()  failed!")
	}
	return value

}

func (smap *SafeMap) PositiveMapGet(key interface{}) *SafeMap {
	smap.lock.Lock()
	value, ok := smap.mapdata[key]
	smap.lock.Unlock()
	if !ok {
		panic("*SafeMap.PositiveMapGet()  failed!")
	}
	return value.(*SafeMap)

}

func (smap *SafeMap) Get(key interface{}) (interface{}, bool) {
	smap.lock.Lock()
	value, ok := smap.mapdata[key]
	smap.lock.Unlock()
	return value, ok
}

func (smap *SafeMap) Delete(key interface{}) bool {
	smap.lock.Lock()
	_, ok := smap.mapdata[key]
	if ok {
		delete(smap.mapdata, key)
		smap.usedSize -= 1
	}
	smap.lock.Unlock()
	return true
}
