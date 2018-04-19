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
	//高牌（high card）
	//既不是同一花色也不是同一点数的五张牌组成。
	//平手牌：如果不止一人抓到此牌，则比较点数最大者，
	//如果点数最大的相同，则比较第二、第三、第四和第五大的，如果所有牌都相同，则平分彩池。
	return En(HIGH_CARD, ToValue(*this))
}
