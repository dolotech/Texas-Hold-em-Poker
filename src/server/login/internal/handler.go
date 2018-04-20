package internal

import (
	"reflect"
	"server/protocol"
	"github.com/dolotech/leaf/gate"
	"server/game"
	"github.com/golang/glog"
	"server/model"
	"github.com/dolotech/leaf/room"
)

func init() {
	handler(&protocol.UserLoginInfo{}, handlLoginUser)
	handler(&protocol.Version{}, handlVersion)
	handler(&protocol.RoomList{}, onRoomList) //
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handlVersion(m *protocol.Version, a gate.Agent) {
	glog.Infoln(m)
	a.WriteMsg(m)
}

func handlLoginUser(m *protocol.UserLoginInfo, a gate.Agent) {
	user := &model.User{UnionId: m.UnionId}
	exist, err := user.GetByUnionId()
	if err != nil {
		a.WriteMsg(protocol.MSG_DB_Error)
		return
	}

	if !exist {
		user = &model.User{Nickname: m.Nickname,
			UnionId: m.UnionId}
		err := user.Insert()
		if err != nil {
			a.WriteMsg(protocol.MSG_User_Not_Exist)
			return
		}
	}

	resp := &protocol.UserLoginInfoResp{
		Nickname: user.Nickname,
		Account:  user.Account,
		UnionId:  user.UnionId,
	}

	a.WriteMsg(resp)
	game.ChanRPC.Go(model.Agent_Login, user, a)
}

func onRoomList(m *protocol.RoomList, a gate.Agent) {

	msg := &protocol.RoomListResp{}

	array := room.GetRooms()
	rooms := make([]*protocol.Room, len(array))

	for k, v := range array {
		d := v.Data()
		data := d.(*model.Room)
		rooms[k] = &protocol.Room{

			Number:      data.Number,
			MaxCap:      v.Cap(),
			Cap:         v.Len(),
			DraginChips: data.DraginChips,
			CreatedAt:   data.CreatedTime(),
			Rid:         data.Rid,
		}
	}
	msg.Room = rooms
	a.WriteMsg(msg)
}
