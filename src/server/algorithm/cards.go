package algorithm

type Cards []Card

func (this *Cards) Take() Card {
	card := (*this)[0]
	(*this) = (*this)[1:]
	return card
}
func (this *Cards) Append(cards Cards) Cards {
	cs := make([]Card,0,len(cards) + len(*this))
	cs = append(cs,(*this)...)
	cs = append(cs, cards...)
	return cs
}


// todo 两对和四张起脚牌的判定

//顺子（Straight，亦称“蛇”）
//此牌由五张顺序扑克牌组成。
//平手牌：如果不止一人抓到此牌，则五张牌中点数最大的赢得此局，
// 如果所有牌点数都相同，平分彩池。
func (this *Cards) Straight() Card {
	lenght := int8(len(*this))
	for i := lenght - 1; i > 2; i-- {
		icard := (*this)[i] & 0xF
		var last = icard
		var count int8
		for j := i - 1; j >= 0; j-- {
			jcard := (*this)[j] & 0xF
			if jcard+1 == last {
				last = jcard
				count ++
				if count == TYPE_LEN-1 {
					return icard
				} else if count == TYPE_LEN-2 && icard == 0x5 {
					for n := 0; n < SUITSIZE; n++ {
						index := int(lenght) - 1 - n
						tc := (*this)[index]
						if tc&0xF != 0xE {
							break
						}
						return icard
					}
				}
			}
		}
	}
	return 0
}

//同花顺（Straight Flush）
//五张同花色的连续牌。
//平手牌：如果摊牌时有两副或多副同花顺，连续牌的头张牌大的获得筹码。
//如果是两副或多副相同的连续牌，平分筹码。
func (this *Cards) StraightFlush() Card {
	lenght := int8(len(*this))
	for k := int8(0); k < SUITSIZE; k++ {
		for i := lenght - 1; i > 2; i-- {
			if int8((*this)[i]>>SUITSIZE) == k {
				icard := (*this)[i] & 0xF
				var last = icard
				var count int8
				for j := i - 1; j >= 0; j-- {
					jcard := (*this)[j] & 0xF
					if int8((*this)[j]>>SUITSIZE) != k {
						continue
					} else if jcard+1 != last {
						break
					}
					last = jcard
					count ++
					if count == TYPE_LEN-1 {
						return icard
					} else if count == TYPE_LEN-2 && icard == 0x5 {
						for n := 0; n < SUITSIZE; n++ {
							index := int(lenght) - 1 - n
							tc := (*this)[index]
							if tc&0xF != 0xE {
								break
							}
							if int8(tc>>SUITSIZE) == k {
								return icard
							}
						}
					}
				}
			}
		}
	}
	return 0
}

//皇家同花顺（Royal Flush）
//同花色的A, K, Q, J和10。
//平手牌：在摊牌的时候有两副多副皇家同花顺时，平分筹码。
func (this *Cards) RoyalFlush() bool {
	lenght := int8(len(*this))
	var last = (*this)[lenght-1]
	if last&0xF == 0xE {
		var count int8
		for i := lenght - 2; i >= 0; i-- {
			if last&0xF == (*this)[i]&0xF+1 {
				if last>>SUITSIZE == (*this)[i]>>SUITSIZE {
					count ++
					if count == TYPE_LEN-1 {
						return true
					}
					last = (*this)[i]
				}
			} else {
				break
			}
		}
	}
	return false
}

//四条（Four of a Kind，亦称“铁支”、“四张”或“炸弹”）
//其中四张是相同点数但不同花的扑克牌，第五张是随意的一张牌。
//平手牌：如果两组或者更多组摊牌，则四张牌中的最大者赢局，如果一组人持有的四张牌是一样的，
//那么第五张牌最大者赢局（起脚牌,2张起手牌中小的那张就叫做起脚牌）。如果起脚牌也一样，平分彩池。
func (this *Cards) Four() Card {
	lenght := int8(len(*this))
	for i := lenght - 1; i > 2; i-- {
		if (*this)[i-1]&0xF == (*this)[i]&0xF &&
			(*this)[i-2]&0xF == (*this)[i]&0xF &&
			(*this)[i-3]&0xF == (*this)[i]&0xF {
			return (*this)[i] & 0xF
		}
	}
	return 0
}

//满堂彩（Fullhouse，葫芦，三带二）
//由三张相同点数及任何两张其他相同点数的扑克牌组成。
//平手牌：如果两组或者更多组摊牌，那么三张相同点数中较大者赢局。
//如果三张牌都一样，则两张牌中点数较大者赢局，如果所有的牌都一样，则平分彩池。
func (this *Cards) FullFouse() Cards {
	lenght := int8(len(*this))
	for i := lenght - 1; i > 1; i-- {
		icard := (*this)[i] & 0xF
		if icard == (*this)[i-1]&0xF &&
			icard == (*this)[i-2]&0xF {
			for j := int8(0); j < lenght-1; j++ {
				jcard := (*this)[j] & 0xF
				if icard != jcard {
					if jcard == (*this)[j+1]&0xF {
						return Cards{icard, icard, icard, jcard, jcard}
					}
				}
			}
		}
	}
	return Cards{}
}

