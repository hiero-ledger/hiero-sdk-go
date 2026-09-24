package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"sync/atomic"

	"github.com/pkg/errors"
)

type _MirrorNetwork struct {
	_ManagedNetwork
	// nextNode is the round-robin cursor over healthyNodes, shared by every concurrent caller.
	nextNode atomic.Uint64
}

func _NewMirrorNetwork() *_MirrorNetwork {
	return &_MirrorNetwork{
		_ManagedNetwork: _NewManagedNetwork(),
	}
}

func (network *_MirrorNetwork) _SetNetwork(newNetwork []string) (err error) {
	network.healthyNodesMutex.Lock()
	defer network.healthyNodesMutex.Unlock()

	newMirrorNetwork := make(map[string]_IManagedNode)
	for _, url := range newNetwork {
		if newMirrorNetwork[url], err = _NewMirrorNode(url); err != nil {
			return err
		}
	}

	return network._ManagedNetwork._SetNetwork(newMirrorNetwork)
}

func (network *_MirrorNetwork) _GetNetwork() []string {
	network.healthyNodesMutex.RLock()
	defer network.healthyNodesMutex.RUnlock()

	temp := make([]string, 0)
	for url := range network._ManagedNetwork.network { //nolint
		temp = append(temp, url)
	}

	return temp
}

// nolint:unused
// Deprecated
// _SetTransportSecurity is no longer supported, as only secured connections are now allowed.
func (network *_MirrorNetwork) _SetTransportSecurity(transportSecurity bool) *_MirrorNetwork {
	return network
}

// _GetNextMirrorNode selects round-robin rather than at random, so a network with one node down
// is not a fresh coin flip on every call.
func (network *_MirrorNetwork) _GetNextMirrorNode() (*_MirrorNode, error) {
	network.healthyNodesMutex.RLock()
	defer network.healthyNodesMutex.RUnlock()

	if len(network.healthyNodes) == 0 {
		return nil, errors.New("no healthy nodes")
	}

	index := (network.nextNode.Add(1) - 1) % uint64(len(network.healthyNodes))
	node := network.healthyNodes[index]
	if node, ok := node.(*_MirrorNode); ok {
		return node, nil
	}
	return &_MirrorNode{}, nil
}
