package internal

import (
	//"leaf/gate"
	"github.com/dolotech/leaf/gate"
	"github.com/dolotech/leaf/go"
	"server/msg"
)

var userLine *UserLine

type UserLine	struct {
	gate.Agent //申请代理
	*g.LinearContext
	Cards	[]*CardData //牌
	UserData	*UserData
	RoomLine	*Room
	//RoomPosition	int	//房间里的座位号 1-4
	ReadyState	int	//准备状态 是否可以开局 0否 1是
}




//初始化 用户登录后 进行
func (u *UserLine)initUser(userInfo *UserData)  {
	u.UserData = userInfo
	u.RoomLine = nil
	//u.Cards = nil
	u.LinearContext = skeleton.NewLinearContext()
}

//点击准备 异或操作
func (u *UserLine)readyGame()  {
	u.ReadyState = u.ReadyState ^ 1
	if u.ReadyState >0{
		u.RoomLine.ReadySignCount ++
	}
}



/*
	yin
	房间操作
 */
//创建房间
func (u *UserLine)createdRoom(roomInfo *msg.RoomInfo)  {
	room := InitRoom(roomInfo,u)
	if room != nil{//创建成功
		//u.Go()
		u.RoomLine = room
		allRooms = append(allRooms,room) //添加一个新的房间信息
		u.Go(func() {//房间开始监听 看是否可以进行游戏
			for{
				room.gameStart()
			}
		}, func() {

		})
	}else {
		u.WriteMsg(&msg.CodeState{msg.MSG_ROOM_OVERVOLUME,"你已在其他房间！"})
	}
}

//查找房间 通过RoomNumber
func (u *UserLine)FindRoom(roomNumber string)(room *Room){
	room = FindRoomsByRoomNumber(roomNumber)
	if room != nil {//有该房间的信息
	//	返回房间的信息 向前台
		roomData := room.RoomData
		//roomInfo := msg.RoomInfo{roomData.RoomName,}
		u.WriteMsg(&msg.RoomInfo{
				roomData.RoomName,roomData.Volume,
				roomData.GameType,roomData.PayValue,
				roomData.BaseMoney,"",roomData.RoomID,
				roomData.RoomNumber})

	}else {
		u.WriteMsg(&msg.CodeState{msg.MSG_ROOM_NOROOM,"没有该 "+ roomNumber +" 的记录!"})
	}
	return
}


//加入房间
func (u *UserLine)joinRoom(room *Room,pwd string)  {
	if room.CheckPlayerAndPwd(pwd){
		if CheckConditionForJoining(u, room){
			if JoinRoom(u,room) {

			}else {
				u.WriteMsg(&msg.CodeState{msg.MSG_ROOM_OVERVOLUME,"你已在其他房间进行游戏了，请退出房间，后操作！"})
			}

		}else {
			u.WriteMsg(&msg.CodeState{msg.MSG_ROOM_NOMONEY,"你的起始资金不够支付该房间的费用!"})
		}

	}else {
		u.WriteMsg(&msg.CodeState{msg.MSG_ROOM_NOTAUTH,"房间密码输入有误!"})
	}
}

//退出房间
func (u *UserLine)exitRoom()  {
//	房主退出
	room := u.RoomLine
	if room.RoomOwner == u.UserData{
		ChangeRoomOlderToNewer(u)
	}else {
	//	非 房主退出
		ExitRoomUser(u)
	}
}

//修改房间信息
func (u *UserLine)changeRoomInfo(roomInfo *msg.RoomInfo)  {
	if ChangeRoomInfo(u, roomInfo){
		//	修改成功
	}else {
		userLine.WriteMsg(&msg.CodeState{msg.MSG_ROOM_NOTAUTH,"你不是房主"})
	}
}

//解散(关闭)房间
func (u *UserLine)closeRoom()  {
	CloseRoom(u)
}



