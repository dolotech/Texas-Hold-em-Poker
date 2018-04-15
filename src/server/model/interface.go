package model

import (
	"github.com/dolotech/leaf/gate"
)

type IRoom interface {
	Data() *Room
	Close()
	Send(gate.Agent, interface{}) error
	WriteMsg(interface{}, ...uint32)
}
type IOccupant interface {
	Online()
	Offline()
	Write(interface{})
	GetId()uint32
	SetData(d interface{})
}
