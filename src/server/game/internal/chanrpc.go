package internal

import (
	"github.com/dolotech/leaf/gate"
	"server/msg"
	"github.com/golang/glog"
)

func init() { //与gate 进行"交流"
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("LoginAgent", rpcLoginAgent)
	skeleton.RegisterChanRPC("RegisterAgent", rpcRigesterAgent)
}

func rpcNewAgent(a gate.Agent) {
	glog.Errorln("rpcNewAgent ",a)
}

func rpcCloseAgent(a gate.Agent)  {
	glog.Errorln("rpcCloseAgent ",a)
}

func rpcLoginAgent(m *msg.UserLoginInfo,a gate.Agent)  {
	glog.Errorln("rpcLoginAgent", m)
	/*err := login(m)
	if err != nil {
		a.WriteMsg( msg.MSG_DB_Error)
		return
	}*/
}

func rpcRigesterAgent(m *msg.RegisterUserInfo,a gate.Agent,str string)   {
	glog.Errorln("rpcRigesterAgent---",m,str)
	/*if checkExitedUser(m.Name) {
		a.WriteMsg(msg.MSG_Register_Existed)
		return
	}
	err := register(m)
	if err != nil {
		a.WriteMsg(msg.MSG_DB_Error)
		return
	}*/
}

func rpcJoinRoomAgent(args []interface{}) {

}
