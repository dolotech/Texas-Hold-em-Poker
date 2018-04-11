package internal

import (
	"server/msg"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"github.com/name5566/leaf/log"
	"time"
	"strconv"
)

type UserData struct {//数据库的数据
	UserID	int	"_id"	//用户id 自增型的
	AccountID	string	//用户线上看到的id
	NickName	string	//用户的昵称
	Sex int	//性别 0--女 1--男
	TotalCount	int	//比赛总次数
	WinCount	int	//胜利次数
	Money	int	//账号金币
	HeadImgUrl	string	//头像
	CreatedAt	int64	//注册时间
	UnionId	string	//微信id
	AccessToken	string //token
		//Name string
		//Pwd string
		//Age int
		//Address string
}

const USERDB  = "users"

//func init()  {
//	skeletonRegister(&msg.UserLoginInfo{},login)
//	//skeletonRegister(&msg.UserLoginInfo{},register)
//}
//
//func skeletonRegister(m interface{}, h interface{})  {
//	skeleton.RegisterChanRPC(reflect.TypeOf(m),h)
//}

func (data *UserData) initValue() error {
	userID, err := mongoDBNextSeq("users")
	if err != nil {
		return fmt.Errorf("get next users id error: %v", err)
	}

	data.UserID = userID
	data.AccountID = time.Now().Format("0102") + strconv.Itoa(data.UserID)
	data.CreatedAt = time.Now().Unix()
	return nil
}

func register(userInfo *msg.RegisterUserInfo)  (err error) {//注册
	//var user User
	//userInfo := args[0].(*msg.RegisterUserInfo)
	skeleton.Go(func() {
		db := mongoDB.Ref()
		defer mongoDB.UnRef(db)
		err := db.DB(DB).C(USERDB).Insert(userInfo)
		if err != nil{
			//log.Fatal("err register --%v",err)
			log.Fatal("err register - %v, err ",err )
			return
 		}
	},nil)
	return
}

func login(user  *msg.UserLoginInfo)(err error) {
	//var user User
	//fmt.Println("---lognin------",args)
	//user := args[0].(*msg.UserLoginInfo)
	fmt.Println("---userinfo---",user)
	var result UserData
	skeleton.Go(func() {
		db := mongoDB.Ref()
		defer mongoDB.UnRef(db)
		// check user
		err := db.DB(DB).C(USERDB).Find(bson.M{"name":user.Name,"pwd":user.Pwd}).One(&result)
		if err != nil{
			//log.Fatal("login err - %v",err)
			//a := args[1].(gate.Agent)
			//ChanRPC.Go("LoginAgent",&msg.LoginError{1,"no user"})
			//fmt.Println("---over----?")
			//time.Sleep(15*time.Second)
			//a.WriteMsg(&msg.LoginError{State:-1,Message:"no user"})
			log.Fatal("login err - %v",err)

			return

		}
	}, nil)
	return
}

//检查用户是否已注册过
func checkExitedUser(userName string) (err error){
	//skeleton.Go(func() {
	//	db := mongoDB.Ref()
	//	defer mongoDB.UnRef(db)
	//	err := db.DB(DB).C(USERDB).Find(bson.M{"name":bson.M{"$exists":userName}})
	//	if err != nil {
	//
	//	}
	//},nil)

	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	var userInfo msg.RegisterUserInfo
	err = db.DB(DB).C(USERDB).Find(bson.M{"name":userName}).One(&userInfo)
	if err != nil{
		return err
	}
	return nil
}

