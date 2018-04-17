package model

const (
	MaxN     = 10
	MaxLevel = 40
	Timeout  = 30
)

const (
	Agent_Login = "LoginAgent" // 登录
	Agent_New   = "NewAgent"   // 新建链接
	Agent_Close = "CloseAgent" //链接关闭
)

const (
	BET_  = "BET_"  // 等待下注
	BET_CALL  = "call"  //跟注：等于单注额 (call)
	BET_FOLD  = "fold"  //弃牌: <0 (fold)
	BET_CHECK = "check" //看注：= 0 表示看注 (check)
	BET_RAISE = "raise" //加注：大于单注额 (raise)
	BET_ALLIN = "allin" //全押：等于玩家手中所有筹码 (allin)
)
