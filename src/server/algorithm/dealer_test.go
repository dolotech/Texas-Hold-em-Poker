package algorithm

import "testing"

func Test_Dealer(t *testing.T)  {
	d:= &Cards{}

	d.Shuffle()
	d.Shuffle()
	t.Logf("%#v ",d)

}
