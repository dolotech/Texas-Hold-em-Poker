package internal

import (
	"server/msg"
	"github.com/golang/glog"
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

func (r *Room) joinRoom(m *msg.JoinRoom, o *Occupant) {
	r.addOccupant(o)
	rinfo := &msg.RoomInfo{
		Number: r.Number,
	}

	o.RoomID = r.Number
	o.User.UpdateRoomId()

	o.WriteMsg(rinfo)

	r.WriteMsg(&msg.JoinRoom{Uid: o.Uid},o.Uid)
	glog.Errorln("joinRoom", m)
}

func (r *Room) leaveRoom(m *msg.LeaveRoom, o *Occupant) {
	r.removeOccupant(o)
	o.RoomID = ""
	o.room = nil
	o.User.UpdateRoomId()
	leave := &msg.LeaveRoom{
		RoomNumber: r.Number,
		Uid:        o.Uid,
	}
	r.WriteMsg(leave)
	glog.Errorln("leaveRoom", m)
}
func (r *Room) bet(m *msg.Bet, o *Occupant) {



	glog.Errorln("bet", m)
}
