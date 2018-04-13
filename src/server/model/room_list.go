package model

import (
	"sync"
)

var rooms = NewRoomList()

type roomlist struct {
	M map[uint32]IRoom
	sync.RWMutex
}

func NewRoomList() *roomlist {
	return &roomlist{
		M: make(map[uint32]IRoom, 1000),
	}
}

func GetRoom(rid uint32) IRoom {
	rooms.RLock()
	r := rooms.M[rid]
	rooms.RUnlock()
	return r
}

func SetRoom(room IRoom) {
	rooms.Lock()

	id := room.Data().Rid
	room.Data().Rid = id
	rooms.M[id] = room
	rooms.Unlock()
}
func DelRoom(room IRoom) {
	rooms.Lock()
	delete(rooms.M, room.Data().Rid)
	room.Data().Rid = 0
	rooms.Unlock()
}
