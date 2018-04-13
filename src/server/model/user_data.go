package model

import (
	"fmt"
	"github.com/golang/glog"
	"time"
	"github.com/dolotech/lib/db"
	"github.com/dolotech/lib/utils"
)

func (this *User) GetById() (bool, error) {
	return db.C().Engine().Where("uid = ?", this.Uid).Get(this)
}

func (this *User) GetByUnionId() (bool, error) {
	return db.C().Engine().Where("union_id = ?", this.UnionId).Get(this)
}

type User struct {
	Uid        uint32 `xorm:"'uid' pk autoincr BIGINT"`  // 用户id
	DeviceId   string `xorm:"'device_id' VARCHAR(32)"`   // 设备id
	UnionId    string `xorm:"'union_id' VARCHAR(32)"`    // 微信联合id
	Nickname   string `xorm:"'nickname' VARCHAR(32)"`    // 微信昵称
	Sex        uint8  `xorm:"'sex' smallint"`            // 微信性别 0-未知，1-男，2-女
	Profile    string `xorm:"'profile' VARCHAR(64)"`     // 微信头像
	Invitecode string `xorm:"'invitecode' VARCHAR(6)"`   // 绑定的邀请码
	Chips      uint32 `xorm:"'chips'"`                   // 筹码
	Lv         uint8  `xorm:"'lv' smallint"`             // 等级
	CreatedAt  uint32 `xorm:"'created_at'"`              // 注册时间
	LastTime   uint32 `xorm:"'last_time'"`               // 上次登录时间
	LastIp     uint32 `xorm:"'last_ip' BIGINT"`          // 最后登录ip
	Type       uint8  `xorm:"'type'  not null smallint"` // 用户类型
	Disable    bool   `xorm:"'disable'"`                 // 是否禁用
	Signature  string `xorm:"'signature' VARCHAR(64)"`   // 个性签名
	Gps        string `xorm:"'gps' VARCHAR(32)"`         // gps定位数据
	Black      bool   `xorm:"'black'"`                   // 黑名单列表
	RoomID     uint32 `xorm:"'room_id'"`                 // 当前所在房间号，0表示不在房间,用于掉线重连
}

func (u *User) Insert( ) (int64,error) {
	return db.C().Engine().InsertOne(u)
}
func (u *User) UpdateChips(value int32) error {
	sql := fmt.Sprintf(`UPDATE public.user SET
	chips = chips + %d WHERE uid = %d `, value, u.Uid)
	_, err := db.C().Engine().Exec(sql)
	if err != nil {
		glog.Errorln(err)
	}
	return err
}

func (u *User) UpdateLogin(ip string) error {
	sql := fmt.Sprintf(`UPDATE public.user SET
	last_login_time =  %d ,last_login_ip =  %d WHERE uid = %d `, uint32(time.Now().Unix()), utils.InetToaton(ip), u.Uid)
	_, err := db.C().Engine().Exec(sql)
	if err != nil {
		glog.Errorln(err)
	}
	return err
}

/*

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
*/
