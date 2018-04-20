package model

import (
	"time"
	"github.com/dolotech/lib/db"
)

//房间基本信息

type Room struct {
	Rid             uint32    `xorm:"'rid' pk autoincr BIGINT"`
	Number          string    `xorm:"'number' index not null VARCHAR(8)"` // 给玩家展示的房间号
	Pwd             string    `xorm:"'pwd' VARCHAR(16)"`                  //房间锁--密码
	State           uint8     `xorm:"'state' smallint"`                   //房间状态 0默认可用 1不可用
	Name            string    `xorm:"'name' VARCHAR(16)"`                 //房间名字
	CreatedAt       time.Time `xorm:"'created_at' index  created"`        //创建时间
	OriginalOwnerID uint32    `xorm:"'original_owner_id'"`                //原始创建人的信息
	Owner           uint32    `xorm:"'owner'"`                            //房管
	Kind            uint32    `xorm:"'kind'"`                             //游戏类型 即玩法
	DraginChips     uint32    `xorm:"'dragin_chips'"`                     //带入筹码

	//Occupants       []*uint32 `xorm:"'occupants'"`                        // 玩家列表，列表第一项为庄家
}

func (u *Room) Insert() (int64, error) {
	return db.C().Engine().InsertOne(u)
}

func (this *Room) GetById() (bool, error) {
	return db.C().Engine().Where("uid = ?", this.Rid).Get(this)
}


func (r *Room) CreatedTime() uint32 {
	return uint32(r.CreatedAt.Unix())
}