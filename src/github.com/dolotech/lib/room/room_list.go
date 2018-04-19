package room

import (
	"sync"
	"time"
	"strconv"
	"math/rand"
	"reflect"
	"github.com/golang/glog"
	"server/model"
	"server/protocol"
	"github.com/dolotech/leaf/gate"
	"github.com/dolotech/leaf/module"
)

var skeleton *module.Skeleton

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func onMessage(m interface{}, a gate.Agent) {
	o := a.UserData().(model.IOccupant)
	if o.GetRoom() != nil {
		o.GetRoom().Send(o, m)
	} else {
		if r := hand.NoRoomHandler(m); r == nil {
			a.WriteMsg(protocol.MSG_NOT_IN_ROOM)
		} else {

			SetRoom(r)
			r.Send(o, m)
		}
	}
	glog.Errorln(m, o)
}

func Regist(r model.IRoom, m interface{}, h interface{}) {
	r.Regist(m, h)
	handler(m, onMessage)
}

var rooms *roomlist

var hand model.IHandler

func init() {
	rooms = &roomlist{
		M: make(map[string]model.IRoom, 1000),
	}
}

func Init(h model.IHandler, s *module.Skeleton) {
	hand = h

	skeleton = s

	h.NewRoom()
}

type roomlist struct {
	M map[string]model.IRoom
	sync.RWMutex
}

func FindRoom() model.IRoom {
	rooms.Lock()
	for _, v := range rooms.M {
		if v.Len() < v.Cap() {
			return v
		}
	}
	rooms.Unlock()
	return nil
}

func GetRoom(rid string) model.IRoom {
	rooms.RLock()
	r := rooms.M[rid]
	rooms.RUnlock()
	return r
}

func SetRoom(room model.IRoom) string {

	rooms.Lock()
	id := rooms.createNumber()
	room.SetNumber(id)
	rooms.M[id] = room
	rooms.Unlock()
	return id
}
func DelRoom(room model.IRoom) {
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

func GetRooms() []model.IRoom {

	r := make([]model.IRoom, len(rooms.M))
	rooms.RLock()
	var n = 0
	for _, v := range rooms.M {
		r[n] = v
		n ++
	}
	rooms.RUnlock()
	return r
}
