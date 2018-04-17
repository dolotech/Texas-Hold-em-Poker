package internal

import (
	"github.com/golang/glog"
	"server/msg"
	"server/model"
	"server/algorithm"
)

func (r *Room) start() {
	// 2人及以上才开始游戏
	if r.n < 2 {
		return
	}
	// 产生庄
	var dealer *Occupant
	button := r.Button - 1
	r.Each((button+1)%r.Cap(), func(o *Occupant) bool {
		r.Button = o.Pos
		dealer = o
		return false
	})

	if dealer == nil {
		return
	}

	r.remain = 0

	// 剔除筹码小于大盲和离线的玩家
	r.Each(0, func(o *Occupant) bool {
		if o.Chips < r.BB || o.IsOffline() {
			r.removeOccupant(o)
			r.addObserve(o)
			return true
		}
		o.SetGameing()
		return true
	})

	// 洗牌
	r.Cards.Shuffle()

	// 产生小盲
	sb := r.next(dealer.Pos)
	if r.n == 2 { // one-to-one
		sb = dealer
	}
	// 产生大盲
	bb := r.next(sb.Pos)
	bbPos := bb.Pos

	// 通报本局庄家
	r.WriteMsg(&msg.Button{Uid: dealer.Uid})

	// 小大盲下注
	r.betting(sb, int32(r.SB))
	r.betting(bb, int32(r.BB))

	// Round 1 : preflop
	r.ready()
	r.Each(0, func(o *Occupant) bool {
		o.cards = algorithm.Cards{r.Cards.Take(), r.Cards.Take()}

		kindCards, kind := o.cards.GetType()
		m := &msg.PreFlop{
			Cards:     o.cards.Bytes(),
			Kind:      kind,
			KindCards: kindCards.Bytes(),
		}
		o.WriteMsg(m)
		return true
	})
	r.Broadcast(&msg.PreFlop{}, false)

	r.action(0)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

	// Round 2 : Flop
	r.ready()
	r.Cards = algorithm.Cards{r.Cards.Take(), r.Cards.Take(), r.Cards.Take()}
	r.Each(0, func(o *Occupant) bool {
		cs := r.Cards.Append(o.cards...)
		kindCards, kind := cs.GetType()
		m := &msg.Flop{
			Cards:     cs.Bytes(),
			Kind:      kind,
			KindCards: kindCards.Bytes(),
		}
		o.WriteMsg(m)
		return true
	})
	r.Broadcast(&msg.Flop{Cards: r.Cards.Bytes()}, false)

	r.action(0)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

	// Round 3 : Turn
	r.ready()
	r.Cards = r.Cards.Append(r.Cards.Take())
	r.Each(0, func(o *Occupant) bool {
		cs := r.Cards.Append(o.cards...)
		kindCards, kind := cs.GetType()
		m := &msg.Turn{
			Cards:     r.Cards[3].Byte(),
			Kind:      kind,
			KindCards: kindCards.Bytes(),
		}
		o.WriteMsg(m)
		return true
	})
	r.Broadcast(&msg.Turn{Cards: r.Cards[3].Byte()}, false)

	r.action(0)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

	// Round 4 : River
	r.ready()
	r.Cards = r.Cards.Append(r.Cards.Take())
	r.Each(0, func(o *Occupant) bool {

		cs := r.Cards.Append(o.cards...)
		kindCards, kind := cs.GetType()
		m := &msg.River{
			Cards:     r.Cards[4].Byte(),
			Kind:      kind,
			KindCards: kindCards.Bytes(),
		}
		o.WriteMsg(m)

		o.kindCards = kindCards
		o.kind = kind
		return true
	})
	r.Broadcast(&msg.River{Cards: r.Cards[4].Byte()}, false)

	r.action(0)

showdown:
	r.showdown()
	// Final : Showdown
	showdown := &msg.Showdown{}
	for _, o := range r.occupants {
		if o != nil && o.IsGameing() {
			item := &msg.ShowdownItem{
				Uid:      o.Uid,
				ChipsWin: r.Chips[o.Pos-1],
				Chips:    o.Chips,
			}
			showdown.Showdown = append(showdown.Showdown, item)
		}
	}

	r.Broadcast(showdown, true)
	glog.Infoln(sb.Pos, bbPos)
}

