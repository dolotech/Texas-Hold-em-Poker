package internal


type CardData struct {
	CardId	int	//牌的ID 顺序
	CardType	int	//牌的类型(1--万子 2--同子 3--条子 4--其他 0--初始)
	TotalCount	int //牌的数量 4
	AvailableCardCount	int	//可用的牌数量 0-4
	State	int	//状态(0--可用状态(没有被摸牌 初始) 1--有拥有者(user) 2--游离牌(已打出))
	UserID	int //拥有者 当 State==1时
	OldUserId	int	//曾经拥有者   （谁打出的牌记录)	State == 2
	CardValue	int	//牌的具体值 如 1万
	CardPosition int	//牌的位置 随机获得
}

type Cards []*CardData

var cards Cards

func init()  {
	//InitCard(cards)
}

func (c *CardData)InitCardClearUserInfo()  {//初始牌 只需将用户信息清除 以及牌的位置 状态进行初始
	//c.CardType = 0
	//c.CardType
	c.State = 0
	c.UserID = 0
	c.OldUserId = 0
	c.CardPosition = 0
}

//生成的牌 顺序 万 同 条 其他 分别标记为:1 2 3 4

func (cs Cards)InitCard()  {//第一次 点击开始

	for i := 0; i < 108; i++{//只对1-9进行 基本牌初始

		cs[i] = &CardData{CardType:( i / 36) + 1,CardValue:(i % 9)+1}

	}
//	东、南、西、北、中、发、白 (具体值 108-136)操作初始 的值分别为 28 29 30 31 32 33 34
	for i := 108; i< len(cs); i++{
		cs[i] = &CardData{CardType:(i / 36) + 1,CardValue:(i / 4) + 1}
	}

}

//牌拥有者添加 发牌 或者摸牌 同时 修改状态
func (c *CardData)AddCardUserIDAndState(userId int)  {
	c.UserID = userId
	c.State = 1
}

//牌的拥有者信息 进行修改
//当曾经拥有者 添加  处理碰 杠时情况
func (c *CardData)ChangeCardUserID(userId int)  {//
	c.OldUserId = c.UserID
	c.UserID = userId
}

//牌状态修改
//变成游离状态 修改userid
func (c *CardData)ChangeCardStateByThrowAway()  {
	c.State = 2
	c.OldUserId = c.UserID
	c.UserID = 0
}


//为[]*CardData 排序
func (c Cards)Len() int {
	return  len(c)
}
func (c Cards)Less(i, j int) bool  {
	if c[i].CardType == c[j].CardValue{
		return c[i].CardValue < c[j].CardValue
	}else if c[i].CardType < c[j].CardType{
		return true
	}else {
		return false
	}
}
func (c Cards)Swap(i, j int)  {
	c[i], c[j] = c[j], c[i]
}

