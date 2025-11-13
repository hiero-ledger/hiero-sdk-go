package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"crypto/rand"
	"math"
	"math/big"
	"sync"
	"time"
)

type _ManagedNetwork struct {
	network                map[string][]_IManagedNode
	nodes                  []_IManagedNode
	healthyNodes           []_IManagedNode
	healthyNodesMutex      *sync.RWMutex
	maxNodeAttempts        int
	minBackoff             time.Duration
	maxBackoff             time.Duration
	maxNodesPerTransaction *int
	ledgerID               *LedgerID
	transportSecurity      bool
	verifyCertificate      bool
	minNodeReadmitPeriod   time.Duration
	maxNodeReadmitPeriod   time.Duration
	earliestReadmitTime    time.Time
}

func _NewManagedNetwork() _ManagedNetwork {
	return _ManagedNetwork{
		network:                map[string][]_IManagedNode{},
		nodes:                  []_IManagedNode{},
		healthyNodes:           []_IManagedNode{},
		healthyNodesMutex:      &sync.RWMutex{},
		maxNodeAttempts:        -1,
		minBackoff:             8 * time.Second,
		maxBackoff:             1 * time.Hour,
		maxNodesPerTransaction: nil,
		ledgerID:               nil,
		transportSecurity:      false,
		verifyCertificate:      false,
		minNodeReadmitPeriod:   8 * time.Second,
		maxNodeReadmitPeriod:   1 * time.Hour,
	}
}

func (mn *_ManagedNetwork) _SetNetwork(network map[string]_IManagedNode) error {
	newNodes := make([]_IManagedNode, len(mn.nodes))
	newNodeKeys := map[string]bool{}
	newNodeValues := map[string]bool{}

	// Copy all the nodes into the `newNodes` list
	copy(newNodes, mn.nodes)

	// Remove nodes from the old mn which do not belong to the new mn
	for _, index := range _GetNodesToRemove(network, newNodes) {
		node := newNodes[index]

		if err := node._Close(); err != nil {
			return err
		}

		if index == len(newNodes)-1 {
			newNodes = newNodes[:index]
		} else {
			newNodes = append(newNodes[:index], newNodes[index+1:]...)
		}
	}

	for _, node := range newNodes {
		newNodeKeys[node._GetKey()] = true
		newNodeValues[node._GetAddress()] = true
	}

	for key, value := range network {
		_, keyOk := newNodeKeys[key]
		_, valueOk := newNodeValues[value._GetAddress()]

		if keyOk && valueOk {
			continue
		}

		newNodes = append(newNodes, value)
	}

	newNetwork, newHealthyNodes := _CreateNetworkFromNodes(newNodes)

	mn.nodes = newNodes
	mn.network = newNetwork
	mn.healthyNodes = newHealthyNodes

	return nil
}

func (mn *_ManagedNetwork) _ReadmitNodes() {
	now := time.Now()

	mn.healthyNodesMutex.Lock()
	defer mn.healthyNodesMutex.Unlock()

	if mn.earliestReadmitTime.Before(now) {
		nextEarliestReadmitTime := now.Add(mn.maxNodeReadmitPeriod)

		for _, node := range mn.nodes {
			if node._GetReadmitTime() != nil && node._GetReadmitTime().After(now) && node._GetReadmitTime().Before(nextEarliestReadmitTime) {
				nextEarliestReadmitTime = *node._GetReadmitTime()
			}
		}

		mn.earliestReadmitTime = nextEarliestReadmitTime
		if mn.earliestReadmitTime.Before(now.Add(mn.minNodeReadmitPeriod)) {
			mn.earliestReadmitTime = now.Add(mn.minNodeReadmitPeriod)
		}

	outer:
		for _, node := range mn.nodes {
			for _, healthyNode := range mn.healthyNodes {
				if node == healthyNode {
					continue outer
				}
			}

			if node._GetReadmitTime().Before(now) {
				mn.healthyNodes = append(mn.healthyNodes, node)
			}
		}
	}
}

func (mn *_ManagedNetwork) _GetNumberOfNodesForTransaction() int { // nolint
	mn._ReadmitNodes()
	if mn.maxNodesPerTransaction != nil {
		return int(math.Min(float64(*mn.maxNodesPerTransaction), float64(len(mn.healthyNodes))))
	}

	return len(mn.healthyNodes)
}

func (mn *_ManagedNetwork) _SetMaxNodesPerTransaction(max int) {
	mn.maxNodesPerTransaction = &max
}

func (mn *_ManagedNetwork) _SetMaxNodeAttempts(max int) {
	mn.maxNodeAttempts = max
}

