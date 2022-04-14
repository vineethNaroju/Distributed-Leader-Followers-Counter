package main

import (
	"sync"
	"time"
)

type Timer struct {
	count      int
	countMutex *sync.Mutex
}

func NewTimer() *Timer {
	timer := &Timer{
		count:      0,
		countMutex: &sync.Mutex{},
	}

	timer.daemon()

	return timer
}

func (timer *Timer) daemon() {
	go func() {
		for {
			time.Sleep(1 * time.Millisecond)

			timer.countMutex.Lock()
			timer.count++
			timer.countMutex.Unlock()
		}
	}()
}

func (timer *Timer) Get() int {
	timer.countMutex.Lock()
	defer timer.countMutex.Unlock()

	return timer.count
}
