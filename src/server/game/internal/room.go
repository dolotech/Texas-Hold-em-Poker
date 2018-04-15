package internal

import (
	"server/model"
	"sync/atomic"
	"github.com/golang/glog"
	"runtime/debug"
	"server/msg"
	"github.com/dolotech/lib/route"
	"github.com/dolotech/leaf/gate"
	"errors"
)

const (
	RoomStatus_Closed  int32 = 9
	RoomStatus_Started int32 = 1
	RoomStatus_End     int32 = 2
	RoomStatus_Ready   int32 = 0
)

type Room struct {
	*model.Room
	occupants []*Occupant
	closeChan chan struct{}
	msgChan   chan *msgObj
	state     int32
	route.Route
}

func NewRoom(data *model.Room) model.IRoom {
	r := &Room{
		Room:      data,
		closeChan: make(chan struct{}),
		msgChan:   make(chan *msgObj, 64),
		occupants: make([]*Occupant, data.N, data.Max),
	}

	r.Regist(&msg.JoinRoom{}, r.joinRoom)
	r.Regist(&msg.LeaveRoom{}, r.leaveRoom)
	r.Regist(&msg.Bet{}, r.bet)
	go r.msgLoop()
	return r
}

func (r *Room) msgLoop() {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorf("roomid %v err: %v", r.Room.Number, err)
			glog.Error(string(debug.Stack()))
			go r.msgLoop()
		}
	}()
	for {
		select {
		case <-r.closeChan:
			atomic.StoreInt32(&r.state, RoomStatus_Closed)
			return
		case m := <-r.msgChan:
			r.Emit(m.msg, m.o)
		}
	}
}
func (r *Room) WriteMsg(msg interface{}, exc ...uint32) {
	for _, v := range r.occupants {
		if v != nil {
			for _, uid := range exc {
				if uid == v.Uid {
					goto End
				}
			}
			v.WriteMsg(msg)
		}
	End:
	}
}

func (r *Room) Data() *model.Room {
	return r.Room
}

func (r *Room) occupant(uid uint32) *Occupant {
	for _, v := range r.occupants {
		if v != nil {
			if uid == v.Uid {
				return v
			}
		}
	}
	return nil
}
func (r *Room) addOccupant(o *Occupant) {
	for _, v := range r.occupants {
		if v.Uid == o.Uid {
			return
		}
	}
	r.occupants = append(r.occupants, o)
}

func (r *Room) removeOccupant(o *Occupant) {
	for k, v := range r.occupants {
		if v.Uid == o.Uid {
			r.occupants = append(r.occupants[:k], r.occupants[k+1:]...)
			return
		}
	}
}

func (r *Room) Close() {
	if atomic.LoadInt32(&r.state) != RoomStatus_Closed {
		r.closeChan <- struct{}{}
	}
}

type msgObj struct {
	msg interface{}
	o   gate.Agent
}

func (r *Room) Send(o gate.Agent, m interface{}) error {
	if atomic.LoadInt32(&r.state) != RoomStatus_Closed {
		r.msgChan <- &msgObj{m, o}
		return nil
	} else {
		o.WriteMsg(msg.MSG_ROOM_CLOSED)
	}
	return errors.New("room closed")
}
