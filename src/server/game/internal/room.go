package internal

import (
	"server/model"
	"sync/atomic"
	"github.com/golang/glog"
	"runtime/debug"
	"server/msg"
	"github.com/dolotech/lib/route"
	"github.com/dolotech/leaf/gate"
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
		msgChan:   make(chan *msgObj, 64),
		occupants: make([]*Occupant, data.N, data.Max),
	}

	r.Regist(&msg.JoinRoom{}, r.joinRoom)
	r.Regist(&msg.LeaveRoom{}, r.leaveRoom)
	r.Regist(&msg.Bet{}, r.bet)
	go r.msgLoop()
	return r
}

type Room struct {
	*model.Room
	occupants []*Occupant
	closeChan chan struct{}
	msgChan   chan *msgObj
	State     int32
	route.Route
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
			atomic.StoreInt32(&r.State, RoomStatus_Closed)
			return
		case m := <-r.msgChan:
			o := r.occupant(m.agent.UserData().(*model.User).Uid)
			if o != nil {
				r.Emit(m.msg, o)
			} else {
				r.Emit(m.msg, m.agent)
			}
		}
	}
}
func (r *Room) Write(msg interface{}, exc ...uint32) {
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
	if atomic.LoadInt32(&r.State) != RoomStatus_Closed {
		r.closeChan <- struct{}{}
	}
}

type msgObj struct {
	msg   interface{}
	agent gate.Agent
}

func (r *Room) Send(a gate.Agent, msg interface{}) {
	au := a.UserData()
	if au == nil {
		glog.Errorln("agent UserData is nil")
		return
	}
	if u, ok := au.(*model.User); !ok {
		glog.Errorln("agent UserData type error")
		return
	} else {
		if u.Uid == 0 {
			glog.Errorln("agent UserData Uid ==0")
			return
		}
	}
	if atomic.LoadInt32(&r.State) != RoomStatus_Closed {
		r.msgChan <- &msgObj{msg, a}
	}
}