//同花（Flush，简称“花”）
//此牌由五张不按顺序但相同花的扑克牌组成。
//平手牌：如果不止一人抓到此牌相，则牌点最高的人赢得该局，
//如果最大点相同，则由第二、第三、第四或者第五张牌来决定胜负，如果所有的牌都相同，平分彩池。
func (this *Cards) Flush() Cards {
	lenght := int8(len(*this))
	cards := make([]Card, 0, TYPE_LEN)
	for i := lenght - 1; i > 0; i-- {
		icard := (*this)[i] >> SUITSIZE
		cards = cards[:0]
		cards = append(cards, (*this)[i]&0xF)
		for j := i - 1; j >= 0; j-- {
			if (*this)[j]>>SUITSIZE == icard {
				cards = append(cards, (*this)[j]&0xF)
				if len(cards) == 5 {
					return cards
				}
			}
		}
	}

	return Cards{}
}

//三条（Three of a kind，亦称“三张”）
//由三张相同点数和两张不同点数的扑克组成。
//平手牌：如果不止一人抓到此牌，则三张牌中最大点数者赢局，
//如果三张牌都相同，比较第四张牌，必要时比较第五张，点数大的人赢局。如果所有牌都相同，则平分彩池。
func (this *Cards) Three() Cards {
	lenght := int8(len(*this))
	for i := lenght - 1; i > 1; i-- {
		icard := (*this)[i] & 0xF
		if (*this)[i-1]&0xF == icard &&
			(*this)[i-2]&0xF == icard {

			cards := make([]Card, 0, TYPE_LEN)
			cards = append(cards, icard, icard, icard)
			for j := lenght - 1; j >= 0; j-- {
				jcard := (*this)[j] & 0xF
				if jcard != icard {
					cards = append(cards, jcard)
					if len(cards) == TYPE_LEN {
						return cards
					}
				}
			}
		}
	}
	return Cards{}
}

//两对（Two Pairs）
//两对点数相同但两两不同的扑克和随意的一张牌组成。
//平手牌：如果不止一人抓大此牌相，牌点比较大的人赢，如果比较大的牌点相同，那么较小牌点中的较大者赢，
//如果两对牌点相同，那么第五张牌点较大者赢（起脚牌,2张起手牌中小的那张就叫做起脚牌）。如果起脚牌也相同，则平分彩池。
func (this *Cards) TwoPair() Cards {
	lenght := int8(len(*this))
	cards := make([]Card, 0, 2)
	for i := lenght - 1; i > 0; i-- {
		icard := (*this)[i] & 0xF
		for j := i - 1; j >= 0; j-- {
			jcard := (*this)[j] & 0xF
			if jcard == icard {
				cards = append(cards, jcard)
				if len(cards) == 2 {
					return cards
				}
			}
		}
	}

	return Cards{}
}

//一对（One Pair）
//由两张相同点数的扑克牌和另三张随意的牌组成。
//平手牌：如果不止一人抓到此牌，则两张牌中点数大的赢，如果对牌都一样，则比较另外三张牌中大的赢，
//如果另外三张牌中较大的也一样则比较第二大的和第三大的，如果所有的牌都一样，则平分彩池。
func (this *Cards) OnePair() Cards {
	lenght := int8(len(*this))
	for i := lenght - 1; i > 0; i-- {
		icard := (*this)[i] & 0xF
		if (*this)[i-1]&0xF == icard {
			cards := make([]Card, 0, TYPE_LEN)
			cards = append(cards, icard, icard)
			for j := lenght - 1; j >= 0; j-- {
				jcard := (*this)[j] & 0xF
				if jcard != icard {
					cards = append(cards, jcard)
					if len(cards) == TYPE_LEN {
						return cards
					}
				}
			}
		}
	}
	return Cards{}
}

//高牌（high card）
//既不是同一花色也不是同一点数的五张牌组成。
//平手牌：如果不止一人抓到此牌，则比较点数最大者，
//如果点数最大的相同，则比较第二、第三、第四和第五大的，如果所有牌都相同，则平分彩池。
func (this *Cards) HighCard() Cards {
	lenght := int8(len(*this))
	if lenght > 0 {
		l := lenght
		if l > TYPE_LEN {
			l = TYPE_LEN
		}
		cards := make([]Card, l)
		for i := int8(0); i < l; i++ {
			cards[i] = (*this)[lenght-1-i] & 0xF
		}
		return cards
	}
	return Cards{}
}
