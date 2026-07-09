//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newNetworkMockNodes() map[string]AccountID {
	nodes := make(map[string]AccountID, 2)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["2.testnet.hedera.com:50211"] = AccountID{0, 0, 5, nil, nil, nil}
	nodes["3.testnet.hedera.com:50211"] = AccountID{0, 0, 6, nil, nil, nil}
	nodes["4.testnet.hedera.com:50211"] = AccountID{0, 0, 7, nil, nil, nil}
	return nodes
}

func TestUnitNetworkAddressBookGetsSet(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	network._SetTransportSecurity(true)

	ledgerID, err := LedgerIDFromString("mainnet")
	require.NoError(t, err)

	network._SetLedgerID(*ledgerID)
	require.NoError(t, err)

	require.True(t, network.addressBook != nil)
}

func TestUnitNetworkIncreaseBackoffConcurrent(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	nodes := newNetworkMockNodes()
	err := network.SetNetwork(nodes)
	require.NoError(t, err)

	node := network._GetNode()
	require.NotNil(t, node)

	numThreads := 20
	var wg sync.WaitGroup
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			network._IncreaseBackoff(node)
			wg.Done()
		}()
	}
	wg.Wait()

	require.Equal(t, len(nodes)-1, len(network.healthyNodes))
}

func TestUnitConcurrentGetNodeReadmit(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	nodes := newNetworkMockNodes()
	err := network.SetNetwork(nodes)
	network._SetMinNodeReadmitPeriod(0)
	network._SetMaxNodeReadmitPeriod(0)
	require.NoError(t, err)

	for _, node := range network.nodes {
		node._SetMaxBackoff(-1 * time.Minute)
	}

	numThreads := 3
	var wg sync.WaitGroup
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			for i := 0; i < 20; i++ {
				node := network._GetNode()
				network._IncreaseBackoff(node)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	network._ReadmitNodes()
	require.Equal(t, len(nodes), len(network.healthyNodes))
}

func TestUnitConcurrentNodeAccess(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	nodes := newNetworkMockNodes()
	err := network.SetNetwork(nodes)
	network._SetMinNodeReadmitPeriod(0)
	network._SetMaxNodeReadmitPeriod(0)
	require.NoError(t, err)

	for _, node := range network.nodes {
		node._SetMaxBackoff(-1 * time.Minute)
	}

	numThreads := 3
	var wg sync.WaitGroup
	node := network._GetNode()
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			for i := 0; i < 20; i++ {
				network._GetNode()
				network._IncreaseBackoff(node)
				node._IsHealthy()
				node._GetAttempts()
				node._GetReadmitTime()
				node._Wait()
				node._InUse()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	network._ReadmitNodes()
	require.Equal(t, len(nodes), len(network.healthyNodes))
}

func TestUnitConcurrentNodeGetChannel(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	nodes := newNetworkMockNodes()
	err := network.SetNetwork(nodes)
	require.NoError(t, err)

	numThreads := 20
	var wg sync.WaitGroup
	node := network._GetNode()
	wg.Add(numThreads)
	logger := NewLogger("", LoggerLevelError)
	for i := 0; i < numThreads; i++ {
		go func() {
			node._GetChannel(logger)
			wg.Done()
		}()
	}
	wg.Wait()
	network._ReadmitNodes()
	require.Equal(t, len(nodes), len(network.healthyNodes))
}

// Re-applying an unchanged network must keep the existing nodes and their open
// connections rather than closing and recreating them.

func TestUnitNetworkRefreshRemovesNoNodesWhenUnchanged(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	err := network.SetNetwork(newNetworkMockNodes())
	require.NoError(t, err)

	// Incoming address book identical to the installed nodes, keyed by address
	// the same way _Network.SetNetwork builds it.
	incoming := make(map[string]_IManagedNode, len(network.nodes))
	for _, node := range network.nodes {
		incoming[node._GetAddress()] = node
	}

	require.Empty(t, _GetNodesToRemove(incoming, network.nodes),
		"an unchanged address book must not mark any node for removal")
}

func TestUnitNetworkRefreshPreservesNodeObjects(t *testing.T) {
	t.Parallel()

	nodes := newNetworkMockNodes()
	network := _NewNetwork()
	err := network.SetNetwork(nodes)
	require.NoError(t, err)

	before := nodesByAccount(&network)

	// Re-apply the identical address book, as the scheduled network update does.
	err = network.SetNetwork(nodes)
	require.NoError(t, err)

	after := nodesByAccount(&network)

	require.Equal(t, len(before), len(after))
	for account, node := range before {
		require.Samef(t, node, after[account],
			"node 0.0.%d was recreated by an unchanged refresh", account)
	}
}

func TestUnitNetworkRefreshPreservesNodeConnections(t *testing.T) {
	t.Parallel()

	nodes := newNetworkMockNodes()
	network := _NewNetwork()
	err := network.SetNetwork(nodes)
	require.NoError(t, err)

	// Open a node's channel, standing in for an in-flight request holding a live
	// connection. grpc.NewClient is lazy, so nothing is dialed.
	account := AccountID{0, 0, 3, nil, nil, nil}
	node, ok := network._GetNodeForAccountID(account)
	require.True(t, ok)
	_, err = node._GetChannel(NewLogger("", LoggerLevelError))
	require.NoError(t, err)
	require.NotNil(t, node.channel)

	err = network.SetNetwork(nodes)
	require.NoError(t, err)

	refreshed, ok := network._GetNodeForAccountID(account)
	require.True(t, ok)
	require.Same(t, node, refreshed,
		"unchanged refresh replaced the node and dropped its open connection")
	require.NotNil(t, refreshed.channel,
		"unchanged refresh closed the node's open connection")
}

func nodesByAccount(network *_Network) map[uint64]*_Node {
	byAccount := make(map[uint64]*_Node, len(network.nodes))
	for _, node := range network.nodes {
		if n, ok := node.(*_Node); ok {
			byAccount[n.accountID.Account] = n
		}
	}
	return byAccount
}
