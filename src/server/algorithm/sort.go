package algorithm

// 对牌值从小到大排序，采用快速排序算法
func SortCards(arr []byte, start, end int8) {
	if start < end {
		i, j := start, end
		card := arr[(start+end)/2]
		key := card & 0xF
		suit := card >> 4
		for i <= j {
			for (arr[i])&0xF < key || ((arr[i])&0xF == key && arr[i]>>4 < suit) {
				i++
			}
			for (arr[j])&0xF > key || ((arr[j])&0xF == key && arr[j]>>4 > suit) {

				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			SortCards(arr, start, j)
		}
		if end > i {
			SortCards(arr, i, end)
		}
	}
}

func Sort(cards []byte, start, end int8) {
	if start < end {
		i, j := start, end
		card := cards[(start+end)/2]
		for i <= j {
			for cards[i] < card {
				i++
			}
			for cards[j] > card {
				j--
			}
			if i <= j {
				cards[i], cards[j] = cards[j], cards[i]
				i++
				j--
			}
		}
		if start < j {
			Sort(cards, start, j)
		}
		if end > i {
			Sort(cards, i, end)
		}
	}
}

type ColorCounter uint16

func (this *ColorCounter) Set(cards []byte) {
	//*this = 0
	l := uint8(len(cards))
	for i := uint8(0); i < l; i++ {
		card := (cards[i] >> 4) * 3
		count := ((*this >> card) & 0x07) + 1
		*this &= (^(0x07 << card))
		*this |= (count << card)
	}
}

func (this *ColorCounter) Get(card byte) uint8 {
	return uint8((*this >> (card >> 4) * 3 ) & 0x07)
}

type ValueCounter uint64

func (this *ValueCounter) Set(cards []byte) {
	//*this = 0
	l := uint8(len(cards))
	for i := uint8(0); i < l; i++ {
		card := (cards[i] & 0xF) * 3
		count := ((*this >> card) & 0x07) + 1
		*this &= (^(0x07 << card))
		*this |= (count << card)
	}
}

func (this *ValueCounter) Get(card byte) uint8 {
	return uint8((*this >> (card & 0xF * 3 )) & 0x07)
}

func ASort(arr []byte, start, end int8, counter *ValueCounter) {
	if start < end {
		i, j := start, end
		card := arr[(start+end)/2]
		key := card & 0xF
		count := counter.Get(key)
		for i <= j {
			for (counter.Get(arr[i]&0xF) < count) || ((counter.Get(arr[i]&0xF) == count) && (arr[i])&0xF < key ) {
				i++
			}
			for (counter.Get(arr[j]&0xF) > count) || ((counter.Get(arr[j]&0xF) == count) && (arr[j])&0xF > key ) {

				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			ASort(arr, start, j, counter)
		}
		if end > i {
			ASort(arr, i, end, counter)
		}
	}
}
