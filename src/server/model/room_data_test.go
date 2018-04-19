package model

import (
	"testing"
	"github.com/dolotech/lib/db"
)

func init() {
	db.Init("postgres://postgres:haosql@127.0.0.1:5432/postgres?sslmode=disable")
}

func TestUser_UpdateChips(t *testing.T) {
	room := &Room{
	}

	t.Log(room.Insert())

	room = &Room{Rid: 5}


	id,err:= room.GetById()

	t.Log(room.CreatedAt)
	t.Logf("%v %v %#+v",id,err, room)
}
