package internal

import (
	"server/model"
	"sync/atomic"
	"github.com/golang/glog"
	"runtime/debug"
	"server/msg"
	"github.com/dolotech/lib/route"
)

const (
	RoomStatus_Closed  int32 = 9
	RoomStatus_Started int32 = 1
	RoomStatus_End     int32 = 2
	RoomStatus_Ready   int32 = 0
)

func NewRoom(data *model.Room) model.IRoom {
	r := &Room{
		Room:      data,
		closeChan: make(chan struct{}),
		msgChan: make(chan *msgObj,64),
		occupants: make([]*Occupant, data.N, data.Max),
		Route:     route.NewRoute(),
	}

	r.Regist(&msg.JoinRoom{}, r.joinRoom)
	r.Regist(&msg.LeaveRoom{}, r.leaveRoom)
	r.Regist(&msg.Bet{}, r.bet)
	skeleton.Go(r.msgLoop, nil)
	return r
}

type Room struct {
	*model.Room
	occupants []*Occupant
	closeChan chan struct{}
	msgChan   chan *msgObj
	state     int32
	*route.Route
}

func (r *Room) msgLoop() {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorf("roomid %v err: %v", r.Room.Number, err)
			glog.Error(string(debug.Stack()))
			skeleton.Go(r.msgLoop, nil)
		}
	}()
	for {
		select {
		case <-r.closeChan:
			atomic.StoreInt32(&r.state, RoomStatus_Closed)
			return
		case m := <-r.msgChan:
			o := r.occupant(m.uid)
			if o != nil {
				r.Emit(m.msg, o)
			} else {
				r.Emit(m.msg, m.uid)
			}
		}
	}
}
func (r *Room) Write(msg interface{}) {
	for _, v := range r.occupants {
		if v != nil {
			v.WriteMsg(msg)
		}
	}
}

func (r *Room) Data() *model.Room {
	return r.Room
}

func (r *Room) occupant(uid uint32) *Occupant {
	for _, v := range r.occupants {
		if v != nil {
			if uid == v.User.Uid {
				return v
			}
		}
	}
	return nil
}
func (r *Room) Close() {
	if atomic.LoadInt32(&r.state) != RoomStatus_Closed {
		r.closeChan <- struct{}{}
	}
}

type msgObj struct {
	msg interface{}
	uid uint32
}

func (r *Room) Send(uid uint32, msg interface{}) {
	if uid == 0 {
		glog.Errorln("Uid is zero")
		return
	}
	if atomic.LoadInt32(&r.state) != RoomStatus_Closed {
		r.msgChan <- &msgObj{msg, uid}
	}
}
