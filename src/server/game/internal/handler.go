package internal

import (
	"reflect"
	"server/msg"
	"github.com/golang/glog"
	"github.com/dolotech/leaf/gate"
	"server/model"
)

func init() {
	handler(&msg.Hello{}, handleHello)   //具体处理函数调用
	handler(&msg.JoinRoom{}, joinRoom)   //
	handler(&msg.LeaveRoom{}, leaveRoom) //
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func leaveRoom(m *msg.LeaveRoom,a gate.Agent) {
	o:= a.UserData().(*Occupant)

	o.room.Send(o,m)
	glog.Errorln(m,o)
}
func joinRoom(m *msg.JoinRoom,a gate.Agent) {
	o:= a.UserData().(*Occupant)

	if o.room !=nil{
		//o.room.Send(o,m)
	}else{
		r:= NewRoom(&model.Room{})
		model.SetRoom(r)
		r.Send(o,m)
	}

	glog.Errorln(m,o)
}
func handleHello(m *msg.Hello, a gate.Agent) {
	glog.Errorf("hello %v", m.Name)
	a.WriteMsg(&msg.Hello{Name: m.Name,})
}
