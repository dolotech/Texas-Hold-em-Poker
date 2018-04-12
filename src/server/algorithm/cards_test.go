package algorithm

import (
	"testing"
	"strings"
)

func TestFour(t *testing.T) {

	cards := Cards{0x12, 0x02, 0x22, 0xb, 0x1b, 0x7, 0x32}
	//cards := Cards([]byte{0x12, 0x02, 0x22, 0xb, 0x1e, 0x14, 0x32})

	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.Four() == Card(2))

	cards = Cards{0x12, 0x02, 0x22, 0x32}

	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.Four() == 2)
}

func TestThree(t *testing.T) {
	var cs Cards
	cards := Cards{0x12, 0x02, 0x22, 0x3a, 0x2a, 0x1a, 0x33}
	SortCards(cards, 0, int8(len(cards))-1)
	cs = cards.Three()
	t.Log(cs.Equal(Cards{0xa, 0xa, 0xa, 0x3, 0x2}), cs.String())
}

func TestHighCard(t *testing.T) {
	var cs Cards
	cards := Cards{0x1E, 0x02, 0x22, 0x3a, 0x2a, 0x1a, 0x33}
	SortCards(cards, 0, int8(len(cards))-1)
	cs = cards.HighCard()
	t.Log(cs.Equal(Cards{0xe, 0xa, 0xa, 0xa, 0x3}), cs.String())
}

func TestCards_OnePair(t *testing.T) {
	var cs Cards
	cards := Cards{0x12, 0x0e, 0x2e, 0x3a, 0x2a, 0x1a, 0x32}
	SortCards(cards, 0, int8(len(cards))-1)

	cs = cards.OnePair()
	t.Log(cs.Equal(Cards{0xe, 0xe, 0xa, 0xa, 0xa}), cs.String())
}

func TestCards_Straight(t *testing.T) {

	cards := Cards{0x12, 0x03, 0x24, 0x35, 0x26, 0x17, 0x33}
	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.Straight() == 0x7)

	cards = Cards{0x12, 0x03, 0x24, 0x34, 0x26, 0x17, 0x37}
	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.Straight() == 0x0)

	cards = Cards{0x12, 0x03, 0x24, 0x34, 0x25, 0x17, 0x3E}
	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.Straight() == 0x5)
}

func TestCards_Flush(t *testing.T) {

	cards := Cards{0x32, 0x33, 0x34, 0x35, 0x26, 0x37, 0x28}

	SortCards(cards, 0, int8(len(cards))-1)
	t.Logf("%#v %#v ", len(cards.Flush()) == 5, cards.Flush())
	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.Flush())
}

func TestCards_TwoPair(t *testing.T) {

	cards := Cards{0x32, 0x32, 0x34, 0x35, 0x25, 0x38, 0x28}

	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(len(cards.TwoPair()) == 2, cards.TwoPair())
}
func TestCards_RoyalFlush(t *testing.T) {

	cards := Cards{0x3A, 0x3B, 0x3C, 0x3E, 0x3D, 0x38, 0x28}

	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.RoyalFlush() == true)

	cards = Cards{0x3A, 0x3B, 0x3C, 0x3E, 0x2D, 0x38, 0x28}

	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.RoyalFlush() == false)
}

func TestCards_PK(t *testing.T) {

	cards1 := Cards{0x32, 0x32, 0x34, 0x35, 0x25, 0x38, 0x28}
	cards2 := Cards{0x32, 0x32, 0x33, 0x34, 0x25, 0x36}

	t.Log(cards1.PK(&cards2) == -1)
}
func TestCards_StraightFlush(t *testing.T) {
	cards := Cards{0x32, 0x33, 0x34, 0x35, 0x36, 0x27, 0x28}
	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.StraightFlush() == 6)

	cards = Cards{0x32, 0x33, 0x34, 0x35, 0x3E, 0x37, 0x28}
	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.StraightFlush() == 5)

	cards = Cards{0x32, 0x33, 0x34, 0x35, 0x3E, 0x36, 0x28}
	SortCards(cards, 0, int8(len(cards))-1)
	t.Log(cards.StraightFlush() == 6)
}

func TestCards_FullFouse(t *testing.T) {
	var cs Cards
	cards := Cards{0x33, 0x33, 0x33, 0x35, 0x25, 0x35, 0x28}
	SortCards(cards, 0, int8(len(cards))-1)
	cs = cards.FullFouse()
	t.Log(cs.Equal(Cards{0x5, 0x5, 0x5, 0x03, 0x03}), cs.Hex())

	t.Log(cards.String())
	t.Log(cards.Hex())
}

func TestString2Num(t *testing.T) {

	array := Cards{0x33, 0x33, 0x33, 0x35, 0x25, 0x35, 0x3E, 0x1A}

	//for i := 0; i < N; i++ {
	go array.Shuffle()
	t.Log(array.String())
	//}

}

func TestFullFouse(t *testing.T) {
	var s = "A A A K K|" +
		"A A A Q Q|" +
		"A A A J J|" +
		"A A A T T|" +
		"A A A 9 9|" +
		"A A A 8 8|" +
		"A A A 7 7|" +
		"A A A 6 6|" +
		"A A A 5 5|" +
		"A A A 4 4|" +
		"A A A 3 3|" +
		"A A A 2 2|" +
		"K K K A A|" +
		"K K K Q Q|" +
		"K K K J J|" +
		"K K K T T|" +
		"K K K 9 9|" +
		"K K K 8 8|" +
		"K K K 7 7|" +
		"K K K 6 6|" +
		"K K K 5 5|" +
		"K K K 4 4|" +
		"K K K 3 3|" +
		"K K K 2 2|" +
		"Q Q Q A A|" +
		"Q Q Q K K|" +
		"Q Q Q J J|" +
		"Q Q Q T T|" +
		"Q Q Q 9 9"

	array := strings.Split(s, "|")
	for _, v := range array {
		cards := &Cards{}
		cards.SetByString(v)
		t.Log(cards.String(), len(cards.FullFouse()) > 0)
	}
}

func Test_Straight1(t *testing.T) {

	testCards := "T J Q K A|" +
		"9 T J Q K|" +
		"8 9 T J Q|" +
		"7 8 9 T J|" +
		"6 7 8 9 T|" +
		"5 6 7 8 9|" +
		"4 5 6 7 8|" +
		"3 4 5 6 7|" +
		"2 3 4 5 6|" +
		"A 2 3 4 5"

	array := strings.Split(testCards, "|")
	for _, v := range array {
		cards := &Cards{}
		cards.SetByString(v)
		cards.Sort()
		t.Log(cards.String(), cards.Straight() > 0)
	}
	t.Log("--------------------------------------")

	for _, v := range array {
		cards := &Cards{}
		cards.SetByString(v)
		cards.Sort()
		t.Log(cards.String(), cards.StraightFlush() > 0)
	}
	/*	t.Log("--------------------------------------")

		for _, v := range array {
			cards := &Cards{}
			cards.SetByString(v)
			cards.Sort()
			t.Log(cards.String(), cards.RoyalFlush())
		}*/
}
func Test_Append(t *testing.T) {
	cards:=Cards{0x33, 0x35, 0x25, 0x35, 0x28}
	cards = cards.Append(Cards{0x33, 0x33})
	cs, hand := cards.GetType()
	t.Logf("%#v %v  %#v",cs, hand,cards)
}
