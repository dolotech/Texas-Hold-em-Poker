package internal

import (
	"reflect"
	"server/msg"
	"github.com/name5566/leaf/gate"
	"server/game"
)

func init() {
	handler(&msg.RegisterUserInfo{},handlRegisterUserInfo)
	handler(&msg.UserLoginInfo{},handlLoginUser)
}

func handler(m interface{}, h interface{})  {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h )
}

func handlRegisterUserInfo(args []interface{})  {
	//收到注册信息
	m := args[0].(*msg.RegisterUserInfo)
	//获取发送者
	a := args[1].(gate.Agent)

	//交给 game 模块处理
	game.ChanRPC.Go("RegisterAgent",a, m )

	a.WriteMsg(&msg.CodeState{MSG_STATE:msg.MSG_Register_OK})

}

func handlLoginUser(args []interface{})  {
	m := args[0].(*msg.UserLoginInfo)
	a := args[1].(gate.Agent)

	//交给 game
	game.ChanRPC.Go("LoginAgent", a, m )
	a.WriteMsg(&msg.CodeState{MSG_STATE:msg.MSG_Login_OK})
	
}
