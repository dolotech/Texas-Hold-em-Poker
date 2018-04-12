package utils

import (
	"container/heap"
)

type TimeOuter interface {
	TimeOut(int64)
}

type Timer struct {
	TimeOuter
	id       uint32
	end      int64 //结束时间
	interval int64 //迭代周期
	index    int
}

type TimerQueue []*Timer

func (this TimerQueue) Len() int {
	return len(this)
}

func (this TimerQueue) Less(i, j int) bool {
	return this[i].end < this[j].end
}

func (this TimerQueue) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
	this[i].index = i
	this[j].index = j
}

func (this *TimerQueue) Push(x interface{}) {
	tmp := *this
	n := len(tmp)
	tmp = tmp[0 : n+1]
	timer := x.(*Timer)
	timer.index = n
	tmp[n] = timer
	*this = tmp
}

func (this *TimerQueue) Pop() interface{} {
	tmp := *this
	n := len(tmp)
	timer := tmp[n-1]
	tmp[n-1] = nil
	timer.index = -1
	*this = tmp[0 : n-1]
	return timer
}

type TimerManager struct {
	id uint32
	tq TimerQueue
}

func NewTimerManager(size int) *TimerManager {
	if size == 0 {
		size = 1024
	}
	return &TimerManager{tq: make([]*Timer, 0, size)}
}

func (this *TimerManager) AddTimer(i TimeOuter, e int64, iv int64) uint32 {
	if cap(this.tq) <= len(this.tq) {
		return 0
	}
	timer := &Timer{TimeOuter: i, interval: iv, end: e}
	this.id++
	timer.id = this.id
	heap.Push(&this.tq, timer)
	return timer.id
}

func (this *TimerManager) RemoveTimer(id uint32) {
	for _, timer := range this.tq {
		if timer.id == id {
			heap.Remove(&this.tq, timer.index)
			return
		}
	}
}

var queue *Queue = &Queue{}

func (this *TimerManager) Run(now int64, limit int) {
	for len(this.tq) > 0 {
		tmp := this.tq[0]
		if tmp.end <= now {
			timer := heap.Pop(&this.tq).(*Timer)
			queue.Push(timer.TimeOuter)
			if timer.interval > 0 {
				timer.end += timer.interval
				heap.Push(&this.tq, timer)
			}
		} else {
			break
		}
		if limit > 0 && queue.Len() >= limit {
			break
		}
	}

	for queue.Len() > 0 {
		queue.Pop().(TimeOuter).TimeOut(now)
	}
}

func (this *TimerManager) dump() {
	queue := &Queue{}
	for len(this.tq) > 0 {
		timer := heap.Pop(&this.tq).(*Timer)
		queue.Push(timer)
		println("Timer:", timer.id, timer.index, timer.end, timer.interval)
	}
	for queue.Len() > 0 {
		heap.Push(&this.tq, queue.Pop().(*Timer))
	}
}
