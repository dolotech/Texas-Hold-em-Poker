package algorithm

func PK(value1, value2 uint32) int8 {
	kind, value := De(value1)
	kinda, valuea := De(value2)

	if kind < kinda {
		return 1
	} else if kind == kinda {
		if kind == ROYAL_FLUSH {
			return 0
		}

		if value < valuea {
			return 1
		} else if value > valuea {
			return 1
		}
	}
	return 0
}

func (this *Cards) Counter() *ValueCounter {
	var counter ValueCounter
	counter.Set(*this)
	return &counter
}
func (this *Cards) GetType() uint32 {
	if len(*this) == 0 {
		return 0
	}

	counter := this.Counter()
	ASort(*this, 0, int8(len(*this))-1, counter)

	if res := this.RoyalFlush(); res > 0 {
		return res
	}

	if res := this.StraightFlush(); res > 0 {
		return res
	}

	if res := this.Four(counter); res > 0 {
		return res
	}

	if res := this.FullFouse(counter); res > 0 {
		return res
	}

	if res := this.Flush(); res > 0 {
		return res
	}

	if res := this.Straight(); res > 0 {
		return res
	}
	if res := this.Three(counter); res > 0 {
		return res
	}
	if res := this.TwoPair(); res > 0 {
		return res
	}

	if res := this.OnePair(); res > 0 {
		return res
	}

	return this.HighCard()

}
