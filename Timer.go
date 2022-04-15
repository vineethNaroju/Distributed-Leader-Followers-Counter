package main

import (
	"sync"
)

type Timer struct {
	count      int
	countMutex *sync.Mutex
}

func NewTimer() *Timer {
	timer := &Timer{
		count:      -1,
		countMutex: &sync.Mutex{},
	}

	return timer
}

func (timer *Timer) Get() int {
	timer.countMutex.Lock()
	defer timer.countMutex.Unlock()
	timer.count++

	return timer.count
}
