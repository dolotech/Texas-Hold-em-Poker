package internal

import (
	"reflect"
	"server/msg"
	"github.com/name5566/leaf/gate"
	"github.com/golang/glog"
)

func init()  {
	handler(&msg.Hello{}, handleHello)//具体处理函数调用
	
}

func handler(m interface{}, h interface{})  {
	skeleton.RegisterChanRPC(reflect.TypeOf(m),h)
	
}

func handleHello(args []interface{})  {
	m := args[0].(*msg.Hello)
	a := args[1].(gate.Agent)

	glog.Errorf("hello %v",m.Name)
	a.WriteMsg(&msg.Hello{ Name : m.Name, })
}