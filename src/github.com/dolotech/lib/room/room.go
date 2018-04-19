package room

import (
	"github.com/golang/glog"
	"runtime/debug"
	"server/protocol"
	"github.com/dolotech/lib/route"
	"errors"
	"server/model"
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
		if err := recover(); err != nil {
			glog.Errorf("roomid %v err: %v", r, err)
			glog.Error(string(debug.Stack()))
			go r.msgLoop()
		}
	}()
	for {
		select {
		case <-r.closeChan:
			close(r.closedBroadcastChan)
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
	o   model.IOccupant
}

func (r *BaseRoom) Send(o model.IOccupant, m interface{}) error {
	select {
	case r.msgChan <- &msgObj{m, o}:
	default:
		o.WriteMsg(protocol.MSG_ROOM_CLOSED)
	}
	return errors.New("room closed")
}
