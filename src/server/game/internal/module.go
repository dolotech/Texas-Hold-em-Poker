package internal

import (
	"github.com/name5566/leaf/module"
	"server/base"
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

}
