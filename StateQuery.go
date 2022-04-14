package main

type StateQuery struct {
	replicaNode        *Node
	replicaIncListSize int
	recordFetchCount   int
}

func NewStateQuery(replicaNode *Node, replicaIncListSize, recordFetchCount int) *StateQuery {
	return &StateQuery{
		replicaNode:        replicaNode,
		replicaIncListSize: replicaIncListSize,
		recordFetchCount:   recordFetchCount,
	}
}
