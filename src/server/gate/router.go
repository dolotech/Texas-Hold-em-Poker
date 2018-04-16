package gate

import (
	"server/msg"
	"server/game"
	"server/login"
)

func init() {
	msg.Processor.SetRouter(&msg.Hello{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.UserLoginInfo{}, login.ChanRPC)
	msg.Processor.SetRouter(&msg.JoinRoom{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.LeaveRoom{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.SitDown{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.StandUp{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Showdown{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Pot{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Button{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.PreFlop{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Bet{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Fold{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Version{}, game.ChanRPC)
}
