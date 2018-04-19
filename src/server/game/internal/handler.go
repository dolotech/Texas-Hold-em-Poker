package internal

import (
	"reflect"
	"server/protocol"
	"github.com/golang/glog"
	"github.com/dolotech/leaf/gate"
	"server/model"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}
func init() {
	handler(&protocol.JoinRoom{}, joinRoom)   //
	handler(&protocol.LeaveRoom{}, onMessage) //
	handler(&protocol.Bet{}, onMessage)       //
	handler(&protocol.SitDown{}, onMessage)   //
	handler(&protocol.StandUp{}, onMessage)   //
	handler(&protocol.RoomList{}, onRoomList) //
	handler(&protocol.Chat{}, onMessage) //
}

func onRoomList(m *protocol.RoomList, a gate.Agent) {
	msg := &protocol.RoomListResp{}

	array := model.GetRooms()
	rooms := make([]*protocol.Room, len(array))

	for k, v := range array {
		rooms[k] = &protocol.Room{
			Number:      v.GetNumber(),
			MaxCap:      v.Cap(),
			Cap:         v.Len(),
			DraginChips: v.GetDragin(),
			CreatedAt:   v.CreatedTime(),
			Rid:         v.ID(),
		}
	}
	msg.Room = rooms
	a.WriteMsg(msg)
}

func onMessage(m interface{}, a gate.Agent) {
	o := a.UserData().(*Occupant)
	if o.room != nil {
		o.room.Send(o, m)
	} else {
		a.WriteMsg(protocol.MSG_NOT_IN_ROOM)
	}
	glog.Errorln(m, o)
}

func joinRoom(m *protocol.JoinRoom, a gate.Agent) {

	o := a.UserData().(*Occupant)

	// 已经在房间
	if o.room != nil {
		// todo 掉线重连处理
		o.room.Send(o, m)
		return
	}

	if len(m.RoomNumber) == 0 {
		r := model.FindRoom()
		r.Send(o, m)
		return
	}
	r := model.GetRoom(m.RoomNumber)
	if r != nil {
		r.Send(o, m)
		return
	}

	room := NewRoom(9, 5, 10, 1000, model.Timeout)
	model.SetRoom(room)
	room.Insert()
	room.Send(o, m)
}
