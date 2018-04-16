package internal

import (
	"server/model"
	"github.com/dolotech/leaf/gate"
	"server/algorithm"
)

type Occupant struct {
	*model.User
	gate.Agent
	room   *Room
	cards  algorithm.Cards
	status int32 // 1为离线状态
}

const (
	Occupant_status_Offline int32 = 1
	Occupant_status_Online  int32 = 0
)

func (o *Occupant)WriteMsg(msg interface{}){
	if o.status != Occupant_status_Offline{
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

func NewOccupant(data *model.User, conn gate.Agent) *Occupant {
	o := &Occupant{
		User:  data,
		Agent: conn,
	}
	return o
}
