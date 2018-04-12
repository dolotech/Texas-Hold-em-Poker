package chanrpc

import (
	"sync"
	"testing"
)

func TestClient_AsynCall(t *testing.T) {
	s := NewServer(10)

	var wg sync.WaitGroup
	wg.Add(1)

	// goroutine 1
	go func() {
		s.Register("f0", func(args []interface{}) {

		})

		s.Register("f1", func(args []interface{}) interface{} {
			return 1
		})

		s.Register("fn", func(args []interface{}) []interface{} {
			return []interface{}{1, 2, 3}
		})

		s.Register("add", func(args []interface{}) interface{} {
			n1 := args[0].(int)
			n2 := args[1].(int)
			return n1 + n2
		})

		wg.Done()

		for {
			s.Exec(<-s.ChanCall)
		}
	}()

	wg.Wait()
	wg.Add(1)

	// goroutine 2
	go func() {
		c := s.Open(10)

		// sync
		err := c.Call0("f0")
		if err != nil {
			t.Log(err)
		}

		r1, err := c.Call1("f1")
		if err != nil {
			t.Log(err)
		} else {
			t.Log(r1)
		}

		rn, err := c.CallN("fn")
		if err != nil {
			t.Log(err)
		} else {
			t.Log(rn[0], rn[1], rn[2])
		}

		ra, err := c.Call1("add", 1, 2)
		if err != nil {
			t.Log(err)
		} else {
			t.Log(ra)
		}

		// asyn
		c.AsynCall("f0", func(err error) {
			if err != nil {
				t.Log(err)
			}
		})

		c.AsynCall("f1", func(ret interface{}, err error) {
			if err != nil {
				t.Log(err)
			} else {
				t.Log(ret)
			}
		})

		c.AsynCall("fn", func(ret []interface{}, err error) {
			if err != nil {
				t.Log(err)
			} else {
				t.Log(ret[0], ret[1], ret[2])
			}
		})

		c.AsynCall("add", 1, 2, func(ret interface{}, err error) {
			if err != nil {
				t.Log(err)
			} else {
				t.Log(ret)
			}
		})

		c.Cb(<-c.ChanAsynRet)
		c.Cb(<-c.ChanAsynRet)
		c.Cb(<-c.ChanAsynRet)
		c.Cb(<-c.ChanAsynRet)

		// go
		s.Go("f0")

		wg.Done()
	}()

	wg.Wait()

	// Output:
	// 1
	// 1 2 3
	// 3
	// 1
	// 1 2 3
	// 3
}
