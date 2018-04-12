package internal

import (
	"reflect"
	"server/msg"
	"github.com/dolotech/leaf/gate"
	"server/game"
)

func init() {
	handler(&msg.RegisterUserInfo{}, handlRegisterUserInfo)
	handler(&msg.UserLoginInfo{}, handlLoginUser)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}


func handlRegisterUserInfo(m *msg.RegisterUserInfo, a gate.Agent) {
	//交给 game 模块处理
	game.ChanRPC.Go("RegisterAgent",  m,a,"hello")
	a.WriteMsg(msg.MSG_SUCCESS)

}

func handlLoginUser(m *msg.UserLoginInfo, a gate.Agent) {
	//交给 game
	game.ChanRPC.Go("LoginAgent", m,a)
	a.WriteMsg(msg.MSG_SUCCESS)

}
