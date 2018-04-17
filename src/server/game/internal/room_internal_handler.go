package internal

import (
	"server/protocol"
	"github.com/golang/glog"
)

func (r *Room) joinRoom(m *protocol.JoinRoom, o *Occupant) {
	if o.room != nil {
		for k, v := range r.occupants {
			if v.Uid == o.Uid {
				// todo 掉线重连现场数据替换处理
				o.Replace(r.occupants[k])
				r.occupants[k] = o

				if o != v {
					v.Close()
					glog.Infoln("掉线重连处理")
				} else {
					glog.Infoln("同一个链接重复请求加入房间")
				}

				r.WriteMsg(&protocol.UserInfo{Uid: o.Uid}, o.Uid)
				return
			}
		}
	}
	glog.Errorln(o)
	rinfo := &protocol.RoomInfo{
		Number: r.Number,
	}
	userinfos := make([]*protocol.UserInfo, 0, r.Cap())
	r.Each(0, func(o *Occupant) bool {
		userinfo := &protocol.UserInfo{
			Nickname: o.Nickname,
			Uid:      o.Uid,
			Account:  o.Account,
			Sex:      o.Sex,
			Profile:  o.Profile,
			Chips:    o.chips,
		}
		userinfos = append(userinfos, userinfo)
		return true
	})

	pos := r.addOccupant(o)

	// 坐下失败转为旁观
	if pos == 0 {
		r.addObserve(o)
	} else {
		userInfo := &protocol.UserInfo{
			Nickname: o.Nickname,
			Uid:      o.Uid,
			Account:  o.Account,
			Sex:      o.Sex,
			Profile:  o.Profile,
			Chips:    o.chips,
		}
		r.Broadcast(&protocol.JoinRoomBroadcast{UserInfo: userInfo}, true, o.Uid)
	}

	o.RoomID = r.Number
	o.UpdateRoomId()
	o.room = r
	o.WriteMsg(&protocol.JoinRoomResp{UserInfos: userinfos, RoomInfo: rinfo})

	glog.Errorln("joinRoom", m)
}

func (r *Room) leaveAndRecycleChips(o *Occupant) {
	if r.removeOccupant(o) > 0 {
		// 玩家站起回收带入筹码
		gap := int32(o.chips) - int32(r.LvChips)
		if gap == 0 {
			o.UpdateChips(gap)
		}
	}
}
func (r *Room) leaveRoom(m *protocol.LeaveRoom, o *Occupant) {


	r.removeObserve(o)
	r.leaveAndRecycleChips(o)

	o.RoomID = ""
	o.room = nil
	o.UpdateRoomId()

	leave := &protocol.LeaveRoom{
		RoomNumber: r.Number,
		Uid:        o.Uid,
	}
	r.Broadcast(leave, true)
	glog.Errorln("leaveRoom", m)
}

func (r *Room) bet(m *protocol.Bet, o *Occupant) {
	if !o.IsGameing() {
		o.WriteMsg(protocol.MSG_NOT_NOT_START)
		return
	}

	if m.Value < 0 {
		o.SetAction(-1)
		o.SetSitdown()
	} else {
		o.SetAction(m.Value)
	}

	glog.Errorln("bet", m)
}

func (r *Room) sitDown(m *protocol.SitDown, o *Occupant) {
	pos := r.addOccupant(o)
	if pos == 0 {
		// 给进入房间的玩家带入筹码
		o.chips = r.LvChips
		r.addObserve(o)
	} else {

	}
	r.Broadcast(&protocol.SitDown{Uid: o.Uid, Pos: o.Pos}, true)

	glog.Errorln("sitDown", m)
}

func (r *Room) standUp(m *protocol.StandUp, o *Occupant) {

	o.SetAction(-1)
	r.leaveAndRecycleChips(o)

	r.addObserve(o)
	r.Broadcast(&protocol.StandUp{Uid: o.Uid}, true)

	glog.Errorln("standUp", m)
}
