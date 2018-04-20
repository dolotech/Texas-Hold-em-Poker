package room

type IOccupant interface {
	GetRoom() IRoom
	WriteMsg(msg interface{})
}
type ICreator interface {
	Create(interface{}) IRoom
}
type IRoom interface {
	Cap() uint8
	Len() uint8
	GetNumber() string
	SetNumber(string)
	Data() interface{}
	SetData(interface{})
	Close()
	Closed() chan struct{}
	Send(IOccupant, interface{}) error
	WriteMsg(interface{}, ...uint32)
	Regist(interface{}, interface{})

	New(interface{})IRoom

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}
