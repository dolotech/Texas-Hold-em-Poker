package db

import (
	//"game/player"
	//"github.com/bmizerany/assert"
	_ "github.com/lib/pq"
	//l4g "lib/log4go"
	"testing"
)

const UCDBURL = "postgres://postgres:postgres@192.168.1.240:3021/postgres?sslmode=disable"

func Test_NewClient(t *testing.T) {
	//l4g.Info("test client ...")

	//client, err := NewClient("postgres", UCDBURL)
	//client.engine.ShowSQL(true)
	//assert.Equal(t, err, nil)
	//defer client.Close()

	//l4g.Info("Start New Client", client)
}

func Test_Insert(t *testing.T) {
	//client, err := NewClient("postgres", UCDBURL)
	//assert.Equal(t, err, nil)
	//defer client.Close()
	//
	//l4g.Info("XXXXX", client)
	//client.Insert(&player.User{},)
}
