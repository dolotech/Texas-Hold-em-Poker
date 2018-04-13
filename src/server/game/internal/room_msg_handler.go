package internal

import (
	"server/msg"
	"github.com/golang/glog"
)

func (r *Room) joinRoom(m *msg.JoinRoom,uid uint32) {
	glog.Errorln("joinRoom",m)
}


func (r *Room) leaveRoom(m *msg.LeaveRoom,o *Occupant) {
	glog.Errorln("leaveRoom",m)
}
func (r *Room) bet(m *msg.LeaveRoom,o *Occupant) {
	glog.Errorln("leaveRoom",m)
}
