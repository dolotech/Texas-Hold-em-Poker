package msg

import (
	//"github.com/dolotech/leaf/network/protobuf"
	"github.com/dolotech/leaf/network/json"
)

//注册消息内容 即 类型(结构体)
//var Processor =  protobuf.NewProcessor()
var Processor = json.NewProcessor()

var (
	// 用户数据
	MSG_SUCCESS          = &CodeState{Code: 0, Message: "success"}       //注册成功
	MSG_Register_Existed = &CodeState{Code: 1, Message: "existed user"}  //注册用户已存在
	MSG_Login_Error      = &CodeState{Code: 2, Message: "login fail"}    //登录失败 信息错误
	MSG_Version_Error    = &CodeState{Code: 3, Message: "version wrong"} //版本号不对
	MSG_DB_Error         = &CodeState{Code: 111, Message: "db error"}    //数据库出错

	//房间信息 1000开始标记
	MSG_ROOM_NOTAUTH    = &CodeState{Code: 1001, Message: "Unauthorized"}     //没有权限
	MSG_ROOM_ERRORPWD   = &CodeState{Code: 1002, Message: "pwd wrong"}        //密码错误
	MSG_ROOM_OVERVOLUME = &CodeState{Code: 1003, Message: "aleady in room"}   //你已经在其他房间了 拒绝加入其他房间
	MSG_ROOM_NOMONEY    = &CodeState{Code: 1004, Message: "not enough money"} //起始资金不够
	MSG_ROOM_NOTEMPTY   = &CodeState{Code: 1005, Message: "room not empty"}   //房子不空
	MSG_ROOM_NOROOM     = &CodeState{Code: 1006, Message: "no room"}          //没有该房子记录
)

func init() {
	Processor.Register(&Hello{})
	Processor.Register(&UserLoginInfo{})
	Processor.Register(&LoginError{})

	Processor.Register(&RegisterUserInfo{})

	Processor.Register(&CodeState{})
	Processor.Register(&Version{})

	//房间会话注册
	Processor.Register(&RoomInfo{})     //基本信息
	Processor.Register(&JoinRoom{}) //用户输入密码 点击进入
}


type Version struct {
	Version string // 版本号
}


type CodeState struct {
	Code    int    // const
	Message string //警告信息
}

type Hello struct {
	Name string
}

type UserLoginInfo struct {
	//登录
	Name string
	Pwd  string
}

type LoginError struct {
	State   int
	Message string
}

type RegisterUserInfo struct {
	//注册
	Name  string
	Pwd   string
	Age   int
	Email string
}

type RoomInfo struct {
	RoomName   string
	Volume     uint8
	GameType   uint32    //游戏类型 即玩法
	PayValue   uint8    //倍数
	BaseMoney  uint32    //最低资本 才能进房间
	RoomPwd    string //房间锁--密码
	RoomID     uint32
	RoomNumber string
}

type JoinRoom struct {
	RoomNumber string
	RoomPwd    string
}

type Bet struct {
	RoomNumber string
	RoomPwd    string
}

type LeaveRoom struct {
	RoomNumber string
	RoomPwd    string
}

type RoomPWDJoinCondition struct {
	Pwd string
}



