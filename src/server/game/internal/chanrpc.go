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

func rpcNewAgent(args []interface{}) {
	glog.Infoln("--rpcNew--", args)
	a := args[0].(gate.Agent)
	glog.Infoln("args[0]:", a)
	glog.Infoln("len():", len(args))
	for i := 0; i < len(args); i++ {
		//fmt.Fprintln("i=%d,arg[%d]=%v",i,i,args[i])
		glog.Infof("i=%d,arg[%d]=%v \n", i, i, args[i])
	}

	_ = a
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}

func rpcLoginAgent(args []interface{}) {
	glog.Infoln("-rpclon-:", args)
	a := args[0].(gate.Agent)
	glog.Infoln("get m--:", a)
	glog.Infoln("len--:", len(args))
	m := args[1].(*msg.UserLoginInfo)
	err := login(m)
	if err != nil {
		a.WriteMsg(&msg.CodeState{MSG_STATE: msg.MSG_DB_Error})
		return
	}
}

func rpcRigesterAgent(args []interface{}) {
	glog.Infoln("resiter---")
	a := args[0].(gate.Agent)
	m := args[1].(*msg.RegisterUserInfo)
	err := checkExitedUser(m.Name)
	if err == nil {
		a.WriteMsg(&msg.CodeState{MSG_STATE: msg.MSG_Register_Existed})
		return
	}
	err = register(m)
	if err != nil {
		a.WriteMsg(&msg.CodeState{MSG_STATE: msg.MSG_DB_Error})
		return
	}
}

func rpcJoinRoomAgent(args []interface{}) {

}
