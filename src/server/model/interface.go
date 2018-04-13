package model

type IRoom interface {
	Data() *Room
	Close()
	Send(uint32, interface{})
	Write(interface{})
}
