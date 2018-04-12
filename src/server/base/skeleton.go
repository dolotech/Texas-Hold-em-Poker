package base

import (
	"github.com/dolotech/leaf/chanrpc"
	"github.com/dolotech/leaf/module"
	"server/conf"
)

func NewSkeleton() *module.Skeleton {
	skeleton := &module.Skeleton{
		//GoLen:              conf.GoLen,
		TimerDispatcherLen: conf.TimerDispatcherLen,
		AsynCallLen:        conf.AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(conf.ChanRPCLen),
	}
	skeleton.Init()
	return skeleton
}
