package model

import "github.com/dolotech/leaf/gate"

type IRoom interface {
	Data() *Room
	Close()
	Send(gate.Agent, interface{})
	WriteMsg(interface{},...uint32)
}
type IOccupant interface {
	//Data() *U
	//Close()
	//Send(gate.Agent, interface{})
	Write(interface{})
}
