package main

import (
	"fmt"
	"time"
)

func main() {

	followerNodeCount := 2

	cluster := NewCluster(followerNodeCount)

	fmt.Println(cluster.Get("abc"))

	cluster.Inc("abc", 10)

	time.Sleep(10 * time.Millisecond)

	fmt.Println(cluster.Get("abc"))
}
