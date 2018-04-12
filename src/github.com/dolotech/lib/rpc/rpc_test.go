package rpc

import (
	"net/rpc"
	//"net/http"
	"net"
	"testing"
	"time"

	log "github.com/golang/glog"
)

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *([]string)) error {
	*reply = append(*reply, "test")
	log.Infoln(args)
	return nil
}

func TestProtoFriendListData(t *testing.T) {
	newServer := rpc.NewServer()
	newServer.Register(new(Arith))

	l, e := net.Listen("tcp", "127.0.0.1:1234") // any available address
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}

	go newServer.Accept(l)
	newServer.HandleHTTP("/foo", "/bar")
	time.Sleep(2 * time.Second)

	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	conn, _ := net.DialTCP("tcp", nil, address)
	defer conn.Close()

	client := rpc.NewClient(conn)
	defer client.Close()

	args := &Args{7, 8}
	reply := make([]string, 10)
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Infoln(reply)
}