func (mn *_ManagedNetwork) _GetMaxNodeAttempts() int {
	return mn.maxNodeAttempts
}

func (mn *_ManagedNetwork) _SetMinNodeReadmitPeriod(min time.Duration) {
	mn.minNodeReadmitPeriod = min
	mn.earliestReadmitTime = time.Now().Add(mn.minNodeReadmitPeriod)
}

func (mn *_ManagedNetwork) _GetMinNodeReadmitPeriod() time.Duration {
	return mn.minNodeReadmitPeriod
}

func (mn *_ManagedNetwork) _SetMaxNodeReadmitPeriod(max time.Duration) {
	mn.maxNodeReadmitPeriod = max
}

func (mn *_ManagedNetwork) _GetMaxNodeReadmitPeriod() time.Duration {
	return mn.maxNodeReadmitPeriod
}

func (mn *_ManagedNetwork) _SetMinBackoff(minBackoff time.Duration) {
	mn.minBackoff = minBackoff
	for _, nod := range mn.healthyNodes {
		if nod != nil {
			nod._SetMinBackoff(minBackoff)
		}
	}
}

func (mn *_ManagedNetwork) _GetNode() _IManagedNode {
	mn._ReadmitNodes()
	mn.healthyNodesMutex.RLock()
	defer mn.healthyNodesMutex.RUnlock()

	if len(mn.healthyNodes) == 0 {
		panic("failed to find a healthy working node")
	}

	bg := big.NewInt(int64(len(mn.healthyNodes)))
	index, _ := rand.Int(rand.Reader, bg)
	return mn.healthyNodes[index.Int64()]
}

func (mn *_ManagedNetwork) _GetMinBackoff() time.Duration {
	return mn.minBackoff
}

func (mn *_ManagedNetwork) _SetMaxBackoff(maxBackoff time.Duration) {
	mn.maxBackoff = maxBackoff
	for _, node := range mn.healthyNodes {
		node._SetMaxBackoff(maxBackoff)
	}
}

func (mn *_ManagedNetwork) _GetMaxBackoff() time.Duration {
	return mn.maxBackoff
}

func (mn *_ManagedNetwork) _GetLedgerID() *LedgerID {
	return mn.ledgerID
}

func (mn *_ManagedNetwork) _SetLedgerID(id LedgerID) *_ManagedNetwork {
	mn.ledgerID = &id
	return mn
}

func (mn *_ManagedNetwork) _Close() error {
	for _, conn := range mn.healthyNodes {
		if err := conn._Close(); err != nil {
			return err
		}
	}

	return nil
}

func _CreateNetworkFromNodes(nodes []_IManagedNode) (network map[string][]_IManagedNode, healthyNodes []_IManagedNode) {
	healthyNodes = []_IManagedNode{}
	network = map[string][]_IManagedNode{}

	for _, node := range nodes {
		if node._IsHealthy() {
			healthyNodes = append(healthyNodes, node)
		}

		value, ok := network[node._GetKey()]
		if !ok {
			value = []_IManagedNode{}
		}
		value = append(value, node)
		network[node._GetKey()] = value
	}

	return network, healthyNodes
}

func (mn *_ManagedNetwork) _SetTransportSecurity(transportSecurity bool) (err error) {
	if mn.transportSecurity != transportSecurity {
		if err := mn._Close(); err != nil {
			return err
		}

		newNodes := make([]_IManagedNode, len(mn.nodes))

		copy(newNodes, mn.nodes)

		for i, node := range newNodes {
			if transportSecurity {
				newNodes[i] = node._ToSecure()
			} else {
				newNodes[i] = node._ToInsecure()
			}
		}

		newNetwork, newHealthyNodes := _CreateNetworkFromNodes(newNodes)

		mn.nodes = newNodes
		mn.healthyNodes = newHealthyNodes
		mn.network = newNetwork
	}

	mn.transportSecurity = transportSecurity
	return nil
}

func _GetNodesToRemove(network map[string]_IManagedNode, nodes []_IManagedNode) []int {
	nodeIndices := []int{}

	for i := len(nodes) - 1; i >= 0; i-- {
		if _, ok := network[nodes[i]._GetKey()]; !ok {
			nodeIndices = append(nodeIndices, i)
		}
	}

	return nodeIndices
}

func (mn *_ManagedNetwork) _SetVerifyCertificate(verify bool) *_ManagedNetwork {
	for _, node := range mn.nodes {
		node._SetVerifyCertificate(verify)
	}

	mn.verifyCertificate = verify
	return mn
}

func (mn *_ManagedNetwork) _GetVerifyCertificate() bool {
	return mn.verifyCertificate
}
