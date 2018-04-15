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
	handler(&msg.Bet{}, betRoom)         //
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func betRoom(m *msg.Bet, a gate.Agent) {
	o := a.UserData().(*Occupant)
	if o.room != nil {
		o.room.Send(o, m)

	} else {
		a.WriteMsg(msg.MSG_NOT_IN_ROOM)
	}
}
func leaveRoom(m *msg.LeaveRoom, a gate.Agent) {
	o := a.UserData().(*Occupant)

	if o.room != nil {
		o.room.Send(o, m)
	} else {
		a.WriteMsg(msg.MSG_NOT_IN_ROOM)
	}

	glog.Errorln(m, o)
}

func joinRoom(m *msg.JoinRoom, a gate.Agent) {

	o := a.UserData().(*Occupant)

	// 已经在房间
	if o.room != nil {
		// todo 掉线重连处理
		o.room.Send(o, m)
		return
	}
	var r model.IRoom
	if len(m.RoomNumber) == 0 {
		r = model.FindRoom()
	} else {
		r = model.GetRoom(m.RoomNumber)
	}
	if r == nil {
		r = NewRoom(9, 5, 10)
		model.SetRoom(r)
	}
	r.Send(o, m)
}
func handleHello(m *msg.Hello, a gate.Agent) {
	glog.Errorf("hello %v", m.Name)
	a.WriteMsg(&msg.Hello{Name: m.Name,})
}
