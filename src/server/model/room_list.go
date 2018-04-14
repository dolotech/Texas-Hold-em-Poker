package model

import (
	"sync"
	"time"
	"strconv"
	"math/rand"
)

var rooms = NewRoomList()

type roomlist struct {
	M map[string]IRoom
	sync.RWMutex
}

func NewRoomList() *roomlist {
	return &roomlist{
		M: make(map[string]IRoom, 1000),
	}
}

func GetRoom(rid string) IRoom {
	rooms.RLock()
	r := rooms.M[rid]
	rooms.RUnlock()
	return r
}

func SetRoom(room IRoom) string{
	rooms.Lock()

	id := createNumber()
	room.Data().Number = id
	rooms.M[id] = room
	rooms.Unlock()
	return  id
}
func DelRoom(room IRoom) {
	rooms.Lock()
	delete(rooms.M, room.Data().Number)
	room.Data().Rid = 0
	rooms.Unlock()
}

func createNumber() string {
	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	var n string
	for i:=0;i<100;i++{
		n = strconv.Itoa(int(r.Int31n(999999-100000) + 100000))
		if _,ok:=rooms.M[n];!ok{
			return  n
		}
	}
	return n
}
