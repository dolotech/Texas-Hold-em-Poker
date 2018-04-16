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

	r.Each(0, func(o *Occupant) bool {
		if o.Chips < r.BB || o.IsOffline() {
			r.removeOccupant(o)
			r.addObserve(o)
			return true
		}
		o.SetGameing()
		return true
	})

	// Select Dealer
	button := r.Button - 1
	r.Each((button+1)%r.Cap(), func(o *Occupant) bool {
		if o.IsGameing(){
			r.Button = o.Pos
			dealer = o
			return false
		}
		return true
	})

	if dealer == nil {
		return
	}

	r.Cards.Shuffle()

	// Small Blind
	sb := r.next(dealer.Pos)
	if r.n == 2 { // one-to-one
		sb = dealer
	}
	// Big Blind
	bb := r.next(sb.Pos)
	bbPos := bb.Pos


	r.WriteMsg(&msg.Button{Uid: dealer.Uid})

	r.Each(0, func(o *Occupant) bool {
		if o.IsGameing() {
			c1 := r.Cards.Take()
			c2 := r.Cards.Take()
			m := &msg.PreFlop{[]byte{c1.Byte(), c2.Byte()}}
			o.WriteMsg(m)
		}

		return true
	})

	r.ready()

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

func (r *Room) ready() {
	r.Bet = 0
	r.Each(0, func(o *Occupant) bool {
		if o !=nil{
			o.Bet = 0
		}
		return true
	})
}

func (r *Room) showdown() {
	r.WriteMsg(&msg.Showdown{

	})
}

func (r *Room) betting() {

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
