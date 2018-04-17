package gate

import (
	"server/protocol"
	"server/game"
	"server/login"
)

func init() {
	protocol.Processor.SetRouter(&protocol.UserLoginInfo{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.JoinRoom{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.LeaveRoom{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.SitDown{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.StandUp{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Showdown{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Pot{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Button{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.PreFlop{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Bet{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Version{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.RoomList{}, game.ChanRPC)
}
