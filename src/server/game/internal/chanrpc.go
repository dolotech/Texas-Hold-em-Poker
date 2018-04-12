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
	//glog.Infoln("--rpcNew--", args)
	//a := args[0].(gate.Agent)
	glog.Errorln("rpcNewAgent ",a)


	//glog.Infoln("len():", len(args))
	//for i := 0; i < len(args); i++ {
		//fmt.Fprintln("i=%d,arg[%d]=%v",i,i,args[i])
		//glog.Infof("i=%d,arg[%d]=%v \n", i, i, args[i])
	//}

	//_ = a
}

func rpcCloseAgent(a gate.Agent)  {

	glog.Errorln("rpcCloseAgent ",a)
	//a := args[0].(gate.Agent)
	//_ = a
}

func rpcLoginAgent(m *msg.UserLoginInfo,a gate.Agent)  {
	//glog.Infoln("-rpclon-:", args)
	//a := args[0].(gate.Agent)
	glog.Errorln("rpcLoginAgent", m)
	//glog.Infoln("len--:", len(args))
	//m := args[1].(*msg.UserLoginInfo)
	err := login(m)
	if err != nil {
		a.WriteMsg( msg.MSG_DB_Error)
		return
	}
}

func rpcRigesterAgent(m *msg.RegisterUserInfo,a gate.Agent)   {
	glog.Errorln("rpcRigesterAgent---",m)
	//a := args[0].(gate.Agent)
	//m := args[1].(*msg.RegisterUserInfo)
	err := checkExitedUser(m.Name)
	if err == nil {
		a.WriteMsg(msg.MSG_Register_Existed)
		return
	}
	err = register(m)
	if err != nil {
		a.WriteMsg(msg.MSG_DB_Error)
		return
	}
}

func rpcJoinRoomAgent(args []interface{}) {

}
