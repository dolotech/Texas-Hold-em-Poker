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
	PayValue        uint8     `xorm:"'pay_value' smallint"`               //倍数
	SB              uint8     `xorm:"'sb' smallint"`                      // 小盲注
	BB              uint8     `xorm:"'bb' smallint"`                      // 大盲注
	Cards           [] uint8  `xorm:"'cards'"`                            //公共牌
	Pot             []uint32  `xorm:"'pot'"`                              // 当前奖池筹码数
	Timeout         uint8     `xorm:"'timeout' smallint"`                 // 倒计时超时时间(秒)
	Button          uint32    `xorm:"'button'"`                           // 当前庄家座位号，从1开始
	Occupants       []*uint32 `xorm:"'occupants'"`                        // 玩家列表，列表第一项为庄家
	Chips           []uint32  `xorm:"'chips'"`                            // 玩家本局下注的总筹码数，与occupants一一对应
	Bet             uint32    `xorm:"'bet'"`                              // 当前下注额
	N               uint8     `xorm:"'n' smallint"`                       // 当前玩家人数
	Max             uint8     `xorm:"'max' smallint"`                     // 房间最大玩家人数
	MaxChips        uint32    `xorm:"'maxchips'"`
	MinChips        uint32    `xorm:"'minchips'"`
}

func (u *Room) Insert() (int64, error) {
	return db.C().Engine().InsertOne(u)
}

func (this *Room) GetById() (bool, error) {
	return db.C().Engine().Where("uid = ?", this.Rid).Get(this)
}
