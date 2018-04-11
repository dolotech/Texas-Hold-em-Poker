package internal

import (
	"github.com/name5566/leaf/db/mongodb"
	"github.com/name5566/leaf/log"
	"server/conf"
)

//const DB_INFO  = "mongodb://yin_test:123456@localhost:27017/runmongo"
var mongoDB *mongodb.DialContext

//const DB  = "runmongo"
const DB  = "mahjongdb"

func init()  {

	//db, err := mongodb.Dial(DB_INFO,10)
	//if err != nil{
	//	//fmt.Println("----connecting----")
	//	//log.Fatal("db %v",err)
	//	log.Fatal("db-err %v",err)
	//	//fmt.Println("------connected----")
	//	return
	//}
	//mongoDB = db
	////fmt.Println("------connected----")
	// mongodb
	if conf.Server.DBMaxConnNum <= 0 {
		conf.Server.DBMaxConnNum = 100
	}
	db, err := mongodb.Dial(conf.Server.DBUrl, conf.Server.DBMaxConnNum)
	if err != nil {
		log.Fatal("dial mongodb error: %v", err)
	}
	mongoDB = db

	err = db.EnsureCounter(DB, "counters", "users")
	if err != nil {
		log.Fatal("ensure counter error: %v", err)
	}

	err = db.EnsureCounter(DB, "counters", "rooms")
	if err != nil {
		log.Fatal("ensure counter error: %v", err)
	}
}

func mongoDBDestroy()  {
	mongoDB.Close()
	mongoDB = nil

}

func mongoDBNextSeq(id string) (int, error) {
	return mongoDB.NextSeq(DB, "counters", id)
}
