package main

import (
	"github.com/dolotech/leaf"
	lconf "github.com/dolotech/leaf/conf"
	"server/conf"
	"server/game"
	"server/gate"
	"server/login"
	"net/http"
	"flag"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath
	go leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		glog.Fatalf("ListenAndServe: %v ", err)
	}
}
