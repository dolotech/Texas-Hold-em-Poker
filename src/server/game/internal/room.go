package internal

import (
	"server/msg"
	"time"
	"github.com/dolotech/leaf/gate"
	"server/model"
)

type position int // 位置 1-4
var allRooms []*Room

type Room struct {
	RoomData *model.RoomData
	MapUsers map[uint32]*model.UserData //玩家信息以及座位 int 为位置 1-4
	//RoomState	int	//房间状态
	RoomOwner *model.UserData //房管
	//*g.LinearContext
	gate.Agent
	Players        uint32 //玩家数量
	ReadySignCount uint32 //准备信号 数量
	Playing        uint32 //房间已开始游戏 0--没有开始 1--开始
	//RoomTime	int //房间有效期限
}

func init() {
	allRooms = make([]*Room, 1) //初始化 浪费一个 第一个留着
}

//密码检查
func (r *Room) CheckPlayerAndPwd(pwd string) bool {
	if r.RoomData.RoomPwd == pwd && r.Players < r.RoomData.Volume { //密码正确 且 不满人
		return true
	} else {
		r.WriteMsg(msg.MSG_ROOM_NOTAUTH)
		return false
	}
}

//开始游戏 满人后 5s 自动开始游戏
func (r *Room) gameStart() {
	if (r.ReadySignCount == r.RoomData.Volume) && ( r.Playing == 0 ) {
		//r.WriteMsg(&)
		//	游戏开始
		//	发牌
		var cards Cards
		cards.InitCard()

		r.Playing = 1

	}
}

//查找房间 按房间Number查找
func FindRoomsByRoomNumber(roomNumber string) (room *Room) {
	for _, v := range allRooms {
		if v.RoomData.RoomNumber == roomNumber {
			room = v
		}
	}
	return
}

//通过 房间的 RoomID 确定某个房间
func FindRoomByRoomNameAndOwner(roomID uint32) (room *Room) {
	for _, v := range allRooms {
		if v.RoomData.RoomID == roomID {
			room = v
		}
	}
	return
}

//初始化房间信息 创建房间
func InitRoom(msgRoomInfo *msg.RoomInfo, userLine *UserLine) (room *Room) {
	if userLine.RoomLine == nil { //没有关联任何的房间
		r := &model.RoomData{
			RoomName:  msgRoomInfo.RoomName,
			Volume:    msgRoomInfo.Volume,
			GameType:  msgRoomInfo.GameType,
			PayValue:  msgRoomInfo.PayValue,
			BaseMoney: msgRoomInfo.BaseMoney,
			RoomPwd:   msgRoomInfo.RoomPwd,
			CreatedAt: uint32(time.Now().Unix()),
			RoomState: 1,
		}
		room.Players = 1
		room.RoomData = r
		room.RoomOwner = userLine.UserData
		room.MapUsers[1] = userLine.UserData //房主做一号位 创建房间的
		//room.LinearContext = skeleton.NewLinearContext()
		return room
	} else {
		return nil
	}
	//return room
}

//修改房间基本信息 仅有房主修改
func ChangeRoomInfo(userLine *UserLine, msgRoomInfo *msg.RoomInfo) bool {
	if userLine.RoomLine.RoomOwner != userLine.UserData { //不是房主
		r := &model.RoomData{
			RoomName:  msgRoomInfo.RoomName,
			Volume:    msgRoomInfo.Volume,
			GameType:  msgRoomInfo.GameType,
			PayValue:  msgRoomInfo.PayValue,
			BaseMoney: msgRoomInfo.BaseMoney,
			RoomPwd:   msgRoomInfo.RoomPwd,
			CreatedAt: uint32(time.Now().Unix()),
			//RoomState:0,
		}
		userLine.RoomLine.RoomData = r
		return true
	} else { //不是房主
		return false

	}
}

//加入房间
func JoinRoom(userLine *UserLine, room *Room) bool {
	if userLine.RoomLine == nil { //确定是没有房间的人 可以进入
		//	分配位置 并进行初始赋值
		room.Players = room.Players + 1
		for i := uint32(1); i <= room.RoomData.Volume; i++ {
			if v, ok := room.MapUsers[i]; !ok || v == nil { //空位置 添加玩家
				room.MapUsers[i] = userLine.UserData
				break
			}
		}
		//room.MapUsers[room.Players] = userLine.UserData
		userLine.RoomLine = room
		return true

	} else {
		//userLine.WriteMsg(&msg.CodeState{msg.MSG_ROOM_OVERVOLUME,"你已在其他房间进行游戏了，请退出房间，后操作！"})
		return false
	}
}

//加入房间前 检查是否 有“资本” 进入
func CheckConditionForJoining(userLine *UserLine, room *Room) bool {
	if userLine.UserData.Chips >= room.RoomData.PayValue {
		return true
	} else {
		//userLine.WriteMsg(&msg.CodeState{msg.MSG_ROOM_NOMONEY,"你的资金不足，该房间的要求!"})
		return false
	}
}

//房主 更换 原房主 还在 ----胜利者 当房主
func ChangeRoomOwner(userLine *UserLine) {
	//room.RoomOwner = userLine.UserData
	//userLine.RoomLine = room
	userLine.RoomLine.RoomOwner = userLine.UserData
}

//房主 更换 原房主退出 下家当房主
func ChangeRoomOlderToNewer(userLine *UserLine) {
	room := userLine.RoomLine
	//var position int
	flag := false      //标记 只进行一次使用 处理下一家当房主
	ownerFlag := false //标记 只进行一次使用 处理原房主置空 解决房主是 1 号 为场景
	for i, v := range room.MapUsers {
		if v == userLine.UserData && !ownerFlag {
			//position = i
			room.MapUsers[i] = nil //该房间 某个位置 置空
			ownerFlag = true
			//break
		}
		if v != nil && !flag && !ownerFlag {
			room.RoomOwner = v //随机 有map的迭代下产生 下一个 房主
			flag = true
		}
	}
	//room.MapUsers[]
	//delete(room.RoomOwner,userLine.UserData)
	userLine.RoomLine = nil // 用户(原房主)变成游离状态
	room.Players--
}

//玩家退出
func ExitRoomUser(userLine *UserLine) {
	room := userLine.RoomLine
	for i, v := range room.MapUsers {
		if v == userLine.UserData {
			room.MapUsers[i] = nil
		}
	}
	room.Players--
	userLine.RoomLine = nil
}

//解散(关闭)房间
func CloseRoom(userLine *UserLine) {
	room := userLine.RoomLine
	if room.Players <= 1 { //仅有房主一人
		if room.RoomOwner == userLine.UserData { //确认是房主 有权限
			userLine.RoomLine = nil
			room.Close()
		} else {
			userLine.WriteMsg(msg.MSG_ROOM_NOTAUTH)
			//return false
		}
	} else {
		userLine.WriteMsg(msg.MSG_ROOM_NOTEMPTY)
		//return false
	}
}

//每局后调用 检查玩家是否还有“资格”留在房间 资金量
func CheckPlayerMoney(userLine *UserLine) bool {
	room := userLine.RoomLine
	if room.RoomData.BaseMoney <= userLine.UserData.Chips { //
		return true
	} else {
		return false
	}
}
