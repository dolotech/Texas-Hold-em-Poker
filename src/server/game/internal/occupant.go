package internal

import (
	"server/model"
	"github.com/dolotech/leaf/gate"
	"server/algorithm"
	"time"
)

// todo 下注超时

type Occupant struct {
	*model.User
	gate.Agent
	room   *Room
	cards  algorithm.Cards
	Pos    uint8 // 玩家座位号，从1开始
	status int32 // 1为离线状态

	Bet        uint32 // 当前下注
	actions    chan int32
	actionName string
}

const (
	Occupant_status_InGame  int32 = 3
	Occupant_status_Offline int32 = 1
	Occupant_status_Observe int32 = 2
	Occupant_status_Sitdown int32 = 0
)

func (o *Occupant) SetAction(n int32) {
	if o.actionName == model.BET_ {
		o.actions <- n
	}
}
func (o *Occupant) GetAction(timeout time.Duration) int32 {
	timer := time.NewTimer(timeout)
	o.actionName = model.BET_
	select {
	case n := <-o.actions:
		timer.Stop()
		o.actionName = ""
		return n
	case <-timer.C:
		o.actionName = ""
		timer.Stop()
		return -1 // 超时弃牌
	}
}

func (o *Occupant) WriteMsg(msg interface{}) {
	if o.status != Occupant_status_Offline {
		o.Agent.WriteMsg(msg)
	}
}
func (o *Occupant) SetData(d interface{}) {
	o.User = d.(*model.User)
}
func (o *Occupant) GetId() uint32 {
	return o.Uid
}

func (o *Occupant) SetObserve() {
	o.status = Occupant_status_Observe
}

func (o *Occupant) IsObserve() bool {
	return o.status == Occupant_status_Observe
}

func (o *Occupant) SetOffline() {
	o.status = Occupant_status_Offline
}

func (o *Occupant) IsOffline() bool {
	return o.status == Occupant_status_Offline
}

func (o *Occupant) SetSitdown() {
	o.status = Occupant_status_Sitdown
}

func (o *Occupant) IsSitdown() bool {
	return o.status == Occupant_status_Sitdown
}

func (o *Occupant) SetGameing() {
	o.status = Occupant_status_InGame
}

func (o *Occupant) IsGameing() bool {
	return o.status == Occupant_status_InGame
}

func (o *Occupant) Replace(value *Occupant) {
	o.Pos = value.Pos
	o.cards = value.cards
	o.room = value.room
}

func NewOccupant(data *model.User, conn gate.Agent) *Occupant {
	o := &Occupant{
		User:    data,
		Agent:   conn,
		actions: make(chan int32),
	}
	return o
}
