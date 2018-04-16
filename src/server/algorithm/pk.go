package algorithm

func (this *Cards) PK(cards *Cards) int8 {
	cs, kind := this.GetType()
	csa, kinda := cards.GetType()

	if kind < kinda {
		return 1
	} else if kind == kinda {
		if kind == ROYAL_FLUSH {
			return 0
		}

		for i := 0; i < len(cs); i++ {
			if cs[i] < csa[i] {
				return 1
			}
			if cs[i] > csa[i] {
				return -1
			}
		}
		return 0
	}

	return -1
}




func (this *Cards) GetType() (Cards, uint8) {
	if len(*this) == 0 {
		return []Card{}, 0
	}

	SortCards(*this, 0, int8(len(*this))-1)

	if this.RoyalFlush() {
		return []Card{0xA, 0xB, 0xC, 0xD, 0xE}, ROYAL_FLUSH
	}

	if cards := this.StraightFlush(); cards > 0 {
		return []Card{cards}, STRAIGHT_FLUSH
	}

	if cards := this.Four(); cards > 0 {
		return []Card{cards}, FOUR
	}

	if cards := this.FullFouse(); len(cards) > 0 {
		return cards, FULL_HOUSE
	}

	if cards := this.Flush(); len(cards) > 0 {
		return cards, FLUSH
	}

	if cards := this.Straight(); cards > 0 {
		return []Card{cards}, STRAIGHT
	}
	if cards := this.Three(); len(cards)  > 0 {
		return cards, THREE
	}
	if cards := this.TwoPair(); len(cards) > 0 {
		return cards, TWO_PAIR
	}

	if cards := this.OnePair(); len(cards)  > 0 {
		return cards, ONE_PAIR
	}

	return this.HighCard(), HIGH_CARD

}