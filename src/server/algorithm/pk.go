package algorithm

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

	if res := this.royalFlush(); res > 0 {
		return res
	}

	if res := this.straightFlush(); res > 0 {
		return res
	}

	if res := this.four(counter); res > 0 {
		return res
	}

	if res := this.fullFouse(counter); res > 0 {
		return res
	}

	if res := this.flush(); res > 0 {
		return res
	}

	if res := this.straight(); res > 0 {
		return res
	}
	if res := this.three(counter); res > 0 {
		return res
	}
	if res := this.twoPair(); res > 0 {
		return res
	}

	if res := this.onePair(); res > 0 {
		return res
	}

	return this.highCard()

}
