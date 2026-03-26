//go:build all || e2e

package hiero

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationRegisteredNodeDeleteTransactionCanExecute(t *testing.T) {
	t.Parallel()

	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Create a registered node
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	endpoint := &BlockNodeServiceEndpoint{}
	endpoint.SetIPAddress(net.IPv4(10, 0, 0, 1).To4()).SetPort(8080)
	endpoint.SetEndpointApi(BlockNodeApiStatus)

	createTx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetDescription("test node for delete").
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	createResp, err := createTx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	createReceipt, err := createResp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	registeredNodeId := createReceipt.RegisteredNodeId

	// Delete the registered node
	deleteTx, err := NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(registeredNodeId).
		FreezeWith(client)
	require.NoError(t, err)

	deleteResp, err := deleteTx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	_, err = deleteResp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

func TestIntegrationRegisteredNodeDeleteTransactionFailsIfAlreadyDeleted(t *testing.T) {
	t.Parallel()

	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Create a registered node
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	endpoint := &BlockNodeServiceEndpoint{}
	endpoint.SetIPAddress(net.IPv4(10, 0, 0, 1).To4()).SetPort(8080)
	endpoint.SetEndpointApi(BlockNodeApiStatus)

	createTx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetDescription("test node for double delete").
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	createResp, err := createTx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	createReceipt, err := createResp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	registeredNodeId := createReceipt.RegisteredNodeId

	// First delete should succeed
	deleteTx, err := NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(registeredNodeId).
		FreezeWith(client)
	require.NoError(t, err)

	deleteResp, err := deleteTx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	_, err = deleteResp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// Second delete should fail
	deleteTx2, err := NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(registeredNodeId).
		FreezeWith(client)
	require.NoError(t, err)

	deleteResp2, err := deleteTx2.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	_, err = deleteResp2.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "INVALID_REGISTERED_NODE_ID")
}

func TestIntegrationRegisteredNodeDeleteTransactionFailsIfNonExistentNode(t *testing.T) {
	t.Parallel()

	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Try to delete a non-existent registered node
	deleteTx, err := NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(9999999).
		FreezeWith(client)
	require.NoError(t, err)

	deleteResp, err := deleteTx.Execute(client)
	require.NoError(t, err)

	_, err = deleteResp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "INVALID_REGISTERED_NODE_ID")
}
