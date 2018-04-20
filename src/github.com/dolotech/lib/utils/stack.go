package utils

import (
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
)

// 产生panic时的调用栈打印
func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		glog.Errorln(x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			glog.Errorf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			glog.Errorf("EXRAS#%v DATA:%v\n", k, spew.Sdump(extras[k]))
		}
	}
}
