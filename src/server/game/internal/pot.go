package internal

import (
	"sort"
)

type handBet struct {
	Pos uint8
	Bet uint32
}

type handBets []handBet

func (p handBets) Len() int {
	return len(p)
}

func (p handBets) Less(i, j int) bool {
	return p[i].Bet < p[j].Bet
}

func (p handBets) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type handPot struct {
	Pot  uint32
	OPos []uint32
}

//因为每个人手中的筹码量会不时的产生变化，当出现多个筹码量不等的玩家“全下(ALL-IN)”争夺底池的时候就会出现多个奖池，
//这时候底池将分出主池和边池，主池里的筹码可 由任何一位争夺者胜利后将其拿走，其后每个边池将分别由参与者按牌型大小分别拿走。
//举个例子，假如发到河牌后共有三个人争夺底池，三人分别用ABC代替，
//假设A玩家手中有60，B手中有80，C手中有100。三人ALL-IN后则先把每人的筹码划分成：
//A玩家60；B玩家60+20；C玩家60+20+20，然后分出底池：
//主池 由60乘以三个人(A+B+C)形成一个堆，总额为180；
//边池一 由20乘以剩下的两个人(B+C)多出的部分为一个堆，总额为40；
//边池二 由筹码最多者C的多出的部分20单独组成一个堆，总额为20。
//现在三个底池形成了，然后分别按照牌型的大小来决定谁拿走哪个底池……
//首先，主池可以由三人中的任何一人得到胜利都可以拿走；
// 但是边池一的争夺就只能在B与C之间产生，两人谁赢谁拿走，与A无关(就算A玩家的牌最大也只能拿走他参与的主池)。
//而最后剩下的边池二不论三人谁输谁赢都只能由C拿走，因为边池二里的筹码只有C自己参与了，与其他人无关。
//main pot and side-pot calculation
func calcPot(bets []uint32) (pots []handPot) {
	var obs []handBet
	for i, bet := range bets {
		if bet > 0 {
			obs = append(obs, handBet{Pos: uint8(i) + 1, Bet: bet})
		}
	}
	sort.Sort(handBets(obs))

	for i, ob := range obs {
		if ob.Bet > 0 {
			s := obs[i:]
			hpot := handPot{Pot: ob.Bet * uint32(len(s))}

			for j, _ := range s {
				s[j].Bet -= ob.Bet
				hpot.OPos = append(hpot.OPos, uint32(s[j].Pos))
			}
			pots = append(pots, hpot)
		}
	}
	return
}
