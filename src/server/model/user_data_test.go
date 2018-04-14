package model

import (
	"testing"
	"time"
	"math/rand"
)

func TestUser_Insert(t *testing.T) {
	user:= &User{Nickname:"Michael",UnionId:"abcc"}
	t.Log(user.Insert())
}

func BenchmarkUser_Insert(b *testing.B) {
	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0;i<b.N;i++{
		n:=r.Int31n(9999 -1000) +1000
		//n:=time.Now().Format("0102") + "123"
		b.Log(n)
	}
}

func TestUserData_Login(t *testing.T) {

	n:=rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(9999 -1000) +1000
	//n:=time.Now().Format("0102") + "123"
	t.Log(n)
}
func TestSeq(t *testing.T) {

}
