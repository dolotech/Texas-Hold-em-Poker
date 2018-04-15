package internal

import (
	"server/model"
	"github.com/dolotech/leaf/gate"
	"server/algorithm"
	"sync/atomic"
)

type Occupant struct {
	*model.User
	gate.Agent
	room  *Room
	cards algorithm.Cards
	//Actions chan *ws.Message
	//timer *time.Timer // action timer

	Status int32 // 1为离线状态
}

const (
	Occupant_status_Offline int32 = 1
	Occupant_status_Online  int32 = 0
)

func (o *Occupant) SetData(d interface{}) {
	o.User = d.(*model.User)
}
func (o *Occupant) GetId()uint32 {
	return  o.Uid
}
func (o *Occupant) Online() {
	atomic.StoreInt32(&o.Status, Occupant_status_Online)
}
func (o *Occupant) Offline() {
	atomic.StoreInt32(&o.Status, Occupant_status_Offline)
}

func NewOccupant(data *model.User, conn gate.Agent) *Occupant {
	o := &Occupant{
		User:  data,
		Agent: conn,
	}
	return o
}



