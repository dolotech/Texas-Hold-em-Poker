package internal

import (
	"github.com/golang/glog"
	"server/msg"
	"strconv"
)

func (r *Room) start() {
	if r.n < 2 {
		return
	}

	var dealer *Occupant

	occupants := make([]*Occupant, r.Max)
	copy(occupants, r.occupants)

	each(occupants, 0, func(o *Occupant) bool {
		if o.Chips < r.BB {
			o.Standup()
		}
		return true
	})

	// Select Dealer
	button := r.Button - 1
	each(occupants, (button+1)%r.Cap(), func(o *Occupant) bool {
		r.Button = o.Pos
		dealer = o
		return false
		return true
	})

	if dealer == nil {
		return
	}

	r.Cards.Shuffle()

	// Small Blind
	sb := r.next(occupants,dealer.Pos)
	if r.n == 2 { // one-to-one
		sb = dealer
	}
	// Big Blind
	bb := r.next(occupants ,sb.Pos)
	bbPos := bb.Pos

	each(occupants, 0, func(o *Occupant) bool {
		o.Bet = 0
		return true
	})

	r.WriteMsg(&msg.Button{Uid: dealer.Uid})

	each(occupants, 0, func(o *Occupant) bool {
		c1 := r.Cards.Take()
		c2 := r.Cards.Take()

		m := &msg.PreFlop{[]byte{c1.Byte(), c2.Byte()}}
		o.WriteMsg(m)
		return true
	})

	r.ready(occupants)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

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

func (r *Room) ready(occupants []*Occupant) {
	r.Bet = 0
	each(occupants, 0, func(o *Occupant) bool {
		o.Bet = 0
		return true
	})
}

func (r *Room) showdown() {
	r.WriteMsg(&msg.Showdown{

	})
}

func (r *Room) betting() {

}

func (r *Room) next(occupants []*Occupant,pos uint8) *Occupant {
	volume := r.Cap()
	for i := (pos) % volume; i != pos-1; i = (i + 1) % volume {
		if r.occupants[i] != nil {
			return r.occupants[i]
		}
	}
	return nil
}
