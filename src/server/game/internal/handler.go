package internal

import (
	"reflect"
	"server/msg"
	"github.com/golang/glog"
	"github.com/dolotech/leaf/gate"
)

func init() {
	handler(&msg.Hello{}, handleHello) //具体处理函数调用
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleHello(m *msg.Hello, a gate.Agent) {
	glog.Errorf("hello %v", m.Name)
	a.WriteMsg(&msg.Hello{Name: m.Name,})
}
