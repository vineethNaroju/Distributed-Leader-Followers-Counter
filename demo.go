package main

import (
	"math/rand"
	"sync"
	"time"
)

func Demo() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	followerNodeCount := 3

	cluster := NewCluster(followerNodeCount)

	go func() {
		for {
			time.Sleep(time.Millisecond * 15)
			cluster.Inc("abc", rand.Intn(100))
		}
	}()

	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		cluster.Get("abc")
	}

	wg.Wait()
}
