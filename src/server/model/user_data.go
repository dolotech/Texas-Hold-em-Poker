package model

import (
	"server/msg"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"github.com/golang/glog"
	"time"
	"strconv"
)

type UserData struct {
	//数据库的数据
	UserID    uint32 "_id" //用户id 自增型的
	AccountID string    //用户线上看到的id

	TotalCount uint32 //总局数
	WinCount   uint32 //胜利次数

	CreatedAt   uint32  //注册时间
	UnionId     string //微信id
	AccessToken string //token
	Uid         uint32 // 用户id
	DeviceId    string // 设备id

	Nickname   string // 微信昵称
	Sex        uint8  // 微信性别 0-未知，1-男，2-女
	Profile    string // 微信头像
	Invitecode string // 绑定的邀请码
	Chips      uint32 // 筹码
	Lv         uint8  // 等级

	LastLoginTime uint32 // 上次登录时间
	LastLoginIp   uint32 // 最后登录ip
	Type          uint8  // 用户类型
	Disable       bool   // 是否禁用
	Signature     string // 个性签名
	Gps           string // gps定位数据
	Black         bool   // 黑名单列表
}

//func init()  {
//	skeletonRegister(&msg.UserLoginInfo{},login)
//	//skeletonRegister(&msg.UserLoginInfo{},register)
//}
//
//func skeletonRegister(m interface{}, h interface{})  {
//	skeleton.RegisterChanRPC(reflect.TypeOf(m),h)
//}

func (data *UserData) initValue() error {
	userID, err := mongoDBNextSeq(USERDB)
	if err != nil {
		return fmt.Errorf("get next users id error: %v", err)
	}

	data.UserID = userID
	data.AccountID = time.Now().Format("0102") + strconv.Itoa(int(data.UserID))
	data.CreatedAt =uint32( time.Now().Unix())
	return nil
}

func (data *UserData) Register(userInfo *msg.RegisterUserInfo) error { //注册

	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).Insert(userInfo)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	return nil
}

func (data *UserData) Login(user *msg.UserLoginInfo) error {
	var result UserData
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	// check user
	err := db.DB(DB).C(USERDB).Find(bson.M{"name": user.Name, "pwd": user.Pwd}).One(&result)
	if err != nil {
		//glog.Fatal("login err - %v",err)
		//a := args[1].(gate.Agent)
		//ChanRPC.Go("LoginAgent",&msg.LoginError{1,"no user"})
		//glog.Infoln("---over----?")
		//time.Sleep(15*time.Second)
		//a.WriteMsg(&msg.LoginError{State:-1,Message:"no user"})

		glog.Errorln(err)
		return err
	}
	return nil
}

//检查用户是否已注册过
func (data *UserData) ExitedUser(userName string) (err error) {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	var userInfo msg.RegisterUserInfo
	err = db.DB(DB).C(USERDB).Find(bson.M{"name": userName}).One(&userInfo)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}
