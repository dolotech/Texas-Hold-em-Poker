package internal

import (
	"testing"
	"server/model"
	"time"
	msg2 "server/msg"
	"reflect"
)

func TestRoom_Value(t *testing.T) {
	t.Log(reflect.ValueOf(12))
}
func TestRoom_RecvMsg(t *testing.T) {
	/*room:= NewRoom(&model.Room{})


	msg:= &msg2.JoinRoom{RoomNumber:"9999"}


	room.Send(12,msg)

*/
	//msg1:= &msg2.LeaveRoom{RoomNumber:"9999"}
	//room.RecvMsg(12,msg1)

	time.Sleep(time.Second*2)
}
func BenchmarkCloseRoom(t *testing.B) {


/*

	for i:=0;i<t.N;i++{
		room:= NewRoom(&model.Room{})



		go room.Close()
		go room.SendMsg(111)
	}

	t.Log("adfasdfads")

	<- time.After(time.Minute)*/

	//room.CloseChan <- struct{}{}
	//t.Log(cap(room.CloseChan))
}
