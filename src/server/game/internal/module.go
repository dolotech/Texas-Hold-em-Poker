package internal

import (
	"github.com/dolotech/leaf/module"
	"server/base"
	"github.com/golang/glog"
	"server/protocol"
	"github.com/dolotech/leaf/room"
	"reflect"
)

var ( //定义
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)
func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}
func init()  {
	handler(&protocol.JoinRoom{}, room.OnMessage)
	handler(&protocol.LeaveRoom{}, room.OnMessage)
	handler(&protocol.Bet{}, room.OnMessage)
	handler(&protocol.SitDown{}, room.OnMessage) //
	handler(&protocol.StandUp{}, room.OnMessage) //
	handler(&protocol.Chat{}, room.OnMessage)    //
}

type Module struct {
	//相当于继承父类定义
	*module.Skeleton
}

func (m *Module) OnInit() { //继承初始
	m.Skeleton = skeleton

	room.Init(&Creator{})
}

func (m *Module) OnDestroy() {
	glog.Errorln("OnDestroy")
}
