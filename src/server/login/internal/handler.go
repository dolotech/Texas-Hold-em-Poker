package internal

import (
	"reflect"
	"server/msg"
	"github.com/dolotech/leaf/gate"
	"server/game"
	"github.com/golang/glog"
)

func init() {
	handler(&msg.RegisterUserInfo{}, handlRegisterUserInfo)
	handler(&msg.UserLoginInfo{}, handlLoginUser)
	handler(&msg.Version{}, handlVersion)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handlVersion(m *msg.Version, a gate.Agent) {
	glog.Infoln(m)
}

func handlRegisterUserInfo(m *msg.RegisterUserInfo, a gate.Agent) {
	//交给 game 模块处理
	game.ChanRPC.Go("RegisterAgent", m, a, "hello")
	//a.WriteMsg(msg.MSG_SUCCESS)

}

func handlLoginUser(m *msg.UserLoginInfo, a gate.Agent) {
	//交给 game
	game.ChanRPC.Go("LoginAgent", m, a)
	//a.WriteMsg(msg.MSG_SUCCESS)

}
