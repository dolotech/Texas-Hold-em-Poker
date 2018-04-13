package internal

import (
	"server/model"
	"sync/atomic"
	"github.com/golang/glog"
	"runtime/debug"
	"server/msg"
	"github.com/dolotech/lib/utils"
)

func NewRoom(data *model.Room) model.IRoom {
	r := &Room{
		data:      data,
		closeChan: make(chan struct{}),
		msgChan:   make(chan *MsgObj),
		occupants: make([]*Occupant, data.N, data.Max),
		route:     utils.NewRoute(),
	}

	r.route.Regist(&msg.JoinRoom{}, r.joinRoom)
	r.route.Regist(&msg.LeaveRoom{}, r.leaveRoom)
	skeleton.Go(r.msgLoop, nil)
	return r
}

type Room struct {
	data      *model.Room
	occupants []*Occupant
	closeChan chan struct{}
	msgChan   chan *MsgObj
	state     int32
	route     *utils.Route
}

func (r *Room) msgLoop() {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorf("roomid %v err: %v", r.data.Number, err)
			glog.Error(string(debug.Stack()))
			skeleton.Go(r.msgLoop, nil)
		}
	}()
	for {
		select {
		case <-r.closeChan:
			atomic.StoreInt32(&r.state, model.RoomStatus_Closed)
			return
		case m := <-r.msgChan:
			r.route.Route(m.Msg,m.Agre)
		}
	}
}
func (r *Room) SendMsg(msg interface{}) {
	for _, v := range r.occupants {
		if v != nil {
			v.WriteMsg(msg)
		}
	}
}

func (r *Room) GetData() *model.Room {
	return r.data
}
func (r *Room) SetData(data *model.Room) {
	r.data = data
}

func (r *Room) getOccupant(uid uint32) *Occupant {
	for _, v := range r.occupants {
		if v != nil {
			if uid == v.data.Uid {
				return v
			}
		}
	}
	return nil
}
func (r *Room) Close() {
	if atomic.LoadInt32(&r.state) != model.RoomStatus_Closed {
		r.closeChan <- struct{}{}
	}
}

type MsgObj struct {
	Msg  interface{}
	Agre interface{}
}

func (r *Room) RecvMsg(uid uint32, msg interface{}) {
	if uid == 0{
		glog.Errorln("uid is zero")
		return
	}
	if atomic.LoadInt32(&r.state) != model.RoomStatus_Closed {
		o := r.getOccupant(uid)
		if o != nil {
			r.msgChan <- &MsgObj{ msg,o}
		} else {
			r.msgChan <- &MsgObj{ msg,uid}
		}
	}
}
