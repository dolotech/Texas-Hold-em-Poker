package model

import (
	"fmt"
	"github.com/golang/glog"
	"time"
	"github.com/dolotech/lib/db"
	"github.com/dolotech/lib/utils"
	"math/rand"
)

func (this *User) GetById() (bool, error) {
	return db.C().Engine().Where("uid = ?", this.Uid).Get(this)
}

func (this *User) GetByAccount() (bool, error) {
	return db.C().Engine().Where("account = ?", this.Account).Get(this)
}

func (this *User) GetByUnionId() (bool, error) {
	return db.C().Engine().Where("union_id = ?", this.UnionId).Get(this)
}

type User struct {
	Uid        uint32    `xorm:"'uid' pk autoincr BIGINT"`            // 用户id
	Account    string    `xorm:"'account' index unique  VARCHAR(16)"` // 客户端玩家展示的账号
	DeviceId   string    `xorm:"'device_id' VARCHAR(32)"`             // 设备id
	UnionId    string    `xorm:"'union_id' VARCHAR(32)"`              // 微信联合id
	Nickname   string    `xorm:"'nickname' VARCHAR(32)"`              // 微信昵称
	Sex        uint8     `xorm:"'sex' smallint"`                      // 微信性别 0-未知，1-男，2-女
	Profile    string    `xorm:"'profile' VARCHAR(64)"`               // 微信头像
	Invitecode string    `xorm:"'invitecode' VARCHAR(6)"`             // 绑定的邀请码
	Chips      uint32    `xorm:"'chips'"`                             // 筹码
	Lv         uint8     `xorm:"'lv' smallint"`                       // 等级
	CreatedAt  time.Time `xorm:"'created_at' index  created"`         // 注册时间
	LastTime   time.Time `xorm:"'last_time'"`                         // 上次登录时间
	LastIp     uint32    `xorm:"'last_ip' BIGINT"`                    // 最后登录ip
	Kind       uint8     `xorm:"'kind'  not null smallint"`           // 用户类型
	Disable    bool      `xorm:"'disable'"`                           // 是否禁用
	Signature  string    `xorm:"'signature' VARCHAR(64)"`             // 个性签名
	Gps        string    `xorm:"'gps' VARCHAR(32)"`                   // gps定位数据
	Black      bool      `xorm:"'black'"`                             // 黑名单列表
	RoomID     uint32    `xorm:"'room_id'"`                           // 当前所在房间号，0表示不在房间,用于掉线重连
}

func (u *User) Insert() error {
	_, err := db.C().Engine().InsertOne(u)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	n := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(9999-1000) + 1000

	account := fmt.Sprintf("%v%v", u.Uid, n)
	sql := `UPDATE public.user SET account =$1 WHERE uid = $2 `
	_, err = db.C().Engine().Exec(sql, account, u.Uid)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	glog.Errorln(u)
	return nil
}
func (u *User) UpdateChips(value int32) error {
	_, err := db.C().Engine().Exec(`UPDATE public.user SET
		chips = chips + $1 WHERE uid =$2 `, value, u.Uid)
	if err != nil {
		glog.Errorln(err)
	}

	//glog.Errorln(res.RowsAffected())
	return err

	//s:=db.C().Engine().Table(u).Incr("chips",value)

	//return nil
}

func (u *User) UpdateLogin(ip string) error {
	sql := `UPDATE public.user SET
	last_time =  $1 ,last_ip =  $2 WHERE uid = $3 `
	_, err := db.C().Engine().Exec(sql, time.Now(), utils.InetToaton(ip), u.Uid)
	if err != nil {
		glog.Errorln(err)
	}
	return err
}

func (u *User) UpdateRoomId() error {
	sql := `UPDATE public.user SET
	room_id = chips + $1 WHERE uid = $2 `
	_, err := db.C().Engine().Exec(sql, u.RoomID, u.Uid)
	if err != nil {
		glog.Errorln(err)
	}
	return err
}
