package internal

import (
	"reflect"
	"server/msg"
	"github.com/dolotech/leaf/gate"
	"server/game"
	"github.com/golang/glog"
	"server/model"
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

	u := &model.User{Account: m.Name}

	_, err := u.GetByAccount()

	if err != nil {
		a.WriteMsg(msg.MSG_User_Not_Exist)
		return
	}

	glog.Infoln("login success",m)
	a.SetUserData(u)
	//game.ChanRPC.Go("LoginAgent", m, a)
	a.WriteMsg(msg.MSG_SUCCESS)

}
