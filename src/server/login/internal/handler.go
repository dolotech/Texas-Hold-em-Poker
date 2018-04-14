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
	//handler(&msg.JoinRoom{}, joinRoom)
}

/*func joinRoom(m *msg.JoinRoom, a gate.Agent) {
	glog.Errorln(m)
}*/

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

	user := &model.User{UnionId: m.UnionId}
	exist, _ := user.GetByUnionId()
	if !exist {
		user = &model.User{Nickname: m.Nickname,
			UnionId: m.UnionId}
		err := user.Insert()
		if err != nil {
			a.WriteMsg(msg.MSG_User_Not_Exist)
			return
		}
	}

	//o := game.NewOccupant(a.UserData().(*model.User), a)

	resp := &msg.UserLoginInfoResp{
		Nickname: user.Nickname,
		Account:  user.Account,
		UnionId:  user.UnionId,
	}
	glog.Infoln("login success", m)
	//a.SetUserData(user)

	a.WriteMsg(resp)
	game.ChanRPC.Go("LoginAgent", user, a)
}
