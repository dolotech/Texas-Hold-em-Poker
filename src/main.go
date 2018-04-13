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
	"github.com/dolotech/lib/db"
	"server/model"
)

var Commit = ""
var BUILD_TIME = ""
var VERSION = ""

var createdb bool
func main() {
	flag.StringVar(&conf.Server.WSAddr, "addr", ":8989", "")
	flag.IntVar(&conf.Server.MaxConnNum, "maxconn", 20000, "")
	flag.StringVar(&conf.Server.DBUrl, "pg", "postgres://postgres:haosql@127.0.0.1:5432/postgres?sslmode=disable", "pg addr")
	flag.BoolVar(&createdb, "createdb", false, "initial database")

	flag.Parse()
	db.Init(conf.Server.DBUrl)

	if createdb {
		createDb()
	}

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

func createDb() {
	// 建表,只维护和服务器game里面有关的表
	err := db.C().Sync(model.User{})
	if err != nil {
		glog.Errorln(err)
	}
}
