package room

import (
	"github.com/golang/glog"
	"server/protocol"
	"github.com/dolotech/lib/route"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/dolotech/lib/utils"
)

func NewRoom() *BaseRoom {
	r := &BaseRoom{
		closeChan:           make(chan struct{}),
		closedBroadcastChan: make(chan struct{}),
		msgChan:             make(chan *msgObj, 128),
	}
	go r.msgLoop()
	return r
}

type BaseRoom struct {
	route.Route

	closedBroadcastChan chan struct{}
	closeChan           chan struct{}
	msgChan             chan *msgObj
}

func (r *BaseRoom) Closed() chan struct{} {
	return r.closedBroadcastChan
}
func (r *BaseRoom) msgLoop() {
	defer func() {
		if err := utils.PrintPanicStack(); err != nil {
			go r.msgLoop()
		}
	}()
	for {
		select {
		case <-r.closeChan:
			close(r.closedBroadcastChan)
			DelRoom(r)
			return
		case m := <-r.msgChan:
			r.Emit(m.msg, m.o)
		}
	}
}

func (r *BaseRoom) Close() {
	select {
	case r.closeChan <- struct{}{}:
	default:
	}
}

type msgObj struct {
	msg interface{}
	o   IOccupant
}

func (r *BaseRoom) Send(o IOccupant, m interface{}) error {
	select {
	case r.msgChan <- &msgObj{m, o}:
	default:
		o.WriteMsg(protocol.MSG_ROOM_CLOSED)
	}

	return errors.New("room closed")
}

func (r *BaseRoom) Cap() uint8          { return 0 }
func (r *BaseRoom) Len() uint8          { return 0 }
func (r *BaseRoom) Data() interface{}   { return nil }
func (r *BaseRoom) SetData(interface{}) {}

func (r *BaseRoom) GetNumber() string               { return "" }
func (r *BaseRoom) SetNumber(string)                {}
func (r *BaseRoom) WriteMsg(interface{}, ...uint32) {}
func (r *BaseRoom) New(interface{}) IRoom           { return nil }

func (r *BaseRoom) Info(args ...interface{}) {
	glog.InfoDepth(1, r.parseLog(args)...)
}

func (r *BaseRoom) Infof(format string, args ...interface{}) {
	glog.InfofDepth(1, format, r.parseLog(args)...)
}

func (r *BaseRoom) Error(args ...interface{}) {
	glog.ErrorDepth(1, r.parseLog(args)...)
}
func (r *BaseRoom) Debug(args ...interface{}) {
	for k, v := range args {
		args[k] = spew.Sdump(v)
	}
	glog.InfoDepth(1, r.parseLog(args)...)
}

func (r *BaseRoom) Debugf(format string, args ...interface{}) {
	for k, v := range args {
		args[k] = spew.Sdump(v)
	}
	glog.InfofDepth(1, format, r.parseLog(args)...)
}
func (r *BaseRoom) Errorf(format string, args ...interface{}) {
	glog.ErrorfDepth(1, format, r.parseLog(args)...)
}

func (r *BaseRoom) parseLog(args ...interface{}) []interface{} {
	param := make([]interface{}, len(args)+4)
	param[0] = r.GetNumber()
	param[1] = r.Cap()
	param[2] = r.Len()
	param[3] = r.Data()
	copy(param[4:], args)
	return param
}
