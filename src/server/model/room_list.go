package model

import (
	"sync"
	"time"
	"strconv"
	"math/rand"
)

var rooms = func() *roomlist {
	return &roomlist{
		M: make(map[string]IRoom, 1000),
	}
}()

type roomlist struct {
	M map[string]IRoom
	sync.RWMutex
}

func FindRoom() IRoom {
	rooms.Lock()
	for _, v := range rooms.M {
		if v.Cap()< v.MaxCap() {
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

	id := createNumber()
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

func createNumber() string {
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
