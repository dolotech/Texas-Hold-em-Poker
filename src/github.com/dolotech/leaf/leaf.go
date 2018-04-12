package leaf

import (
	"github.com/dolotech/leaf/cluster"
	"github.com/dolotech/leaf/console"
	"github.com/golang/glog"
	"github.com/dolotech/leaf/module"
	"os"
	"os/signal"
)

func Run(mods ...module.Module) {

	glog.Errorf("Leaf %v starting up", version)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	glog.Errorf("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}
