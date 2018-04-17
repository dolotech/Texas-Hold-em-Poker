package game

import "testing"

func Test_test(t *testing.T) {
	//var StraightValue = []uint32{126, 252, 504, 1008, 2016, 4032, 8064, 16128, 32256, 64512}
	//var StraightValue = []uint32{15872, 7936, 3968, 1984, 992, 496, 248, 124, 62, 31}
	array := []byte{0x31, 0x32, 0x33, 0x34, 0x25}
	var handvalue uint32
	for _, v := range array {
		value := (v & 0xF)
		if value == 0xE{
			handvalue |= 1
		}
		handvalue |= (1 << (value -1))
	}

	t.Logf("%b %d", handvalue, handvalue)



	for i := uint16(1); i <= 10; i++ {
		var value uint16
		for j := i; j <= i+5; j++ {
			value |= (1 << j)
		}

		t.Log(value)
	}

}
