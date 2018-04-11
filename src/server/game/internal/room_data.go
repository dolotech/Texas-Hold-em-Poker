package internal

//房间基本信息

type RoomData struct {
	RoomID	int
	RoomNumber string
	Volume	int //房间容量
	RoomPwd	string	//房间锁--密码
	RoomState	int	//房间状态 0默认可用 1不可用
	RoomName	string	//房间名字
	CreatedAt	int64	//创建时间
	OriginalOwnerID	int //原始创建人的信息
	RoomOwner int	//房管
	BaseMoney	int	//最低资本 才能进房间
	PayValue	int	//倍数
	GameType	int	//游戏类型 即玩法
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




