package internal

import (
	"github.com/dolotech/leaf/module"
	"server/base"
	"github.com/golang/glog"
)

var (//定义
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {//相当于继承父类定义
	*module.Skeleton
}

func (m *Module) OnInit() {//继承初始
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {
	glog.Errorln("OnDestroy")
}
