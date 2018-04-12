package mongodb

import (
	"fmt"
	"github.com/name5566/leaf/db/mongodb"
	"gopkg.in/mgo.v2"
	"github.com/golang/glog"
)

func Example() {
	c, err := mongodb.Dial("localhost", 10)
	if err != nil {
		glog.Infoln(err)
		return
	}
	defer c.Close()

	// session
	s := c.Ref()
	defer c.UnRef(s)
	err = s.DB("test").C("counters").RemoveId("test")
	if err != nil && err != mgo.ErrNotFound {
		glog.Infoln(err)
		return
	}

	// auto increment
	err = c.EnsureCounter("test", "counters", "test")
	if err != nil {
		glog.Infoln(err)
		return
	}
	for i := 0; i < 3; i++ {
		id, err := c.NextSeq("test", "counters", "test")
		if err != nil {
			glog.Infoln(err)
			return
		}
		glog.Infoln(id)
	}

	// index
	c.EnsureUniqueIndex("test", "counters", []string{"key1"})

	// Output:
	// 1
	// 2
	// 3
}
