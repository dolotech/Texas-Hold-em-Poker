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
	AccountID string       //用户线上看到的id

	TotalCount uint32 //总局数
	WinCount   uint32 //胜利次数

	CreatedAt   uint32 //注册时间
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
	Phone      string // 手机号码

	LastLoginTime uint32 // 上次登录时间
	LastLoginIp   uint32 //上次登录ip
	Type          uint8  // 用户类型
	Disable       bool   // 是否禁用
	Signature     string // 个性签名
	Gps           string // gps定位数据
	Black         bool   // 黑名单列表

	RoomID uint32 // 当前所在房间号，0表示不在房间,用于掉线重连
}

func (data *UserData) initValue() error {
	userID, err := mongoDBNextSeq(USERDB)
	if err != nil {
		return fmt.Errorf("get next users id error: %v", err)
	}

	data.UserID = userID
	data.AccountID = time.Now().Format("0102") + strconv.Itoa(int(data.UserID))
	data.CreatedAt = uint32(time.Now().Unix())
	return nil
}

func (data *UserData) GetByWechat() error {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).Find(bson.M{"unionid": data.UnionId}).One(data)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (data *UserData) IncChips(change int) error {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).UpdateId(data.UserID, bson.M{"$inc": bson.M{"chips": change}})
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (data *UserData) UpdateSex() error {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).UpdateId(data.UserID, bson.M{"$set": bson.M{"sex": data.Sex}})
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (data *UserData) Insert() error { //注册
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	err := db.DB(DB).C(USERDB).Insert(data)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (data *UserData) Register() error { //注册
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)

	err := data.Insert()
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
	err := db.DB(DB).C(USERDB).Find(bson.M{"name": user.Name, "pwd": user.Pwd}).One(&result)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

//检查用户是否已注册过
func (data *UserData) ExistByAccountID() bool {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)
	count, _ := db.DB(DB).C(USERDB).Find(bson.M{"accountid": data.AccountID}).Count()
	return count > 0
}
