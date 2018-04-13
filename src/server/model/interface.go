package model


type IRoom interface {
	SendMsg(msg interface{})
	GetData() *Room
	Close()
	SetData(*Room)
	RecvMsg(uid uint32,msg interface{})
}

