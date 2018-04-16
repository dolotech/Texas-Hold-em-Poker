package internal

import (
	"server/model"
	"github.com/dolotech/leaf/gate"
	"server/algorithm"
	"server/msg"
)

type Occupant struct {
	*model.User
	gate.Agent
	room   *Room
	cards  algorithm.Cards
	Pos    uint8 // 玩家座位号，从1开始
	status int32 // 1为离线状态
}

const (
	Occupant_status_Standup int32 = 3
	Occupant_status_Sitdown int32 = 2
	Occupant_status_Offline int32 = 1
	Occupant_status_Online  int32 = 0
)

func (o *Occupant) WriteMsg(msg interface{}) {
	if o.status != Occupant_status_Offline {
		o.Agent.WriteMsg(msg)
	}
}
func (o *Occupant) SetData(d interface{}) {
	o.User = d.(*model.User)
}
func (o *Occupant) GetId() uint32 {
	return o.Uid
}
func (o *Occupant) Online() {
	o.status = Occupant_status_Online
}
func (o *Occupant) Offline() {
	o.status = Occupant_status_Offline
}

func (o *Occupant) Standup() {
	o.status = Occupant_status_Standup
	o.WriteMsg(&msg.StandUp{})
}
func (o *Occupant) Sitdown() {
	o.status = Occupant_status_Sitdown
}

func (o *Occupant) inGame() bool {
	return o.status == Occupant_status_Sitdown
}

func NewOccupant(data *model.User, conn gate.Agent) *Occupant {
	o := &Occupant{
		User:  data,
		Agent: conn,
	}
	return o
}
