package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"server/conf"
	"server/game"
	"server/gate"
	"server/login"
	"github.com/golang/glog"
	"flag"
)

func main() {
	flag.Parse()
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath
	glog.Errorf("%v", 12312)
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)

}
