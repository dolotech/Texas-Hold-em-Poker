package internal

import (
	"server/msg"
	"github.com/golang/glog"
	"github.com/dolotech/leaf/gate"
	"server/model"
)

/*

//共有四轮发牌，按顺序分别为：
ActPreflop = "preflop" //底牌
ActFlop    = "flop"    // 翻牌
ActTurn    = "turn"    // 转牌
ActRiver   = "river"   // 河牌

ActShowdown = "showdown" //摊牌和比牌
ActPot      = "pot"      //通报奖池

ActJoin   = "join"   //通报加入游戏的玩家
ActLeave  = "gone"   //用户离开房间
ActBet    = "bet"    //玩家下注
ActButton = "button" //通报本局庄家
ActState  = "state"  //房间信息*/

func (r *Room) joinRoom(m *msg.JoinRoom, a gate.Agent) {
	o := NewOccupant(a.UserData().(*model.User), a)
	r.addOccupant(o)
	// todo 下发玩家的座位和其他已在房间的玩家的数据
	glog.Errorln("joinRoom", m)
}

func (r *Room) leaveRoom(m *msg.LeaveRoom, o *Occupant) {
	r.removeOccupant(o)
	// todo 下房间其他完广播玩家离开房间消息
	glog.Errorln("leaveRoom", m)
}
func (r *Room) bet(m *msg.LeaveRoom, o *Occupant) {
	glog.Errorln("bet", m)
}
