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
	"server/algorithm"
	"time"
)

const (
	RoomStatus_Closed  int32 = 9
	RoomStatus_Started int32 = 1
	RoomStatus_End     int32 = 2
	RoomStatus_Ready   int32 = 0
)

type Room struct {
	*model.Room
	route.Route

	occupants []*Occupant
	observes  []*Occupant // 旁观列表

	closeChan chan struct{}
	msgChan   chan *msgObj
	state     int32

	remain int
	allin  int
	n      uint8

	SB       uint32          // 小盲注
	BB       uint32          // 大盲注
	Cards    algorithm.Cards //公共牌
	Pot      []uint32        // 当前奖池筹码数
	Timeout  time.Duration   // 倒计时超时时间(秒)
	Button   uint8           // 当前庄家座位号，从1开始
	Chips    []uint32        // 玩家本局下注的总筹码数，与occupants一一对应
	Bet      uint32          // 当前下注额
	Max      uint8           // 房间最大玩家人数
	MaxChips uint32
	MinChips uint32
}

func NewRoom(max uint8, sb, bb uint32) model.IRoom {
	if max <= 0 || max > 9 {
		max = 9 // default 9 occupants
	}

	r := &Room{
		Room:      &model.Room{},
		closeChan: make(chan struct{}),
		msgChan:   make(chan *msgObj, 64),
		occupants: make([]*Occupant, max),
		Chips:     make([]uint32, max),
		Pot:       make([]uint32, 1),
		Timeout:   time.Second * 10,
		SB:        sb,
		BB:        bb,
		Max:       max,
	}

	r.Regist(&msg.JoinRoom{}, r.joinRoom)
	r.Regist(&msg.LeaveRoom{}, r.leaveRoom)
	r.Regist(&msg.Bet{}, r.bet)
	r.Regist(&msg.Fold{}, r.fold)
	r.Regist(&msg.SitDown{}, r.sitDown) //
	r.Regist(&msg.StandUp{}, r.standUp) //
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

func (r *Room) addOccupant(o *Occupant) uint8 {
	for _, v := range r.occupants {
		if v != nil && v.Uid == o.Uid {
			return 0
		}
	}

	for k, v := range r.occupants {
		if v == nil {
			r.occupants[k] = o
			r.n ++
			o.Pos = uint8(k + 1)
			o.SetSitdown()
			return o.Pos
		}
	}
	return 0
}

func (r *Room) removeOccupant(o *Occupant) {
	for k, v := range r.occupants {
		if v != nil && v.Uid == o.Uid {
			v.Pos = 0
			r.occupants[k] = nil
			r.n --
			return
		}
	}
}

func (r *Room) addObserve(o *Occupant) uint8 {
	for _, v := range r.observes {
		if v != nil && v.Uid == o.Uid {
			return 0
		}
	}
	o.SetObserve()
	r.observes = append(r.observes, o)

	return 0
}

func (r *Room) removeObserve(o *Occupant) {
	for k, v := range r.observes {
		if v != nil && v.Uid == o.Uid {
			r.observes = append(r.observes[:k], r.observes[k+1:]...)
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

// start starts from 0
func (r *Room) Each(start uint8, f func(o *Occupant) bool) {
	volume := r.Cap()
	end := (volume + start - 1) % volume
	i := start
	for ; i != end; i = (i + 1) % volume {
		if r.occupants[i] != nil && r.occupants[i].IsGameing() && !f(r.occupants[i]) {
			return
		}
	}

	// end
	if r.occupants[i] != nil && r.occupants[i].IsGameing() {
		f(r.occupants[i])
	}
}
func (r *Room) Cap() uint8 {
	return uint8(len(r.occupants))
}
func (r *Room) MaxCap() uint8 {
	return r.Max
}
func (r *Room) GetNumber() string {
	return r.Number
}
func (r *Room) SetNumber(value string) {
	r.Number = value
}
