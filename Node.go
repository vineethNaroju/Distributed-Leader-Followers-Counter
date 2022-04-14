package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	isLeader bool
	timer    *Timer
	leader   *Node

	incList      []*Incop
	incListMutex *sync.Mutex

	logChannel           chan *Incop
	stateQueryChannel    chan StateQuery
	stateResponseChannel chan *StateQueryResponse
	stateIncListSize     int

	store      map[string]*Data
	storeMutex *sync.Mutex
}

func NewNode(isLeader bool, timer *Timer, leader *Node) *Node {

	node := &Node{
		isLeader:     isLeader,
		timer:        timer,
		leader:       leader,
		incListMutex: &sync.Mutex{},
		logChannel:   make(chan *Incop),
		store:        make(map[string]*Data),
		storeMutex:   &sync.Mutex{},
	}

	if isLeader {
		node.stateIncListSize = -1
		node.stateQueryChannel = make(chan StateQuery)
		node.leaderDaemon()
	} else {
		node.stateIncListSize = 0
		node.stateResponseChannel = make(chan *StateQueryResponse)
		node.followerDaemon()
	}

	return node
}

func (node *Node) Get(key string) *Data {
	node.storeMutex.Lock()
	defer node.storeMutex.Unlock()

	if data, ok := node.store[key]; ok {
		return data
	}

	return nil
}

func (node *Node) Inc(key string, val int) error {
	if node.isLeader {
		createdOn := node.timer.Get()
		data := NewData(createdOn, val)

		node.leaderWriteToLog(NewIncop(key, data))

		node.storeMutex.Lock()
		defer node.storeMutex.Unlock()

		if d, ok := node.store[key]; ok {
			d.createdOn = createdOn
			d.value += val
		} else {
			node.store[key] = data
		}

		return nil
	}

	return fmt.Errorf("Node is not leader, so cant accept writes")
}

func (node *Node) leaderWriteToLog(op *Incop) {
	node.logChannel <- op
}

func (node *Node) leaderLogIngressDaemon() {
	go func() {
		for {
			op := <-node.logChannel

			node.incListMutex.Lock()
			node.incList = append(node.incList, op)
			node.incListMutex.Unlock()
		}
	}()
}

func (node *Node) leaderStateQueryIngressDaemon() {
	go func() {
		for {
			query := <-node.stateQueryChannel
			response := NewStateQueryResponse(false, []*Incop{})

			node.incListMutex.Lock()

			if len(node.incList) <= query.replicaIncListSize {
				response.end = true
			} else {

				i := query.replicaIncListSize

				for query.recordFetchCount > 0 && i < len(node.incList) {
					response.incList = append(response.incList, node.incList[i])
					i++
				}

				if query.replicaIncListSize+query.recordFetchCount < len(node.incList) {
					response.end = false
				} else {
					response.end = true
				}
			}

			node.incListMutex.Unlock()

			query.replicaNode.stateResponseChannel <- response
		}
	}()
}

func (node *Node) leaderDaemon() {
	node.leaderLogIngressDaemon()
	node.leaderStateQueryIngressDaemon()
}

func (node *Node) followerStateSyncDaemon() {
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)

		}
	}()
}

func (node *Node) followerDaemon() {
	node.followerStateSyncDaemon()
}
