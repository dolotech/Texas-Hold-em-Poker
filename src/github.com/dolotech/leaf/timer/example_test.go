package timer

import (
	"time"
	"testing"
)

func TestTimer(t *testing.T) {
	d := NewDispatcher(10)

	// timer 1
	d.AfterFunc(1, func() {
		t.Log("My name is Leaf")
	})

	// timer 2
	tim := d.AfterFunc(1, func() {
		t.Log("will not print")
	})
	tim.Stop()

	// dispatch
	(<-d.ChanTimer).Cb()

	// Output:
	// My name is Leaf
}

func TestNewCronExpr(t *testing.T) {
	cronExpr, err := NewCronExpr("0 * * * *")
	if err != nil {
		return
	}

	t.Log(cronExpr.Next(time.Date(
		2000, 1, 1,
		20, 10, 5,
		0, time.UTC,
	)))

	// Output:
	// 2000-01-01 21:00:00 +0000 UTC
}

func TestDispatcher_CronFunc(t *testing.T) {
	d := NewDispatcher(10)

	// cron expr
	cronExpr, err := NewCronExpr("* * * * * *")
	if err != nil {
		t.Log(err)
		return
	}

	// cron
	var c *Cron
	c = d.CronFunc(cronExpr, func() {
		t.Log("My name is Leaf")
		c.Stop()
	})

	// dispatch
	(<-d.ChanTimer).Cb()

	// Output:
	// My name is Leaf
}
