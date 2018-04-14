package model

import "github.com/dolotech/leaf/gate"

type IRoom interface {
	Data() *Room
	Close()
	Send(gate.Agent, interface{})
	Write(interface{},...uint32)
}
