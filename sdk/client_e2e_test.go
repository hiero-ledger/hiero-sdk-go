//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationClientCanExecuteSerializedTransactionFromAnotherClient(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	client2, err := ClientForNetworkV2(env.Client.GetNetwork())
	require.NoError(t, err)
	client2.SetOperator(env.OperatorID, env.OperatorKey)

	tx, err := NewTransferTransaction().AddHbarTransfer(env.OperatorID, HbarFromTinybar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(1)).SetNodeAccountIDs([]AccountID{{Account: 3}}).FreezeWith(env.Client)
	require.NoError(t, err)
	txBytes, err := tx.ToBytes()
	FromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)
	txFromBytes, ok := FromBytes.(TransferTransaction)
	require.True(t, ok)
	resp, err := txFromBytes.Execute(client2)
	require.NoError(t, err)
	reciept, err := resp.SetValidateStatus(true).GetReceipt(client2)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, reciept.Status)
}

func TestIntegrationClientCanFailGracefullyWhenDoesNotHaveNodeOfAnotherClient(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Get one of the nodes of the network from the original client
	var address string
	for key := range env.Client.GetNetwork() {
		address = key
		break
	}
	// Use that node to create a network for the second client but with a different node account id
	var network = map[string]AccountID{
		address: {Account: 99},
	}

	client2, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	client2.SetOperator(env.OperatorID, env.OperatorKey)

	// Create a transaction with a node using original client
	tx, err := NewTransferTransaction().AddHbarTransfer(env.OperatorID, HbarFromTinybar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(1)).SetNodeAccountIDs([]AccountID{{Account: 3}}).FreezeWith(env.Client)
	require.NoError(t, err)
	txBytes, err := tx.ToBytes()
	FromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)
	txFromBytes, ok := FromBytes.(TransferTransaction)
	require.True(t, ok)

	// Try to execute it with the second client, which does not have the node
	_, err = txFromBytes.Execute(client2)
	require.Error(t, err)
	require.Equal(t, err.Error(), "Invalid node AccountID was set for transaction: 0.0.3")
}

func DisabledTestIntegrationClientPingAllBadNetwork(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	netwrk := _NewNetwork()
	netwrk.SetNetwork(env.Client.GetNetwork())

	tempClient := _NewClient(netwrk, env.Client.GetMirrorNetwork(), env.Client.GetLedgerID(), true, 0, 0)
	tempClient.SetOperator(env.OperatorID, env.OperatorKey)

	tempClient.SetMaxNodeAttempts(1)
	tempClient.SetMaxNodesPerTransaction(2)
	tempClient.SetMaxAttempts(3)
	net := tempClient.GetNetwork()
	assert.True(t, len(net) > 1)

	keys := make([]string, len(net))
	val := make([]AccountID, len(net))
	i := 0
	for st, n := range net {
		keys[i] = st
		val[i] = n
		i++
	}

	tempNet := make(map[string]AccountID, 2)
	tempNet["in.process.ew:3123"] = val[0]
	tempNet[keys[1]] = val[1]

	err := tempClient.SetNetwork(tempNet)
	require.NoError(t, err)

	tempClient.PingAll()

	net = tempClient.GetNetwork()
	i = 0
	for st, n := range net {
		keys[i] = st
		val[i] = n
		i++
	}

	_, err = NewAccountBalanceQuery().
		SetAccountID(val[0]).
		Execute(tempClient)
	require.NoError(t, err)

	assert.Equal(t, 1, len(tempClient.GetNetwork()))

}

func TestClientInitWithMirrorNetwork(t *testing.T) {
	t.Parallel()
	mirrorNetworkString := "testnet.mirrornode.hedera.com:443"
	client, err := ClientForMirrorNetwork([]string{mirrorNetworkString})
	require.NoError(t, err)

	mirrorNetwork := client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, mirrorNetworkString, mirrorNetwork[0])
	assert.NotEmpty(t, client.GetNetwork())

	client, err = ClientForMirrorNetworkWithShardAndRealm([]string{mirrorNetworkString}, 0, 0)
	require.NoError(t, err)

	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, mirrorNetworkString, mirrorNetwork[0])
	assert.NotEmpty(t, client.GetNetwork())
}

func TestClientIntegrationForMirrorNetworkWithShardAndRealm(t *testing.T) {
	t.Parallel()

	mirrorNetworkString := "testnet.mirrornode.hedera.com:443"
	client, err := ClientForMirrorNetworkWithShardAndRealm([]string{mirrorNetworkString}, 0, 0)
	require.NoError(t, err)
	require.NotNil(t, client)

	// TODO enable when we have non-zero realm and shard env
	// client, err = ClientForMirrorNetworkWithShardAndRealm([]string{mirrorNetworkString}, 5, 3)
	// require.NoError(t, err)
	// require.NotNil(t, client)
	// require.Equal(t, uint64(5), client.GetShard())
	// require.Equal(t, uint64(3), client.GetRealm())

	client, err = ClientForMirrorNetworkWithShardAndRealm([]string{}, 0, 0)
	require.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to query address book: no healthy nodes")
}

// Every node of a live network answers the probe, and PingAll probes every node.
func TestIntegrationClientPingAllNetworkNodes(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	for address, nodeID := range env.Client.GetNetwork() {
		require.NoError(t, env.Client.Ping(nodeID), "node %s at %s should be reachable", nodeID, address)
	}

	// PingAll swallows errors, so a rising use count is the proof it probed - valid only if healthy first.
	usesBefore := make(map[AccountID]int64)
	for address, nodeID := range env.Client.GetNetwork() {
		node, ok := env.Client.network._GetNodeForAccountID(nodeID)
		require.True(t, ok)
		require.True(t, node._IsHealthy(), "node %s at %s should be healthy before PingAll", nodeID, address)
		usesBefore[nodeID] = node._GetUseCount()
	}

	env.Client.PingAll()
	for address, nodeID := range env.Client.GetNetwork() {
		node, ok := env.Client.network._GetNodeForAccountID(nodeID)
		require.True(t, ok)
		assert.Greater(t, node._GetUseCount(), usesBefore[nodeID],
			"PingAll must probe node %s at %s", nodeID, address)
		// Only transport failures back a node off, so this catches unreachability, not every failure.
		assert.True(t, node._IsHealthy(), "node %s at %s should stay healthy after PingAll", nodeID, address)
	}
}

// The probe works against a live node with no operator configured: a COST_ANSWER query is free and
// unsigned, so it attaches no payment.
func TestIntegrationClientPingWithoutOperator(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	client, err := ClientForNetworkV2(env.Client.GetNetwork())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, client.Close())
	}()
	require.Nil(t, client.operator)

	for _, nodeID := range client.GetNetwork() {
		require.NoError(t, client.Ping(nodeID))
	}
}
