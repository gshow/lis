package safemap

/**

@author ricolau<ricolau@qq.com>
@version 2016-09-21
@usage

@warning, presently, not validate the "key" for map,
for users,please make sure touse a right key type yourself

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

type tuple struct {
	Key   interface{}
	Value interface{}
}

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

func (this *SafeMap) LockForSet(key interface{}, value interface{}) func(bool) {

	this.lock.Lock()
	return func(callback bool) {

		if callback {
			this.mapdata[key] = value

		}
		this.lock.Unlock()

	}

}

func (this *SafeMap) LockForDelete(key interface{}) (interface{}, func(bool), bool) {

	this.lock.Lock()
	ret, ok := this.mapdata[key]
	if !ok {
		this.lock.Unlock()
		return nil, func(bool) {}, false
	}
	return ret, func(callback bool) {

		if callback {
			delete(this.mapdata, key)

		}
		this.lock.Unlock()

	}, true

}

func (this *SafeMap) Iterate() <-chan tuple {
	ch := make(chan tuple, 1)
	go func() {
		this.lock.Lock()
		for key, val := range this.mapdata {
			ch <- tuple{key, val}
		}

		this.lock.Unlock()
		close(ch)
	}()

	return ch
}

func (this *SafeMap) Set(key interface{}, value interface{}) bool {
	this.lock.Lock()

	this.mapdata[key] = value
	this.usedSize += 1

	this.lock.Unlock()

	return true

}

func (this *SafeMap) SetNotExist(key interface{}, value interface{}) bool {
	this.lock.Lock()

	_, ok := this.mapdata[key]
	if ok {
		this.lock.Unlock()
		return false

	}
	this.mapdata[key] = value

	this.lock.Unlock()
	this.usedSize += 1

	return true
}

func (this *SafeMap) Size() int {
	return this.usedSize
}

func (this *SafeMap) Len() int {
	return this.usedSize
}

func (this *SafeMap) Exist(key interface{}) bool {
	this.lock.Lock()
	_, ok := this.mapdata[key]
	this.lock.Unlock()
	return ok
}

func (this *SafeMap) PositiveGet(key interface{}) interface{} {
	this.lock.Lock()
	value, ok := this.mapdata[key]
	this.lock.Unlock()
	if !ok {
		panic("*SafeMap.PositiveGet()  failed!")
	}
	return value

}

func (this *SafeMap) PositiveMapGet(key interface{}) *SafeMap {
	this.lock.Lock()
	value, ok := this.mapdata[key]
	this.lock.Unlock()
	if !ok {
		panic("*SafeMap.PositiveMapGet()  failed!")
	}
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
