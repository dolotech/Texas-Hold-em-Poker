package internal

import (
	"time"
	"server/model"
	"github.com/dolotech/leaf/gate"
)

type Occupant struct {
	*model.User
	gate.Agent
	room *Room
	//Actions chan *ws.Message
	timer *time.Timer // action timer
}

func NewOccupant(data *model.User, conn gate.Agent) *Occupant {
	o := &Occupant{
		User:  data,
		Agent: conn,
	}
	return o
}
