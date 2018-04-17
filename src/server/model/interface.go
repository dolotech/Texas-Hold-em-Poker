package model

import (
	"github.com/dolotech/leaf/gate"
)

type IRoom interface {
	ID() uint32
	Cap() uint8
	Len() uint8
	GetDragin() uint32
	CreatedTime() uint32
	GetNumber() string
	SetNumber(string)
	Close()
	Send(gate.Agent, interface{}) error
	WriteMsg(interface{}, ...uint32)
}
