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
		msgChan:   make(chan interface{}),
		occupants: make([]*Occupant, data.N, data.Max),
		route:     utils.NewRoute(),
	}

	r.route.Regist(&msg.JoinRoom{},r.joinRoom)
	r.route.Regist(&msg.LeaveRoom{},r.leaveRoom)
	skeleton.Go(r.msgLoop, nil)
	return r
}

type Room struct {
	data      *model.Room
	occupants []*Occupant
	closeChan chan struct{}
	msgChan   chan interface{}
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
			r.route.Route(m)
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

func (r *Room) Close() {
	if atomic.LoadInt32(&r.state) != model.RoomStatus_Closed {
		r.closeChan <- struct{}{}
	}
}
func (r *Room) RecvMsg(msg interface{}) {
	if atomic.LoadInt32(&r.state) != model.RoomStatus_Closed {
		r.msgChan <- msg
	}
}