// 计算并通报奖池
func (r *Room) calc() (pots []handPot) {
	pots = calcPot(r.Chips)
	r.Pot = r.Pot[:]
	var ps []uint32
	for _, pot := range pots {
		r.Pot = append(r.Pot, pot.Pot)
		ps = append(ps, pot.Pot)
	}
	r.Broadcast(&msg.Pot{Pot: ps}, true)
	return
}

func (r *Room) action(pos uint8) {
	if r.allin+1 >= r.remain {
		return
	}
	var skip uint8
	if pos == 0 { // start from left hand of button
		pos = (r.Button)%r.Cap() + 1
	}
	for {
		var raised uint8
		r.Each(pos-1, func(o *Occupant) bool {
			if r.remain <= 1 {
				return false
			}
			if o.Pos == skip || o.Chips == 0 {
				return true
			}
			r.WriteMsg(&msg.BetPrompt{})
			n := o.GetAction(r.Timeout)
			if r.remain <= 1 {
				return false
			}
			if r.betting(o, n) {
				raised = o.Pos
				return false
			}
			return true
		})
		if raised == 0 {
			break
		}
		pos = raised
		skip = pos
	}
}

func (r *Room) ready() {
	r.Bet = 0
	r.Each(0, func(o *Occupant) bool {
		o.Bet = 0
		o.actionName = ""
		r.remain++
		o.kind = 0
		o.kindCards = nil
		return true
	})
}

// 比牌
func (r *Room) showdown() {
	pots := r.calc()

	for i, _ := range r.Chips {
		r.Chips[i] = 0
	}

	for _, pot := range pots {
		var maxO *Occupant
		for _, pos := range pot.OPos {
			o := r.occupants[pos-1]
			if o != nil && len(o.cards) > 0 {
				if maxO == nil {
					maxO = o
					continue
				}
				if maxO.PK(o) > 0 {
					maxO = o
				}
			}
		}

		var winners []uint8

		for _, pos := range pot.OPos {
			o := r.occupants[pos-1]
			if o != nil && maxO.PK(o) == 0 && o.IsGameing() {
				winners = append(winners, o.Pos)
			}
		}

		if len(winners) == 0 {
			glog.Errorln("!!!no winners!!!")
			return
		}

		for _, winner := range winners {
			r.Chips[winner-1] += pot.Pot / uint32(len(winners))
		}
		r.Chips[winners[0]-1] += pot.Pot % uint32(len(winners)) // odd chips
	}

	for i, _ := range r.Chips {
		if r.occupants[i] != nil {
			r.occupants[i].Chips += r.Chips[i]
		}
	}
}

func (r *Room) betting(o *Occupant, n int32) (raised bool) {
	value := n
	if n < 0 {
		o.actionName = model.BET_FOLD
		n = 0
		r.remain--
	} else if n == 0 {
		o.actionName = model.BET_CHECK
	} else if uint32(n)+o.Bet <= r.Bet {
		o.actionName = model.BET_CALL
		o.Chips -= uint32(n)
		o.Bet += uint32(n)
	} else {
		o.actionName = model.BET_RAISE
		o.Chips -= uint32(n)
		o.Bet += uint32(n)
		r.Bet = o.Bet
		raised = true
	}
	if o.Chips == 0 {
		o.actionName = model.BET_ALLIN
	}
	r.Chips[o.Pos-1] += uint32(n)

	r.Broadcast(&msg.BetBroadcast{
		Uid:   o.Uid,
		Kind:  o.actionName,
		Value: value,
	}, true)

	return
}

func (r *Room) next(pos uint8) *Occupant {
	volume := r.Cap()
	for i := (pos) % volume; i != pos-1; i = (i + 1) % volume {
		if r.occupants[i] != nil && r.occupants[i].IsGameing() {
			return r.occupants[i]
		}
	}
	return nil
}
