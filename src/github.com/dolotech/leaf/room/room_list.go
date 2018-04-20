package room

import (
	"sync"
	"time"
	"strconv"
	"math/rand"
	"github.com/golang/glog"
	"server/protocol"
	"github.com/dolotech/leaf/gate"
	"github.com/dolotech/lib/utils"
)


func OnMessage(m interface{}, a gate.Agent) {
	defer utils.PrintPanicStack()
	o := a.UserData().(IOccupant)
	if o.GetRoom() != nil {
		o.GetRoom().Send(o, m)
	} else {
		if r := hand.Create(m); r == nil {
			a.WriteMsg(protocol.MSG_NOT_IN_ROOM)
		} else {
			SetRoom(r)
			r.Send(o, m)
		}
	}
	glog.Errorln(m, o)
}

func Regist(r IRoom, m interface{}, h interface{}) {
	r.Regist(m, h)
}

var rooms *roomlist

var hand ICreator

func init() {
	rooms = &roomlist{
		M: make(map[string]IRoom, 1000),
	}
}

func Init(h ICreator) {
	hand = h
}

type roomlist struct {
	M map[string]IRoom
	sync.RWMutex
}

func FindRoom() IRoom {
	rooms.Lock()
	for _, v := range rooms.M {
		if v.Len() < v.Cap() {
			return v
		}
	}
	rooms.Unlock()
	return nil
}

func GetRoom(rid string) IRoom {
	rooms.RLock()
	r := rooms.M[rid]
	rooms.RUnlock()
	return r
}

func SetRoom(room IRoom) string {

	rooms.Lock()
	id := rooms.createNumber()
	room.SetNumber(id)
	rooms.M[id] = room
	rooms.Unlock()
	return id
}
func DelRoom(room IRoom) {
	rooms.Lock()
	delete(rooms.M, room.GetNumber())
	rooms.Unlock()
}

func (this *roomlist) createNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var n string
	for i := 0; i < 100; i++ {
		n = strconv.Itoa(int(r.Int31n(999999-100000) + 100000))
		if _, ok := rooms.M[n]; !ok {
			return n
		}
	}
	return n
}

func Each(f func(o IRoom) bool) {
	rooms.RLock()
	for _, v := range rooms.M {
		if !f(v) {
			break
		}
	}
	rooms.RUnlock()
}
func GetRooms() []IRoom {
	r := make([]IRoom, len(rooms.M))
	rooms.RLock()
	var n = 0
	for _, v := range rooms.M {
		r[n] = v
		n ++
	}

	rooms.RUnlock()
	return r
}
