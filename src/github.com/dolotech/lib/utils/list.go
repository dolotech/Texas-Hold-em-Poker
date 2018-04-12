// 线程安全的数组封装，注意：写锁接口不能嵌套调用，比如：Range接口不能调用删除接口Del
package utils

import (
	"sync"
)

func NewList() *List {
	return &List{elems: make([]interface{}, 0)}
}

type List struct {
	sync.RWMutex
	elems []interface{}
}

func (this *List) Get(f func(interface{}) bool) interface{} {
	this.RLock()
	defer this.RUnlock()
	for _, v := range this.elems {
		if f(v) {
			return v
		}
	}
	return nil
}

func (this *List) Add(value interface{}) {
	this.Lock()
	defer this.Unlock()
	this.elems = append(this.elems, value)
}

func (this *List) Pure() {
	this.Lock()

	defer this.Unlock()
	this.elems = this.elems[:0]

}
func (this *List) Del(value interface{}) {
	this.Lock()
	defer this.Unlock()
	for k, v := range this.elems {
		if value == v {
			this.elems = append(this.elems[:k], this.elems[k+1:]...)
			break
		}
	}
}

func (this *List) Delete(f func(interface{}) bool) {
	this.Lock()
	defer this.Unlock()
	for k, v := range this.elems {
		if f(v) {
			this.elems = append(this.elems[:k], this.elems[k+1:]...)
			break
		}
	}
}

func (this *List) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.elems)
}

func (this *List) Replace(value interface{}, f func(interface{}) bool) {
	this.RLock()
	defer this.RUnlock()
	for k, v := range this.elems {
		if f(v) {
			this.elems = append(this.elems[:k], this.elems[k+1:]...)
			this.elems = append(this.elems, value)
			break
		}
	}
}

func (this *List) Range(f func(interface{}) bool) {
	this.RLock()
	defer this.RUnlock()
	for _, v := range this.elems {
		if f(v) {
			break
		}
	}
}

func (this *List) LRange(f func(interface{}) bool) {
	this.Lock()
	defer this.Unlock()
	for _, v := range this.elems {
		if f(v) {
			break
		}
	}
}
