package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Cluster struct {
	timer         *Timer
	leaderNode    *Node
	followerNodes []*Node

	clusterMutex *sync.Mutex
}

func NewCluster(followerNodeCount int) *Cluster {

	if followerNodeCount < 1 {
		followerNodeCount = 1
	}

	cluster := &Cluster{
		timer:        NewTimer(),
		clusterMutex: &sync.Mutex{},
	}

	cluster.leaderNode = NewNode("leader", true, cluster.timer, nil, -1)

	for i := 1; i <= followerNodeCount; i++ {
		cluster.followerNodes = append(cluster.followerNodes, NewNode(fmt.Sprint("follower-", i), false, cluster.timer, cluster.leaderNode, 5*i))
	}

	return cluster
}

func (cluster *Cluster) Get(key string) int {

	// get from a bunch of follower nodes and decide based on timestamp i guess

	idx := rand.Intn(len(cluster.followerNodes))

	readNode := cluster.followerNodes[idx]

	return cluster.getFromNode(readNode, key)

}

func (cluster *Cluster) getFromNode(readNode *Node, key string) int {
	data := readNode.Get(key)

	if data != nil {
		fmt.Println(time.Now(), "GET", key, "responded by", readNode.name, "val:", data.value, "createdOn:", data.createdOn)
		return data.value
	}

	fmt.Println(time.Now(), "GET", key, "responded by", readNode.name, "didn't contain key.")

	return 0
}

func (cluster *Cluster) Inc(key string, val int) {
	if err := cluster.leaderNode.Inc(key, val); err != nil {
		log.Fatalf(err.Error())
	}
}
