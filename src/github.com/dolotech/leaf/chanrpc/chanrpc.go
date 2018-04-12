package chanrpc

import (
	"errors"
	"fmt"
	"github.com/dolotech/leaf/conf"
	"github.com/golang/glog"
	"runtime"
	"reflect"
)

// one server per goroutine (goroutine not safe)
// one client per goroutine (goroutine not safe)
type Server struct {
	// id -> function
	//
	// function:
	// func(args []interface{})
	// func(args []interface{}) interface{}
	// func(args []interface{}) []interface{}
	functions map[interface{}]*reflect.Value
	ChanCall  chan *CallInfo
}

type CallInfo struct {
	f       *reflect.Value
	args    []interface{}
	chanRet chan *RetInfo
	cb      interface{}
}

type RetInfo struct {
	// nil
	// interface{}
	// []interface{}
	ret interface{}
	err error
	// callback:
	// func(err error)
	// func(ret interface{}, err error)
	// func(ret []interface{}, err error)
	cb interface{}
}

type Client struct {
	s               *Server
	chanSyncRet     chan *RetInfo
	ChanAsynRet     chan *RetInfo
	pendingAsynCall int
}

func NewServer(l int) *Server {
	s := new(Server)
	s.functions = make(map[interface{}]*reflect.Value)
	s.ChanCall = make(chan *CallInfo, l)
	return s
}

func assert(i interface{}) []interface{} {
	if i == nil {
		return nil
	} else {
		return i.([]interface{})
	}
}

// you must call the function before calling Open and Go
func (s *Server) Register(id interface{}, f interface{}) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		glog.Warning("消息处理不是函数", id)
		return
	}
	if _, ok := s.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	v := reflect.ValueOf(f)
	s.functions[id] = &v
}

func (s *Server) ret(ci *CallInfo, ri *RetInfo) (err error) {
	if ci.chanRet == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	ri.cb = ci.cb
	ci.chanRet <- ri
	return
}

func (s *Server) Exec(ci *CallInfo) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				fmt.Errorf("%v: %s", r, buf[:l])
			} else {
				fmt.Errorf("%v", r)
			}

			s.ret(ci, &RetInfo{err: fmt.Errorf("%v", r)})
		}
	}()

	agrs := make([]reflect.Value, 0, len(ci.args))
	for i := 0; i < len(ci.args); i++ {
		agrs = append(agrs, reflect.ValueOf(ci.args[i]))
	}
	ci.f.Call(agrs)
}

// goroutine safe
func (s *Server) Go(id interface{}, args ...interface{}) {
	f, _ := s.functions[id]
	if f == nil {
		return
	}

	defer func() {
		recover()
	}()

	s.ChanCall <- &CallInfo{
		f:    f,
		args: args,
	}
}

// goroutine safe
func (s *Server) Call0(id interface{}, args ...interface{}) error {
	return s.Open(0).Call0(id, args...)
}

// goroutine safe
func (s *Server) Call1(id interface{}, args ...interface{}) (interface{}, error) {
	return s.Open(0).Call1(id, args...)
}

// goroutine safe
func (s *Server) CallN(id interface{}, args ...interface{}) ([]interface{}, error) {
	return s.Open(0).CallN(id, args...)
}

func (s *Server) Close() {
	close(s.ChanCall)

	for ci := range s.ChanCall {
		s.ret(ci, &RetInfo{
			err: errors.New("chanrpc server closed"),
		})
	}
}

// goroutine safe
func (s *Server) Open(l int) *Client {
	c := NewClient(l)
	c.Attach(s)
	return c
}

func NewClient(l int) *Client {
	c := new(Client)
	c.chanSyncRet = make(chan *RetInfo, 1)
	c.ChanAsynRet = make(chan *RetInfo, l)
	return c
}

func (c *Client) Attach(s *Server) {
	c.s = s
}

func (c *Client) call(ci *CallInfo, block bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	if block {
		c.s.ChanCall <- ci
	} else {
		select {
		case c.s.ChanCall <- ci:
		default:
			err = errors.New("chanrpc channel full")
		}
	}
	return
}

func (c *Client) f(id interface{}, n int) (*reflect.Value,  error) {
	if c.s == nil {
		return nil, errors.New("server not attached")
	}
	f := c.s.functions[id]
	if f == nil {
		return nil,fmt.Errorf("function id %v: function not registered", id)
	}
	return f,nil
}

func (c *Client) Call0(id interface{}, args ...interface{}) error {
	f, err := c.f(id, 0)
	if err != nil {
		return err
	}

	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return err
	}

	ri := <-c.chanSyncRet
	return ri.err
}

func (c *Client) Call1(id interface{}, args ...interface{}) (interface{}, error) {
	f, err := c.f(id, 1)
	if err != nil {
		return nil, err
	}
	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return nil, err
	}

	ri := <-c.chanSyncRet
	return ri.ret, ri.err
}

func (c *Client) CallN(id interface{}, args ...interface{}) ([]interface{}, error) {
	f, err := c.f(id, 2)
	if err != nil {
		return nil, err
	}

	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return nil, err
	}

	ri := <-c.chanSyncRet
	return assert(ri.ret), ri.err
}

func (c *Client) asynCall(id interface{}, args []interface{}, cb interface{}, n int) {
	f, err := c.f(id, n)
	if err != nil {
		c.ChanAsynRet <- &RetInfo{err: err, cb: cb}
		return
	}
	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.ChanAsynRet,
		cb:      cb,
	}, false)
	if err != nil {
		c.ChanAsynRet <- &RetInfo{err: err, cb: cb}
		return
	}
}

func (c *Client) AsynCall(id interface{}, _args ...interface{}) {
	if len(_args) < 1 {
		panic("callback function not found")
	}

	args := _args[:len(_args)-1]
	cb := _args[len(_args)-1]

	var n int
	switch cb.(type) {
	case func(error):
		n = 0
	case func(interface{}, error):
		n = 1
	case func([]interface{}, error):
		n = 2
	default:
		panic("definition of callback function is invalid")
	}

	// too many calls
	if c.pendingAsynCall >= cap(c.ChanAsynRet) {
		execCb(&RetInfo{err: errors.New("too many calls"), cb: cb})
		return
	}

	c.asynCall(id, args, cb, n)
	c.pendingAsynCall++
}

func execCb(ri *RetInfo) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				glog.Error("%v: %s", r, buf[:l])
			} else {
				glog.Error("%v", r)
			}
		}
	}()

	// execute
	switch ri.cb.(type) {
	case func(error):
		ri.cb.(func(error))(ri.err)
	case func(interface{}, error):
		ri.cb.(func(interface{}, error))(ri.ret, ri.err)
	case func([]interface{}, error):
		ri.cb.(func([]interface{}, error))(assert(ri.ret), ri.err)
	default:
		panic("bug")
	}
	return
}

func (c *Client) Cb(ri *RetInfo) {
	c.pendingAsynCall--
	execCb(ri)
}

func (c *Client) Close() {
	for c.pendingAsynCall > 0 {
		c.Cb(<-c.ChanAsynRet)
	}
}

func (c *Client) Idle() bool {
	return c.pendingAsynCall == 0
}
