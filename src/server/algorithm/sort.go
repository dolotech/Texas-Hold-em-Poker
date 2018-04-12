package algorithm



// 对牌值从小到大排序，采用快速排序算法
func SortCards(arr Cards, start, end int8) {
	if start < end {
		i, j := start, end
		card := arr[(start+end)/2]
		key := card &0xF
		suit := card >> 4
		for i <= j {
			for (arr[i]) &0xF < key || ((arr[i]) &0xF == key && arr[i]>>4 < suit) {
				i++
			}
			for (arr[j]) &0xF > key || ((arr[j]) &0xF == key && arr[j]>>4 > suit) {

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