package internal

import (
	"server/model"
	"server/protocol"
	"server/algorithm"
	"time"
	"github.com/dolotech/leaf/room"
	"github.com/golang/glog"
)

const (
	RUNNING  uint8= 1
	GAMEOVER uint8= 0
)

type Room struct {
	*model.Room
	*room.MsgLoop
	*room.Log
	Occupants   []*Occupant
	observes    []*Occupant // 站起的玩家
	AutoSitdown []*Occupant // 自动坐下队列

	remain int
	allin  int
	n      uint8
	status uint8

	SB       uint32          // 小盲注
	BB       uint32          // 大盲注
	Cards    algorithm.Cards // 公共牌
	Pot      []uint32        // 奖池筹码数, 第一项为主池，其他项(若存在)为边池
	Timeout  time.Duration   // 倒计时超时时间(秒)
	Button   uint8           // 当前庄家座位号，从1开始
	Chips    []uint32        // 玩家本局下注的总筹码数，与occupants一一对应
	Bet      uint32          // 当前回合 上一玩家下注额
	Max      uint8           // 房间最大玩家人数
	MaxChips uint32
	MinChips uint32
}

func NewRoom(max uint8, sb, bb uint32, chips uint32, timeout uint8) *Room {
	if max <= 0 || max > 9 {
		max = 9 // default 9 Occupants
	}

	r := &Room{
		Room:      &model.Room{DraginChips: chips,},
		MsgLoop:   room.NewMsgLoop(),
		Chips:     make([]uint32, max),
		Occupants: make([]*Occupant, max),
		Pot:       make([]uint32, 0, max),
		Timeout:   time.Second * time.Duration(timeout),
		SB:        sb,
		BB:        bb,
		Max:       max,
	}

	r.Log = room.NewLog(r)
	r.Regist(&protocol.JoinRoom{}, r.joinRoom)
	r.Regist(&protocol.LeaveRoom{}, r.leaveRoom)
	r.Regist(&protocol.Bet{}, r.bet)
	r.Regist(&protocol.SitDown{}, r.sitDown) //
	r.Regist(&protocol.StandUp{}, r.standUp) //
	r.Regist(&protocol.Chat{}, r.chat)       //
	r.Regist(&protocol.Chat{}, r.chat)       //
	r.Regist(&startDelay{}, r.startDelay)    //

	return r
}

type startDelay struct {
	kind uint8
}

func (r *Room) New(m interface{}) room.IRoom {
	glog.Errorln(r, m)
	if msg, ok := m.(*protocol.JoinRoom); ok {
		if len(msg.RoomNumber) == 0 {
			r := room.FindRoom()
			return r
		}
		r := room.GetRoom(msg.RoomNumber)
		if r != nil {
			return r
		}
		room := NewRoom(9, 5, 10, 1000, model.Timeout)
		room.Insert()
		return room
	}
	return nil
}

func (r *Room) WriteMsg(msg interface{}, exc ...uint32) {
	for _, v := range r.Occupants {
		if v != nil {
			for _, uid := range exc {
				if uid == v.GetUid() {
					goto End
				}
			}
			v.WriteMsg(msg)
		}
	End:
	}
}

func (r *Room) Broadcast(msg interface{}, all bool, exc ...uint32) {
	for _, v := range r.Occupants {
		if v != nil && (all || !v.IsGameing()) {
			for _, uid := range exc {
				if uid == v.GetUid() {
					goto End1
				}
			}
			v.WriteMsg(msg)
		}
	End1:
	}
	for _, v := range r.observes {
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
	for _, v := range r.Occupants {
		if v != nil && v.GetUid() == o.Uid {
			return 0
		}
	}

	for k, v := range r.Occupants {
		if v == nil {
			r.Occupants[k] = o
			o.SetRoom(r)
			o.Pos = uint8(k + 1)
			o.SetSitdown()
			return o.Pos
		}
	}
	return 0
}

func (r *Room) removeOccupant(o *Occupant) uint8 {
	for k, v := range r.Occupants {
		if v != nil && v.GetUid() == o.Uid {
			v.SetPos(0)
			r.Occupants[k] = nil
			return uint8(k + 1)
		}
	}
	return 0
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

// start starts from 0
func (r *Room) Each(start uint8, f func(o *Occupant) bool) {
	volume := r.Cap()
	end := (volume + start - 1) % volume
	i := start
	for ; i != end; i = (i + 1) % volume {
		if r.Occupants[i] != nil && r.Occupants[i].IsGameing() && !f(r.Occupants[i]) {
			return
		}
	}

	// end
	if r.Occupants[i] != nil && r.Occupants[i].IsGameing() {
		f(r.Occupants[i])
	}
}

func (r *Room) Cap() uint8 {
	return r.Max
}
func (r *Room) Len() uint8 {
	var num uint8
	for _, v := range r.Occupants {
		if v != nil {
			num ++
		}
	}
	return num
}

func (r *Room) GetNumber() string {
	return r.Number
}
func (r *Room) SetNumber(value string) {
	r.Number = value
}

func (r *Room) Data() interface{} { return r.Room }
func (r *Room) SetData(d interface{}) {
	r.Room = d.(*model.Room)
}
