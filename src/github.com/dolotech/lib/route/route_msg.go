package route

import (
	"reflect"
	"github.com/golang/glog"
)

type Route map[string]*reflect.Value

func (this *Route) Emit(msg interface{}, arg ...interface{}) {
	if *this != nil {
		msgID := reflect.TypeOf(msg).Elem().Name()
		if f, ok := (*this)[msgID]; ok {
			array := make([]reflect.Value, len(arg)+1)
			array[0] = reflect.ValueOf(msg)
			for i := 0; i < len(arg); i++ {
				array[i+1] = reflect.ValueOf(arg[i])
			}
			f.Call(array)
		}
	}
}
func (this *Route) Regist(msg, f interface{}) {
	if *this == nil {
		*this = make(map[string]*reflect.Value, 18)
	}
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		glog.Error("message pointer required")
	}
	msgID := msgType.Elem().Name()
	if msgID == "" {
		glog.Error("unnamed  message")
	}
	if _, ok := (*this)[msgID]; ok {
		glog.Errorf("message %v is already registered", msgID)
	}
	v := reflect.ValueOf(f)
	(*this)[msgID] = &v
}
