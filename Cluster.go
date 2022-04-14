package main

import "sync"

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

	cluster.leaderNode = NewNode(true, cluster.timer, nil)

	for i := 0; i < followerNodeCount; i++ {
		cluster.followerNodes = append(cluster.followerNodes, NewNode(false, cluster.timer, cluster.leaderNode))
	}

	return cluster
}

func (cluster *Cluster) Get(key string) int {

	// get from a bunch of follower nodes and decide based on timestamp i guess

	return 0
}

func (cluster *Cluster) Inc(key string, val int) {

}
