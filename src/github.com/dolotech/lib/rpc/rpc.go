/**
* Created by Michael on 2016/8/5.
*	web服务器与推送服务器通讯的RPC调用
*
*
 */
package rpc

import (
	"net"
	"net/rpc"
	//"strconv"

	log "github.com/golang/glog"
)

type RPCServer struct {
	rpc  *rpc.Server
	port string
}

type RPCClient struct {
	rpc       *rpc.Client
	addr      string
	conn      *net.TCPConn
	Connected bool
}

// 新建RPC客户端
func CreateClient(addr string) *RPCClient {
	return &RPCClient{addr: addr}

}

// 新建RPC服务端
func CreateServer() *RPCServer {
	server := &RPCServer{}
	server.rpc = rpc.NewServer()
	return server
}

// serviceMethod  方法名字包括类名 eg: "Receiver.Receive"
// args  传参数
// reply	返回的数据结构，一定是指针
func (this *RPCClient) Call(serviceMethod string, args interface{}, reply interface{}) error {

	if this.Connected == false {
		e := this.connect()
		if e != nil {
			return e
		}
		this.Connected = true
	}
	//	defer this.rpc.Close()

	err := this.rpc.Call(serviceMethod, args, reply)
	if err != nil {
		log.Errorln("error:", err)
		e := this.connect()
		if e == nil {
			this.rpc.Call(serviceMethod, args, reply)
		}
	}

	return nil
}

func (this *RPCClient) connect() error {
	address, err := net.ResolveTCPAddr("tcp", this.addr)
	if err != nil {
		log.Errorln(err)
		return err
	}
	if this.conn != nil {
		this.conn.Close()
	}
	this.conn, err = net.DialTCP("tcp", nil, address)
	if err != nil {
		log.Errorln(err)
		return err
	}

	this.rpc = rpc.NewClient(this.conn)
	if err != nil {
		log.Errorln(err)
		return err
	}
	//	defer conn.Close()

	return nil
}

// 注册服务端方法,eg.  server.Register(new(Receiver))
func (this *RPCServer) Register(rcvr interface{}) {
	this.rpc.Register(rcvr)
}

func (this *RPCServer) Listen(addr string) {
	//this.port = strconv.Itoa(port)
	l, e := net.Listen("tcp", addr) // any available address
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}
	this.rpc.Accept(l)
}
