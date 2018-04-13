package internal

import (
	//"leaf/gate"
	"github.com/dolotech/leaf/gate"
	"server/model"
)

type UserLine struct {
	gate.Agent //申请代理
	//*g.LinearContext
	Cards    []*CardData //牌
	UserData *model.User
	RoomLine *Room
	//RoomPosition	int	//房间里的座位号 1-4
	ReadyState int //准备状态 是否可以开局 0否 1是
}

//初始化 用户登录后 进行
func (u *UserLine) initUser(userInfo *model.User) {
	u.UserData = userInfo
	u.RoomLine = nil
	//u.Cards = nil
	//u.LinearContext = skeleton.NewLinearContext()
}

/*
//点击准备 异或操作
func (u *UserLine) readyGame() {
	u.ReadyState = u.ReadyState ^ 1
	if u.ReadyState > 0 {
		u.RoomLine.ReadySignCount ++
	}
}
*/
