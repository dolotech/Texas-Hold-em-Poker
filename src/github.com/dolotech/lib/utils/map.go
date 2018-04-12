package utils

import (
	"sync"
)

func NewMap() *Map {
	return &Map{elems: map[interface{}]interface{}{}}
}

type Map struct {
	sync.RWMutex
	elems map[interface{}]interface{}
}

func (this *Map) Get(key interface{}) interface{} {
	this.RLock()
	defer this.RUnlock()
	if value, ok := this.elems[key]; ok {
		return value
	} else {
		return nil
	}
}
func (this *Map) Set(key interface{}, value interface{}) {
	this.Lock()
	defer this.Unlock()
	this.elems[key] = value
}

func (this *Map) Del(key interface{}) {
	this.Lock()
	defer this.Unlock()
	delete(this.elems, key)
}

func (this *Map) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.elems)
}

func (this *Map) Range(f func(interface{}, interface{}) bool) {
	this.RLock()
	defer this.RUnlock()
	for k, v := range this.elems {
		if f(k, v) {
			break
		}
	}
}

func (this *Map) LRange(f func(interface{}, interface{}) bool) {
	this.Lock()
	defer this.Unlock()
	for k, v := range this.elems {
		if f(k, v) {
			break
		}
	}
}
