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
	handler(&msg.Bet{}, betRoom) //
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func betRoom(m *msg.Bet, a gate.Agent) {
	o := a.UserData().(*Occupant)
	if o.room != nil{
		err:=o.room.Send(o, m)
		if err !=nil{
			a.WriteMsg(msg.MSG_ROOM_CLOSED)
		}
	}else{
		a.WriteMsg(msg.MSG_NOT_IN_ROOM)
	}
}
func leaveRoom(m *msg.LeaveRoom, a gate.Agent) {
	o := a.UserData().(*Occupant)

	if o.room != nil{
		err:=o.room.Send(o, m)
		if err !=nil{
			a.WriteMsg(msg.MSG_ROOM_CLOSED)
		}
	}else{
		a.WriteMsg(msg.MSG_NOT_IN_ROOM)
	}

	glog.Errorln(m, o)
}
func joinRoom(m *msg.JoinRoom, a gate.Agent) {
	o := a.UserData().(*Occupant)

	//todo   从房间列表查找房间
	if o.room != nil {
		//o.room.Send(o,m)
	} else {
		r := NewRoom(&model.Room{})
		model.SetRoom(r)
		r.Send(o, m)
	}

	glog.Errorln(m, o)
}
func handleHello(m *msg.Hello, a gate.Agent) {
	glog.Errorf("hello %v", m.Name)
	a.WriteMsg(&msg.Hello{Name: m.Name,})
}
