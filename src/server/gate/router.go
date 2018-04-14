package gate

import (
	"server/msg"
	"server/game"
	"server/login"
)

//消息在此进行交割
func init() {
	msg.Processor.SetRouter(&msg.Hello{},game.ChanRPC)//参数消息内容 通信桥chanRPC
	//msg.Processor.SetRouter(&msg.UserLoginInfo{},game.ChanRPC)

	//用注册
	msg.Processor.SetRouter(&msg.RegisterUserInfo{},login.ChanRPC)
	//登录
	msg.Processor.SetRouter(&msg.UserLoginInfo{},login.ChanRPC)

	msg.Processor.SetRouter(&msg.JoinRoom{},game.ChanRPC)
	msg.Processor.SetRouter(&msg.LeaveRoom{},game.ChanRPC)
}
