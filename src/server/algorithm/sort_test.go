package algorithm

import "testing"

func Test10Base(t *testing.T) {
	arr := Cards([]Card{1, 2, 3, 4, 4, 5, 6})


	arr.Shuffle()


	t.Logf("%#v ",arr)

	SortCards(arr,0, int8(len(arr))-1)

	t.Logf("%#v ",arr)
}