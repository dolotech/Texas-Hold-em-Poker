package algorithm

import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
)

func TestFour(t *testing.T) {
	cards := Cards{0x12, 0x02, 0x22, 0xb, 0x1b, 0x7, 0x32}
	k,_:=De(cards.GetType())
	assert.Equal(t, k, FOUR)

	cards = Cards{0x12, 0x02, 0x22, 0x32}
	k,_=De(cards.GetType())
	assert.Equal(t, k, FOUR)
}

func TestThree(t *testing.T) {
	cards := Cards{0x12, 0x02, 0x22, 0x3a, 0x2a, 0x1a, 0x33}
	k,_:=De(cards.GetType())
	assert.Equal(t, k, FULL_HOUSE)

}

func TestHighCard(t *testing.T) {
	//var cs Cards
	cards := Cards{0x1E, 0x02, 0x22, 0x3a, 0x2a, 0x1a, 0x33}
	k, _ := De(cards.GetType())
	assert.Equal(t, k, FULL_HOUSE)
	//assert.Equal(t, v, uint32(0xE))
}

func TestCards_OnePair(t *testing.T) {
	cards := Cards{0x12, 0x0e, 0x2e, 0x3a, 0x2a, 0x1a, 0x32}

	k, _ := De(cards.GetType())
	assert.Equal(t, k, FULL_HOUSE)
}

func TestCards_Straight(t *testing.T) {

	cards := Cards{0x12, 0x03, 0x24, 0x35, 0x26, 0x17, 0x33}
	k, v := De(cards.GetType())
	assert.Equal(t, k, STRAIGHT)
	assert.Equal(t, v, uint32(7))

	cards = Cards{0x12, 0x03, 0x24, 0x34, 0x26, 0x17, 0x37}
	assert.Equal(t, k, STRAIGHT)
	assert.Equal(t, v, uint32(7))

	cards = Cards{0x12, 0x03, 0x24, 0x34, 0x25, 0x17, 0x3E}

	k, v = De(cards.GetType())
	assert.Equal(t, k, STRAIGHT)
	assert.Equal(t, v, uint32(0x5))


	cards = Cards{0x12, 0x03, 0x2a, 0x3c, 0x2b, 0x1d, 0x3E}
	k, v = De(cards.GetType())
	assert.Equal(t, k, STRAIGHT)
	assert.Equal(t, v, uint32(0xe))
}

func TestCards_Flush(t *testing.T) {

	cards := Cards{0x32, 0x33, 0x34, 0x35, 0x26, 0x37, 0x28}

	k, _ := De(cards.GetType())
	assert.Equal(t, k,FLUSH)
	//assert.Equal(t, v, uint32(6))
}

func TestCards_TwoPair(t *testing.T) {

	cards := Cards{0x12, 0x22, 0x34, 0x35, 0x25, 0x38, 0x28}

	k, _ := De(cards.GetType())
	assert.Equal(t, k, TWO_PAIR)
}
func TestCards_RoyalFlush(t *testing.T) {

	cards := Cards{0x3A, 0x3B, 0x3C, 0x3E, 0x3D, 0x38, 0x28}
	k, _ := De(cards.GetType())
	assert.Equal(t, k, ROYAL_FLUSH)

	cards = Cards{0x3A, 0x3B, 0x3C, 0x3E, 0x2D, 0x38, 0x28}

	k, _ = De(cards.GetType())
	assert.Equal(t, k, FLUSH)
}


func TestCards_FLUSH1(t *testing.T) {
	cards1 := Cards{0x22, 0x22, 0x22, 0x25, 0x38, 0x28}
	cards2 := Cards{0x32, 0x32, 0x33, 0x34, 0x25, 0x36}
	v1:=cards1.GetType()
	v2:=cards2.GetType()

	assert.Equal(t,v1 > v2 ,true)


	cards1 = Cards{0x22, 0x22, 0x2E, 0x25, 0x38, 0x28}
	cards2 = Cards{0x32, 0x32, 0x33, 0x3a, 0x2a, 0x3E}
	assert.Equal(t,cards1.GetType() > cards2.GetType() ,false)
}

func TestCards_FLUSH(t *testing.T) {
	cards1 := Cards{0x22, 0x22, 0x24, 0x25, 0x38, 0x28}
	cards2 := Cards{0x32, 0x32, 0x33, 0x34, 0x25, 0x36}
	v1:=cards1.GetType()
	v2:=cards2.GetType()

	assert.Equal(t,v1 > v2 ,true)
}

