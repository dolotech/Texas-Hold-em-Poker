package internal

import (
	"github.com/golang/glog"
	"server/msg"
	"strconv"
	"server/model"
	"server/algorithm"
)

func (r *Room) start() {
	// 2人及以上才开始游戏
	if r.n < 2 {
		return
	}
	// 产生庄位
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

	// 产生小盲位
	sb := r.next(dealer.Pos)
	if r.n == 2 { // one-to-one
		sb = dealer
	}
	// 产生大盲位
	bb := r.next(sb.Pos)
	bbPos := bb.Pos

	// 通报本局庄家
	r.WriteMsg(&msg.Button{Uid: dealer.Uid})

	// 小大盲下注
	r.betting(sb.Pos, int32(r.SB))
	r.betting(bb.Pos, int32(r.BB))

	// Round 1 : preflop
	r.Each(0, func(o *Occupant) bool {
		if o.IsGameing() {
			o.Bet = 0
			o.actionName = ""
			r.remain++
			o.cards = algorithm.Cards{r.Cards.Take(), r.Cards.Take()}
			m := &msg.PreFlop{}
			m.Cards = o.cards.Bytes()
			o.WriteMsg(m)
		}
		return true
	})

	r.action(0)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

	// Round 2 : Flop
	r.ready()
	r.Each(0, func(o *Occupant) bool {
		if o.IsGameing() {
			o.Bet = 0
			o.actionName = ""
			r.remain++
			r.Cards = algorithm.Cards{r.Cards.Take(), r.Cards.Take(), r.Cards.Take()}
			cs := r.Cards.Append(o.cards...)
			kindCards, kind := cs.GetType()
			m := &msg.Flop{
				Cards:     cs.Bytes(),
				Kind:      kind,
				KindCards: kindCards.Bytes(),
			}
			o.WriteMsg(m)
		}
		return true
	})

	r.action(0)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

	// Round 3 : Turn
	r.ready()
	r.Each(0, func(o *Occupant) bool {
		if o.IsGameing() {
			o.Bet = 0
			o.actionName = ""
			r.remain++
			card := r.Cards.Take()

			r.Cards = r.Cards.Append(card)

			cs := r.Cards.Append(o.cards...)
			kindCards, kind := cs.GetType()
			m := &msg.Turn{
				Cards:     card.Byte(),
				Kind:      kind,
				KindCards: kindCards.Bytes(),
			}
			o.WriteMsg(m)
		}
		return true
	})

	r.action(0)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

	// Round 4 : River

	r.ready()
	r.Each(0, func(o *Occupant) bool {
		if o.IsGameing() {
			o.Bet = 0
			o.actionName = ""
			r.remain++
			card := r.Cards.Take()

			r.Cards = r.Cards.Append(card)

			cs := r.Cards.Append(o.cards...)
			kindCards, kind := cs.GetType()
			m := &msg.Turn{
				Cards:     card.Byte(),
				Kind:      kind,
				KindCards: kindCards.Bytes(),
			}
			o.WriteMsg(m)
		}
		return true
	})

	r.action(0)

showdown:
	r.showdown()
	// Final : Showdown

	glog.Infoln(sb.Pos, bbPos)
}

func (r *Room) calc() (pots []handPot) {
	pots = calcPot(r.Chips)
	r.Pot = nil
	var ps []string
	for _, pot := range pots {
		r.Pot = append(r.Pot, pot.Pot)
		ps = append(ps, strconv.Itoa(int(pot.Pot)))
	}

	r.WriteMsg(&msg.Pot{
	})

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

			if o.Pos == skip || o.Chips == 0 || o.cards.Len() == 0 {
				return true
			}

			r.WriteMsg(&msg.Bet{Uid: o.Uid})

			n := o.GetAction(r.Timeout)
			if r.remain <= 1 {
				return false
			}

			if r.betting(o.Pos, n) {
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
		if o != nil {
			o.Bet = 0
		}
		return true
	})
}

func (r *Room) showdown() {
	//pots := r.calc()

	for i, _ := range r.Chips {
		r.Chips[i] = 0
	}
	/*

	for _, pot := range pots {
		var maxHand uint32
		for _, pos := range pot.OPos {
			o := r.occupants[pos-1]
			if o != nil && o.Hand > maxHand {
				maxHand = o.Hand
			}
		}

		var winners []int

		for _, pos := range pot.OPos {
			o := r.occupants[pos-1]
			if o != nil && o.Hand == maxHand && len(o.cards) > 0 {
				winners = append(winners, pos)
			}
		}

		if len(winners) == 0 {
			fmt.Println("!!!no winners!!!")
			return
		}

		for _, winner := range winners {
			r.Chips[winner-1] += pot.Pot / len(winners)
		}
		r.Chips[winners[0]-1] += pot.Pot % len(winners) // odd chips
	}
	*/

	for i, _ := range r.Chips {
		if r.occupants[i] != nil {
			r.occupants[i].Chips += r.Chips[i]
		}
	}
}

func (r *Room) betting(pos uint8, n int32) (raised bool) {
	if pos <= 0 {
		return
	}

	o := r.occupants[pos-1]
	if o == nil {
		return
	}

	value := n
	if n < 0 {
		o.actionName = model.BET_FOLD
		o.cards = nil
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

	r.WriteMsg(&msg.Bet{
		Uid:   o.Uid,
		Kind:  o.actionName,
		Value: value,
	})

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
