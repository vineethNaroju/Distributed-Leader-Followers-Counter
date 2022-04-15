package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	name                    string
	syncPeriodInMillisecond time.Duration

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

func NewNode(name string, isLeader bool, timer *Timer, leader *Node, syncPeriodInMillisecond int) *Node {

	node := &Node{
		name:         name,
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
		node.syncPeriodInMillisecond = time.Duration(syncPeriodInMillisecond) * time.Millisecond
		node.stateResponseChannel = make(chan *StateQueryResponse)
		node.followerDaemon()
	}

	fmt.Println(time.Now(), "created node", node.name, "with sync time in", syncPeriodInMillisecond, "ms")

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
			fmt.Println(time.Now(), "leader incList appended by key:", op.key, "value:", op.data.value, "createdOn:", op.data.createdOn)
			node.incListMutex.Unlock()
		}
	}()
}

func (node *Node) leaderStateQueryIngressDaemon() {
	go func() {
		for {
			query := <-node.stateQueryChannel

			// fmt.Println(time.Now(), "leader got state query by", query.replicaNode.name, " with replicaIncListSize", query.replicaIncListSize, "and recordFetchCount", query.recordFetchCount)

			response := NewStateQueryResponse([]*Incop{})

			node.incListMutex.Lock()

			if len(node.incList) > query.replicaIncListSize {

				i := query.replicaIncListSize

				for query.recordFetchCount > 0 && i < len(node.incList) {
					response.incList = append(response.incList, node.incList[i])
					i++
					query.recordFetchCount--
				}

				// if query.replicaIncListSize+query.recordFetchCount < len(node.incList) {
				// 	response.end = false
				// } else {
				// 	response.end = true
				// }
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
			time.Sleep(node.syncPeriodInMillisecond)

			node.leader.stateQueryChannel <- *NewStateQuery(node, node.stateIncListSize, 2)

			res := <-node.stateResponseChannel

			if len(res.incList) == 0 {
				continue
			}

			fmt.Println(time.Now(), node.name, "got state query response from leader as", res)

			node.storeMutex.Lock()

			for _, op := range res.incList {
				if _, ok := node.store[op.key]; ok {
					node.store[op.key].createdOn = op.data.createdOn
					node.store[op.key].value += op.data.value
				} else {
					node.store[op.key] = op.data
				}

				fmt.Println(time.Now(), node.name, "updated store for key:", op.key, "as value:", node.store[op.key].value, "createdOn:", op.data.createdOn)
			}

			node.stateIncListSize += len(res.incList)

			node.storeMutex.Unlock()
		}
	}()
}

func (node *Node) followerDaemon() {
	node.followerStateSyncDaemon()
}