func TestCards_PK(t *testing.T) {

	cards1 := Cards{0x32, 0x32, 0x34, 0x35, 0x25, 0x38, 0x28}
	cards2 := Cards{0x32, 0x32, 0x33, 0x34, 0x25, 0x36}
	v1:=cards1.GetType()
	v2:=cards2.GetType()
	assert.Equal(t,PK(v1,v2) ,int8(1))

}
func TestCards_StraightFlush(t *testing.T) {
	cards := Cards{0x32, 0x33, 0x34, 0x35, 0x36, 0x27, 0x28}
	k, v := De(cards.GetType())
	assert.Equal(t, k, uint8(9))
	assert.Equal(t, v, uint32(6))

	cards = Cards{0x32, 0x33, 0x34, 0x35, 0x3E, 0x37, 0x28}
	k, v = De(cards.GetType())
	assert.Equal(t, k, uint8(9))
	assert.Equal(t, v, uint32(5))

	cards = Cards{0x32, 0x33, 0x34, 0x35, 0x3E, 0x36, 0x28}
	k, v = De(cards.GetType())
	assert.Equal(t, k, uint8(9))
	assert.Equal(t, v, uint32(6))
}

func TestCards_FullFouse(t *testing.T) {

	cards := Cards{0x33, 0x33, 0x33, 0x35, 0x25, 0x35, 0x28}
	k, _ := De(cards.GetType())

	//t.Logf("%v %v %#v ",k,v,cards)
	assert.Equal(t, k, FULL_HOUSE)
	//assert.Equal(t, v, FULL_HOUSE)

	//t.Log(cards.String())
	//t.Log(cards.Hex())
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

	var oldValue uint32
	for _, v := range array {
		cards := &Cards{}
		cards.SetByString(v)
		//t.Log(cards.String())

		k, value := De(cards.GetType())
		if oldValue == 0{
			oldValue = value
			continue
		}

		assert.Equal(t, oldValue> value, true)

		assert.Equal(t, k, FULL_HOUSE)
		//assert.Equal(t,v,uint32(6))

		//t.Log(De(cards.FullFouse(cards.Counter())))
	}
}

func Test_Straight1(t *testing.T) {

	/*	testCards := "T J Q K A|" +
		"9 T J Q K|" +
		"8 9 T J Q|" +
		"7 8 9 T J|" +
		"6 7 8 9 T|" +
		"5 6 7 8 9|" +
		"4 5 6 7 8|" +
		"3 4 5 6 7|" +
		"2 3 4 5 6|" +
		"A 2 3 4 5"

	//array := strings.Split(testCards, "|")
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
	}*/
	/*	t.Log("--------------------------------------")

		for _, v := range array {
			cards := &Cards{}
			cards.SetByString(v)
			cards.Sort()
			t.Log(cards.String(), cards.RoyalFlush())
		}*/
}

func Test_AnalyseCards(t *testing.T) {
	var a ColorCounter
	cards := []byte{0x33, 0x35, 0x25, 0x35, 0x28}
	a.Set(cards)

	t.Log(a.Get(0x33), a.Get(0x35), a.Get(0x28))
}
func Test_Append(t *testing.T) {
	cards := Cards{0x33, 0x35, 0x25, 0x35, 0x28}
	cards = cards.Append(Cards{0x33, 0x33}...)

	t.Logf("%#v ", cards.GetType())
}

func Test_turnToValue1(t *testing.T) {
	v1 := []byte{0x33, 0x35, 0x25, 0x35, 0x28}
	v2 := []byte{0x35, 0x35, 0x25, 0x35, 0x27}

	t.Logf("%#v %#v ", v1, v2)
	b1 := ToValue(v1)
	b2 := ToValue(v2)

	t.Log(b1, b2)

}
func Test_ColorCounter(t *testing.T) {
	v2 := []byte{0x35, 0x35, 0x25, 0x35, 0x27}

	var colorCounter ColorCounter

	colorCounter.Set(v2)

	t.Log(v2)
}
func Test_turnToValue(t *testing.T) {

	v := ToValue([]byte{0x3E, 0x3E, 0x3E, 0x3E, 0x3E, 0x3E, 0x3E})

	//var v3 uint16 = v
	v1 := v | ( 10 << 24)

	//t.Log(v3)

	t.Logf("%v %v %v", v, v1>>24, v1&0xFFFFFF)
	t.Logf("%v ", ToValue([]byte{0x33, 0x35, 0x25, 0x35, 0x28}))

}
