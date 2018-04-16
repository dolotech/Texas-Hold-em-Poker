package model

import (
	"github.com/dolotech/leaf/gate"
)

type IRoom interface {
	Cap() uint8
	MaxCap() uint8
	//Data() *Room


	GetNumber()string
	SetNumber(string)
	Close()
	Send(gate.Agent, interface{}) error
	WriteMsg(interface{}, ...uint32)
}
