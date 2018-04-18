package algorithm

import "testing"

func Test11Base(t *testing.T) {

	cards := Cards{0x12, 0x03, 0x24, 0x35, 0x26, 0x17, 0x33}
	var a ValueCounter
	a.Set(cards)
	ASort(cards, 0, int8(len(cards) -1), &a)

	t.Logf("%#v ",cards)
}
func Test10Base(t *testing.T) {
	arr := Cards([]byte{1, 2, 3, 4, 4, 5, 6})

	arr.Shuffle()

	t.Logf("%#v ", arr)

	t.Logf("%#v ", arr)
}
