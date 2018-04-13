package utils

import (
	"reflect"
	"github.com/golang/glog"
)

type Route struct {
	hash map[string]*reflect.Value
}

func (this *Route) Route(msg interface{}, arg ...interface{}) {
	msgID := reflect.TypeOf(msg).Elem().Name()
	if f, ok := this.hash[msgID]; ok {
		array := make([]reflect.Value, len(arg)+1)
		array[0] = reflect.ValueOf(msg)
		for i := 0; i < len(arg); i++ {
			array[i+1] = reflect.ValueOf(arg[i])
		}
		f.Call(array)
	}
}
func (this *Route) Regist(msg, f interface{}) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		glog.Error("message pointer required")
	}
	msgID := msgType.Elem().Name()
	if msgID == "" {
		glog.Error("unnamed  message")
	}
	if _, ok := this.hash[msgID]; ok {
		glog.Errorf("message %v is already registered", msgID)
	}
	v := reflect.ValueOf(f)
	this.hash[msgID] = &v
}

func NewRoute() *Route {
	return &Route{
		hash: make(map[string]*reflect.Value),
	}
}
