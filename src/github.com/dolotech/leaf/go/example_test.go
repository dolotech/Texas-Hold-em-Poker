package g

import (
	"testing"
)

func TestGo_Go(t *testing.T) {
	d := New(10)

	// go 1
	var res int
	d.Go(func() {
		t.Log("1 + 1 = ?")
		res = 1 + 1
	}, func() {
		t.Log(res)
	})

	d.Cb(<-d.ChanCb)

	// go 2
	d.Go(func() {
		t.Log("My name is ")
	}, func() {
		t.Log("Leaf")
	})

	d.Close()

	// Output:
	// 1 + 1 = ?
	// 2
	// My name is Leaf
}

/*

func TestGo_NewLinearContext(t *testing.T) {
	d := New(10)

	// parallel
	d.Go(func() {
		time.Sleep(time.Second / 2)
		t.Log("1")
	}, nil)
	d.Go(func() {
		t.Log("2")
	}, nil)

	d.Cb(<-d.ChanCb)
	d.Cb(<-d.ChanCb)

	// linear
	c := d.NewLinearContext()
	c.Go(func() {
		time.Sleep(time.Second / 2)
		t.Log("1")
	}, nil)
	c.Go(func() {
		t.Log("2")
	}, nil)

	d.Close()

	// Output:
	// 2
	// 1
	// 1
	// 2
}
*/
