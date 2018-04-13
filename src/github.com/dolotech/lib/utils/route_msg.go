package utils

import (
	"reflect"
	"github.com/golang/glog"
)

type Route struct {
	hash map[string]*reflect.Value
}

func (this *Route) Route(msg ,arg interface{}) {
	msgType := reflect.TypeOf(msg)
	msgID := msgType.Elem().Name()
	if f, ok := this.hash[msgID]; ok {
		//glog.Errorln("Route")
		f.Call([]reflect.Value{reflect.ValueOf(msg),reflect.ValueOf(arg)})
	}
}
func (this *Route) Regist(msg, f interface{}) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		glog.Fatal("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	if msgID == "" {
		glog.Fatal("unnamed json message")
	}
	if _, ok := this.hash[msgID]; ok {
		glog.Fatal("message %v is already registered", msgID)
	}
	v := reflect.ValueOf(f)
	this.hash[msgID] = &v
}

func NewRoute() *Route {
	return &Route{
		hash: make(map[string]*reflect.Value),
	}
}

//var route = NewRoute()
