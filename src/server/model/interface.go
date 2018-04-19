package model

type IOccupant interface {
	GetRoom() IRoom
	SetRoom(IRoom)
	WriteMsg(msg interface{})
	GetUid() uint32
	IsGameing() bool
	GetPos() uint8
	SetPos(uint8)
}
type IHandler interface {
	NewRoom() IRoom
	NoRoomHandler(m interface{}) IRoom
}
type IRoom interface {
	ID() uint32
	Cap() uint8
	Len() uint8
	GetDragin() uint32
	CreatedTime() uint32
	GetNumber() string
	SetNumber(string)
	Close()
	Closed() chan struct{}
	Send(IOccupant, interface{}) error
	WriteMsg(interface{}, ...uint32)
	Regist(interface{}, interface{})
}
