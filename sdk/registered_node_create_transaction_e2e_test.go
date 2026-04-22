//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrationRegisteredNodeCreateTransactionCanExecute(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Generate admin key
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Build a block node service endpoint with an IPv4 address
	endpoint := &BlockNodeServiceEndpoint{}
	endpoint.SetIPAddress(net.IPv4(192, 168, 1, 1).To4()).
		SetPort(50211)
	endpoint.AddEndpointApi(BlockNodeApiPublish)

	// Execute the RegisteredNodeCreateTransaction
	tx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetDescription("e2e test registered node").
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	resp, err := tx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// Verify the receipt contains a non-zero registeredNodeId
	require.Greater(t, receipt.RegisteredNodeId, uint64(0), "registeredNodeId should be non-zero")

	time.Sleep(time.Second * 5)

	// Query the registered node from the mirror node and verify fields
	book, err := NewRegisteredNodeAddressBookQuery().
		SetRegisteredNodeId(receipt.RegisteredNodeId).
		Execute(client)
	require.NoError(t, err)
	require.Len(t, book.RegisteredNodes, 1)
	require.Equal(t, "e2e test registered node", book.RegisteredNodes[0].Description)
	require.Equal(t, receipt.RegisteredNodeId, book.RegisteredNodes[0].RegisteredNodeID)
	require.Equal(t, adminKey.PublicKey().String(), book.RegisteredNodes[0].AdminKey.String())
	require.Len(t, book.RegisteredNodes[0].ServiceEndpoints, 1)

	ep, ok := book.RegisteredNodes[0].ServiceEndpoints[0].(*BlockNodeServiceEndpoint)
	require.True(t, ok, "expected *BlockNodeServiceEndpoint")
	require.Equal(t, net.IPv4(192, 168, 1, 1).To4(), net.IP(ep.GetIPAddress()))
	require.Equal(t, uint32(50211), ep.GetPort())
	require.Equal(t, []BlockNodeApi{BlockNodeApiPublish}, ep.GetEndpointApis())
}

func TestIntegrationRegisteredNodeCreateTransactionMirrorNodeEndpointSucceeds(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Generate admin key
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Build a mirror node service endpoint
	endpoint := &MirrorNodeServiceEndpoint{}
	endpoint.SetIPAddress(net.IPv4(10, 0, 0, 1).To4()).SetPort(443)

	// Execute the RegisteredNodeCreateTransaction
	tx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	resp, err := tx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// Verify the receipt contains a non-zero registeredNodeId
	require.Greater(t, receipt.RegisteredNodeId, uint64(0), "registeredNodeId should be non-zero")
}

func TestIntegrationRegisteredNodeCreateTransactionRpcRelayEndpointSucceeds(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Generate admin key
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Build an RPC relay service endpoint
	endpoint := &RpcRelayServiceEndpoint{}
	endpoint.SetDomainName("rpc.example.com").SetPort(8545)

	// Execute the RegisteredNodeCreateTransaction
	tx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	resp, err := tx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// Verify the receipt contains a non-zero registeredNodeId
	require.Greater(t, receipt.RegisteredNodeId, uint64(0), "registeredNodeId should be non-zero")
}

func TestIntegrationRegisteredNodeCreateTransactionGeneralServiceEndpointSucceeds(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Generate admin key
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Build a general service endpoint with a description
	endpoint := &GeneralServiceEndpoint{}
	endpoint.SetDomainName("general.example.com").
		SetPort(9000).
		SetDescription("custom-service")

	// Execute the RegisteredNodeCreateTransaction
	tx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	resp, err := tx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// Verify the receipt contains a non-zero registeredNodeId
	require.Greater(t, receipt.RegisteredNodeId, uint64(0), "registeredNodeId should be non-zero")
}

func TestIntegrationRegisteredNodeCreateTransactionMixedEndpointsSucceeds(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Generate admin key
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Build a block node service endpoint
	blockEndpoint := &BlockNodeServiceEndpoint{}
	blockEndpoint.SetIPAddress(net.IPv4(192, 168, 1, 1).To4()).SetPort(50211)
	blockEndpoint.AddEndpointApi(BlockNodeApiPublish)

	// Build a mirror node service endpoint
	mirrorEndpoint := &MirrorNodeServiceEndpoint{}
	mirrorEndpoint.SetIPAddress(net.IPv4(10, 0, 0, 1).To4()).SetPort(443)

	// Build an RPC relay service endpoint
	rpcEndpoint := &RpcRelayServiceEndpoint{}
	rpcEndpoint.SetDomainName("rpc.example.com").SetPort(8545)

	// Build a general service endpoint
	generalEndpoint := &GeneralServiceEndpoint{}
	generalEndpoint.SetDomainName("general.example.com").
		SetPort(9000).
		SetDescription("custom-service")

	// Execute the RegisteredNodeCreateTransaction with all four endpoint types
	tx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetServiceEndpoints([]RegisteredServiceEndpoint{blockEndpoint, mirrorEndpoint, rpcEndpoint, generalEndpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	resp, err := tx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// Verify the receipt contains a non-zero registeredNodeId
	require.Greater(t, receipt.RegisteredNodeId, uint64(0), "registeredNodeId should be non-zero")
}

func TestIntegrationRegisteredNodeCreateTransactionWithDescriptionSucceeds(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Generate admin key
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Build a block node service endpoint
	endpoint := &BlockNodeServiceEndpoint{}
	endpoint.SetIPAddress(net.IPv4(192, 168, 1, 1).To4()).SetPort(50211)
	endpoint.AddEndpointApi(BlockNodeApiPublish)

	description := "e2e test registered node with description"

	// Execute the RegisteredNodeCreateTransaction with a description
	tx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetDescription(description).
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	resp, err := tx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// Verify the receipt contains a non-zero registeredNodeId
	require.Greater(t, receipt.RegisteredNodeId, uint64(0), "registeredNodeId should be non-zero")
}

func TestIntegrationRegisteredNodeCreateTransactionFailsIfNoAdminKeySet(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Build a block node service endpoint
	endpoint := &BlockNodeServiceEndpoint{}
	endpoint.SetIPAddress(net.IPv4(192, 168, 1, 1).To4()).SetPort(50211)
	endpoint.AddEndpointApi(BlockNodeApiPublish)

	// Execute the RegisteredNodeCreateTransaction without an admin key
	tx, err := NewRegisteredNodeCreateTransaction().
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = tx.Execute(client)
	require.ErrorContains(t, err, "KEY_REQUIRED")
}

func TestIntegrationRegisteredNodeCreateTransactionFailsIfEmptyEndpoints(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Generate admin key
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Execute the RegisteredNodeCreateTransaction with empty service endpoints
	tx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetServiceEndpoints([]RegisteredServiceEndpoint{}).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = tx.Sign(adminKey).Execute(client)
	require.ErrorContains(t, err, "INVALID_REGISTERED_ENDPOINT")
}
