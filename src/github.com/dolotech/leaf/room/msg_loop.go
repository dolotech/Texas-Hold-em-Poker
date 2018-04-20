package room

import (
	"github.com/dolotech/lib/utils"
	"server/protocol"
	"errors"
	"github.com/dolotech/lib/route"
	"github.com/golang/glog"
	"github.com/davecgh/go-spew/spew"
)

type MsgLoop struct {
	closedBroadcastChan chan struct{}
	closeChan           chan IRoom
	msgChan             chan *msgObj
	route.Route
}

func NewMsgLoop() *MsgLoop {
	m := &MsgLoop{
		closeChan:           make(chan IRoom, 1),
		closedBroadcastChan: make(chan struct{}),
		msgChan:             make(chan *msgObj, 128),
	}
	go m.msgLoop()

	return m
}

func (r *MsgLoop) Closed() chan struct{} {
	return r.closedBroadcastChan
}
func (r *MsgLoop) msgLoop() {
	defer func() {
		if err := utils.PrintPanicStack(); err != nil {
			go r.msgLoop()
		}
	}()
	for {
		select {
		case m := <-r.closeChan:
			close(r.closedBroadcastChan)
			DelRoom(m.(IRoom))
			return
		case m := <-r.msgChan:
			r.Emit(m.msg, m.o)
		}
	}
}

func (r *MsgLoop) Close(m IRoom) {
	select {
	case r.closeChan <- m:
	default:
	}
}

type msgObj struct {
	msg interface{}
	o   IOccupant
}

func (r *MsgLoop) Send(o IOccupant, m interface{}) error {
	select {
	case r.msgChan <- &msgObj{m, o}:
	default:
		o.WriteMsg(protocol.MSG_ROOM_CLOSED)
	}

	return errors.New("room closed")
}

type Log struct {
	room IRoom
}

func NewLog(room IRoom) *Log {
	return &Log{room: room}
}

func (r *Log) Info(args ...interface{}) {
	glog.InfoDepth(1, r.parseLog(args)...)
}

func (r *Log) Infof(format string, args ...interface{}) {
	glog.InfofDepth(1, format, r.parseLog(args)...)
}

func (r *Log) Error(args ...interface{}) {
	glog.ErrorDepth(1, r.parseLog(args)...)
}
func (r *Log) Debug(args ...interface{}) {
	for k, v := range args {
		args[k] = spew.Sdump(v)
	}
	glog.InfoDepth(1, r.parseLog(args)...)
}

func (r *Log) Debugf(format string, args ...interface{}) {
	for k, v := range args {
		args[k] = spew.Sdump(v)
	}
	glog.InfofDepth(1, format, r.parseLog(args)...)
}
func (r *Log) Errorf(format string, args ...interface{}) {
	glog.ErrorfDepth(1, format, r.parseLog(args)...)
}

func (r *Log) parseLog(args ...interface{}) []interface{} {
	param := make([]interface{}, len(args)+4)
	param[0] = r.room.GetNumber()
	param[1] = r.room.Cap()
	param[2] = r.room.Len()
	param[3] = r.room.Data()
	copy(param[4:], args)
	return param
}
