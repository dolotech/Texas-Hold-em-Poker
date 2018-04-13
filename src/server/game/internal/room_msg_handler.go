package internal

import (
	"server/msg"
	"github.com/golang/glog"
)

func (r *Room) joinRoom(m *msg.JoinRoom) {
	glog.Errorln("joinRoom",m)
}


func (r *Room) leaveRoom(m *msg.LeaveRoom) {
	glog.Errorln("leaveRoom",m)
}
