package mongodb

import (
	"gopkg.in/mgo.v2"
	"testing"
)

func TestDial(t *testing.T) {
	c, err := Dial("localhost", 10)
	if err != nil {
		t.Log(err)
		return
	}
	defer c.Close()

	// session
	s := c.Ref()
	defer c.UnRef(s)
	err = s.DB("test").C("counters").RemoveId("test")
	if err != nil && err != mgo.ErrNotFound {
		t.Log(err)
		return
	}

	// auto increment
	err = c.EnsureCounter("test", "counters", "test")
	if err != nil {
		t.Log(err)
		return
	}
	for i := 0; i < 3; i++ {
		id, err := c.NextSeq("test", "counters", "test")
		if err != nil {
			t.Log(err)
			return
		}
		t.Log(id)
	}

	// index
	c.EnsureUniqueIndex("test", "counters", []string{"key1"})

	// Output:
	// 1
	// 2
	// 3
}
