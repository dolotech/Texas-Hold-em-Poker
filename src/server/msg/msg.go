package msg

import (
	//"github.com/dolotech/leaf/network/protobuf"
	"github.com/dolotech/leaf/network/json"
)

var Processor = json.NewProcessor()

var (
	// 用户数据
	MSG_SUCCESS          = &CodeState{Code: 0, Message: "success"}        //注册成功
	MSG_Register_Existed = &CodeState{Code: 1, Message: "existed user"}   //注册用户已存在
	MSG_Login_Error      = &CodeState{Code: 2, Message: "login fail"}     //登录失败 信息错误
	MSG_Version_Error    = &CodeState{Code: 3, Message: "version wrong"}  //版本号不对
	MSG_User_Not_Exist   = &CodeState{Code: 4, Message: "user not exist"} //用户不存在
	MSG_DB_Error         = &CodeState{Code: 111, Message: "db error"}     //数据库出错

	//房间错误信息 1000开始标记
	MSG_ROOM_NOTAUTH    = &CodeState{Code: 1001, Message: "Unauthorized"}     //没有权限
	MSG_ROOM_ERRORPWD   = &CodeState{Code: 1002, Message: "pwd wrong"}        //密码错误
	MSG_ROOM_OVERVOLUME = &CodeState{Code: 1003, Message: "aleady in room"}   //你已经在其他房间了 拒绝加入其他房间
	MSG_ROOM_NOMONEY    = &CodeState{Code: 1004, Message: "not enough money"} //起始资金不够
	MSG_ROOM_NOTEMPTY   = &CodeState{Code: 1005, Message: "room not empty"}   //房子不空
	MSG_ROOM_NOROOM     = &CodeState{Code: 1006, Message: "no room"}          //没有该房子记录
	MSG_ROOM_FULL       = &CodeState{Code: 1007, Message: "room full"}        // 房间已满
	MSG_NOT_IN_ROOM     = &CodeState{Code: 1008, Message: "not in room"}      // 你不在房间
	MSG_ROOM_CLOSED     = &CodeState{Code: 1009, Message: "room closed"}      // 房间已经关闭
)

func init() {
	Processor.Register(&Hello{})
	Processor.Register(&UserLoginInfo{})
	Processor.Register(&UserLoginInfoResp{})

	Processor.Register(&CodeState{})
	Processor.Register(&Version{})

	//房间会话注册
	Processor.Register(&RoomInfo{})  //基本信息
	Processor.Register(&JoinRoom{})  //
	Processor.Register(&LeaveRoom{}) //

	Processor.Register(&Showdown{})
	Processor.Register(&Deal{})
	Processor.Register(&Pot{})
	Processor.Register(&Bet{})
	Processor.Register(&Button{})
	Processor.Register(&StandUp{})
	Processor.Register(&SitDown{})
}

// 版本号
type Version struct {
	Version string
}

type CodeState struct {
	Code    int    // const
	Message string //警告信息
}

type Hello struct {
	Name string
}

//登录
type UserLoginInfo struct {
	UnionId  string
	Nickname string
}

//登录
type UserLoginInfoResp struct {
	UnionId  string
	Nickname string
	Account  string
}

type RoomInfo struct {
	Number     string
	Volume     uint8
	GameType   uint32 //游戏类型 即玩法
	PayValue   uint8  //倍数
	BaseMoney  uint32 //最低资本 才能进房间
	RoomPwd    string //房间锁--密码
	RoomID     uint32
	RoomNumber string
}

type StandUp struct {
	Uid uint32
}

type SitDown struct {
	Uid uint32
}
type JoinRoom struct {
	Uid        uint32
	RoomNumber string
	RoomPwd    string
}

// 发牌
type Deal struct {
	RoomNumber string
	RoomPwd    string
}

//通报本局庄家
type Button struct {
	RoomNumber string
	RoomPwd    string
}

//玩家下注
type Bet struct {
	RoomNumber string
	RoomPwd    string
}

//通报奖池
type Pot struct {
	RoomNumber string
	RoomPwd    string
}

//摊牌和比牌
type Showdown struct {
	RoomNumber string
	RoomPwd    string
}

type LeaveRoom struct {
	RoomNumber string
	Uid        uint32
}
