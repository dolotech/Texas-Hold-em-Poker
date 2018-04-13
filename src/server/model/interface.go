package model


type IRoom interface {
	SendMsg(msg interface{})
	GetData() *Room
	Close()
	SetData(*Room)
	RecvMsg(msg interface{})
}

