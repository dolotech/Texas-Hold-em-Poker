package internal

import (
	"github.com/dolotech/leaf/gate"
	"github.com/golang/glog"
	"server/model"
)

func init() { //与gate 进行"交流"
	skeleton.RegisterChanRPC(model.Agent_New, rpcNewAgent)
	skeleton.RegisterChanRPC(model.Agent_Close, rpcCloseAgent)
	skeleton.RegisterChanRPC(model.Agent_Login, rpcLoginAgent)
}

func rpcNewAgent(a gate.Agent) {
	glog.Errorln("新建链接 ", a)
}

func rpcCloseAgent(a gate.Agent) {
	glog.Errorln("链接关闭 ", a)
}

func rpcLoginAgent(u *model.User, a gate.Agent) {

	o := NewOccupant(u, a)
	a.SetUserData(o)

	if len(u.RoomID) > 0 {
		o.room = model.GetRoom(u.RoomID).(*Room)
	}
	glog.Errorln("rpcLoginAgent", u)
}
