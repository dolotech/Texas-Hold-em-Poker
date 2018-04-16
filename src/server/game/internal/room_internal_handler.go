package internal

import (
	"server/msg"
	"github.com/golang/glog"
	"sync/atomic"
)

func (r *Room) joinRoom(m *msg.JoinRoom, o *Occupant) {
	if o.room != nil {
		for k, v := range r.occupants {
			if v.Uid == o.Uid {
				// todo 掉线重连现场数据替换处理
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
	glog.Errorln(o)
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

func (r *Room) sitDown(m *msg.SitDown, o *Occupant) {
	glog.Errorln("sitDown", m)
}

func (r *Room) standUp(m *msg.StandUp, o *Occupant) {
	glog.Errorln("standUp", m)
}
