package internal

import "testing"

func Test_Pot(t *testing.T)  {
	bets := []uint32{60,80,90,0,0,0,0,0,0}
	res:=calcPot(bets)
	//[{180 [1 2 3]} {40 [2 3]} {10 [3]}]
	t.Log(res)


	bets = []uint32{60,60,60,0,0,0,0,0,0}
	res=calcPot(bets)
	//[{180 [1 2 3]}]
	t.Log(res)

	a:= []byte{}

	a = nil
	t.Log(len(a))

}
