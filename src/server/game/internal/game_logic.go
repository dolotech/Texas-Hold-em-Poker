package internal

import "github.com/golang/glog"

func (r *Room) start() {
	if r.inGameN() < 2 {
		return
	}

	var dealer *Occupant

	r.Each(0, func(o *Occupant) bool {
		if o.Chips < r.BB {
			o.Standup()
		}
		return true
	})

	// Select Dealer
	button := r.Button - 1
	r.Each((button+1)%r.Cap(), func(o *Occupant) bool {
		if o.inGame() {
			r.Button = o.Pos
			dealer = o
			return false
		}
		return true
	})

	r.Cards.Shuffle()

	// Small Blind
	sb := r.next(dealer.Pos)
	if r.inGameN() == 2 { // one-to-one
		sb = dealer
	}
	// Big Blind
	bb := r.next(sb.Pos)
	bbPos := bb.Pos


	glog.Infoln(sb.Pos,bbPos)
}

func (r *Room) ready() {

}

func (r *Room) showdown() {

}

func (r *Room) betting() {

}

func (r *Room) next(pos uint8) *Occupant {
	volume := r.Cap()
	for i := (pos) % volume; i != pos-1; i = (i + 1) % volume {
		if r.occupants[i] != nil && r.occupants[i].inGame() {
			return r.occupants[i]
		}
	}
	return nil
}

func (r *Room) EachInGame(start uint8, f func(o *Occupant) bool) {
	volume := r.Cap()
	end := (volume + start - 1) % volume
	i := start
	for ; i != end; i = (i + 1) % volume {
		if r.occupants[i] != nil && r.occupants[i].inGame() && !f(r.occupants[i]) {
			return
		}
	}

	// end
	if r.occupants[i] != nil && r.occupants[i].inGame() {
		f(r.occupants[i])
	}
}
func (r *Room) inGameN() uint8 {
	var n uint8
	r.EachInGame(0, func(o *Occupant) bool {
		n ++
		return true
	})
	return n
}
