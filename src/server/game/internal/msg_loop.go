package internal

import (
	"github.com/golang/glog"
	"github.com/davecgh/go-spew/spew"
	"github.com/dolotech/lib/utils"
	"server/protocol"
	"errors"
	"github.com/dolotech/leaf/room"
)

func (r *Room) Closed() chan struct{} {
	return r.closedBroadcastChan
}
func (r *Room) msgLoop() {
	defer func() {
		if err := utils.PrintPanicStack(); err != nil {
			go r.msgLoop()
		}
	}()
	for {
		select {
		case <-r.closeChan:
			close(r.closedBroadcastChan)
			room.DelRoom(r)
			return
		case m := <-r.msgChan:
			r.Emit(m.msg, m.o)
		}
	}
}

func (r *Room) Close() {
	select {
	case r.closeChan <- struct{}{}:
	default:
	}
}

type msgObj struct {
	msg interface{}
	o   room.IOccupant
}

func (r *Room) Send(o room.IOccupant, m interface{}) error {
	select {
	case r.msgChan <- &msgObj{m, o}:
	default:
		o.WriteMsg(protocol.MSG_ROOM_CLOSED)
	}

	return errors.New("room closed")
}

func (r *Room) Info(args ...interface{}) {
	glog.InfoDepth(1, r.parseLog(args)...)
}

func (r *Room) Infof(format string, args ...interface{}) {
	glog.InfofDepth(1, format, r.parseLog(args)...)
}

func (r *Room) Error(args ...interface{}) {
	glog.ErrorDepth(1, r.parseLog(args)...)
}
func (r *Room) Debug(args ...interface{}) {
	for k, v := range args {
		args[k] = spew.Sdump(v)
	}
	glog.InfoDepth(1, r.parseLog(args)...)
}

func (r *Room) Debugf(format string, args ...interface{}) {
	for k, v := range args {
		args[k] = spew.Sdump(v)
	}
	glog.InfofDepth(1, format, r.parseLog(args)...)
}
func (r *Room) Errorf(format string, args ...interface{}) {
	glog.ErrorfDepth(1, format, r.parseLog(args)...)
}

func (r *Room) parseLog(args ...interface{}) []interface{} {
	param := make([]interface{}, len(args)+4)
	param[0] = r.GetNumber()
	param[1] = r.Cap()
	param[2] = r.Len()
	param[3] = r.Data()
	copy(param[4:], args)
	return param
}
