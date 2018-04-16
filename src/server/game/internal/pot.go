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

// main pot and side-pot calculation
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
