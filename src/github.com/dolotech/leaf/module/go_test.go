package module

import (
	"testing"
	"time"
)

func TestSkeleton_Go(t *testing.T) {

	s := Skeleton{}
	s.Init()
	f := func() {
		<-time.After(time.Second)
		t.Log("ffff")
	}

	ch:= make(chan int)
	c := func() {
		<-time.After(time.Second)
		t.Log("ccc")

		ch <- 0
	}

	s.Go(f, c)

	t.Log("end",<-ch)

}
