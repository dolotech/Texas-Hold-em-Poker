package model

import "server/algorithm"

//房间基本信息

type RoomData struct {
	RoomID          uint32
	RoomNumber      string
	Volume          uint32    //房间容量
	RoomPwd         string //房间锁--密码
	RoomState       uint8    //房间状态 0默认可用 1不可用
	RoomName        string //房间名字
	CreatedAt       uint32  //创建时间
	OriginalOwnerID uint32    //原始创建人的信息
	RoomOwner       uint32    //房管
	BaseMoney       uint32    //最低资本 才能进房间
	PayValue        uint32    //倍数
	GameType        uint32    //游戏类型 即玩法

	Id        string
	SB        uint32             // 小盲注
	BB        uint32             // 大盲注
	Cards     algorithm.Cards //公共牌
	Pot       []uint32           // 当前奖池筹码数
	Timeout   uint32             // 倒计时超时时间(秒)
	Button    uint32             // 当前庄家座位号，从1开始
	Occupants []*uint32          // 玩家列表，列表第一项为庄家
	Chips     []uint32           // 玩家本局下注的总筹码数，与occupants一一对应
	Bet       uint32             // 当前下注额
	N         uint32             // 当前玩家人数
	Max       uint32             // 房间最大玩家人数
	MaxChips  uint32
	MinChips  uint32
}

//func (R *RoomData) InitValue(userID int) error {
//	roomID, err := mongoDBNextSeq("rooms")
//	if err != nil {
//		return fmt.Errorf("get next rooms id error: %v", err)
//	}
//
//	R.RoomID = roomID
//	R.RoomNumber = fmt.Sprintf("%06d", roomID)
//	R.OriginalOwnerID = userID
//	R.RoomOwner = userID
//	R.CreatedAt = time.Now().Unix()
//	return nil
//}

//创建 往数据库写入
