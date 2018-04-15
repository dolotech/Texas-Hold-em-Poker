package internal

import (
	"server/msg"
	"github.com/golang/glog"
	"sync/atomic"
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
	// todo 掉线重连处理
	if o.room != nil {
		for k, v := range r.occupants {
			if v.Uid == o.Uid {
				r.occupants[k] = o
				if o != v {
					v.Close()
					glog.Infoln("掉线重连处理")
				} else {
					glog.Infoln("同一个链接重复请求加入房间")
				}

				r.WriteMsg(&msg.JoinRoom{Uid: o.Uid}, o.Uid)
				return
			}
		}
	}

	r.addOccupant(o)
	rinfo := &msg.RoomInfo{
		Number: r.Number,
	}

	o.RoomID = r.Number
	o.UpdateRoomId()
	o.room = r

	o.WriteMsg(rinfo)

	r.WriteMsg(&msg.JoinRoom{Uid: o.Uid}, o.Uid)
	glog.Errorln("joinRoom", m)
}

func (r *Room) leaveRoom(m *msg.LeaveRoom, o *Occupant) {
	if atomic.LoadInt32(&r.state) == RoomStatus_Started {
		o.Offline()
		return
	}

	r.removeOccupant(o)
	o.RoomID = ""
	o.room = nil
	o.UpdateRoomId()
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
