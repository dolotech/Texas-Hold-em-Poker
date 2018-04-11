package internal

import (
	"math/rand"
	"time"
)

const  (
	TOTALCARDS = 136 //牌总数
	INIT_PERSIONALCARDS = 13 //初始阶段 每个人13张牌的数量
)

func init()  {
	rand.Seed(time.Now().UnixNano())
}

//发牌 发13张(每个人) 开始游戏阶段  返回一个二维数组[4][13] 即4个用户的牌
func DealCards(mCards map[int]int)(u [4][]int)  {
	l := len(mCards)
	for i := 0; i < len(u); i++{
		for l < INIT_PERSIONALCARDS * ( i + 1){
			n := rand.Intn(TOTALCARDS)
			//mCards[n] = n
			if _, ok := mCards[n]; !ok{
				mCards[n] = n //记录已发过的牌
				u[i] = append(u[i],n)//用户的牌
			}
		}
	}
	return u
}

//摸牌 一张 返回下标
//参数 l 为已经分配过的牌的数量
func DrawCard(mCards map[int]int) int  {
	l := len(mCards)
	var i int
	for l < l + 1{
		i = rand.Intn(TOTALCARDS)
		mCards[i] = i
	}
	return i
}



