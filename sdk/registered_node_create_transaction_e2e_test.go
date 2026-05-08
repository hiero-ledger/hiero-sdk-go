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

	// Verify the receipt contains a registered node ID
	require.NotNil(t, receipt.RegisteredNodeId, "registeredNodeId should be set on the receipt")

	time.Sleep(time.Second * 5)

	// Query the registered node from the mirror node and verify fields
	book, err := NewRegisteredNodeAddressBookQuery().
		SetRegisteredNodeId(*receipt.RegisteredNodeId).
		Execute(client)
	require.NoError(t, err)
	require.Len(t, book.RegisteredNodes, 1)
	require.Equal(t, "e2e test registered node", book.RegisteredNodes[0].Description)
	require.Equal(t, *receipt.RegisteredNodeId, book.RegisteredNodes[0].RegisteredNodeID)
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

	// Verify the receipt contains a registered node ID
	require.NotNil(t, receipt.RegisteredNodeId, "registeredNodeId should be set on the receipt")

	time.Sleep(time.Second * 5)

	// Query the registered node from the mirror node and verify the endpoint
	book, err := NewRegisteredNodeAddressBookQuery().
		SetRegisteredNodeId(*receipt.RegisteredNodeId).
		Execute(client)
	require.NoError(t, err)
	require.Len(t, book.RegisteredNodes, 1)
	require.Equal(t, *receipt.RegisteredNodeId, book.RegisteredNodes[0].RegisteredNodeID)
	require.Equal(t, adminKey.PublicKey().String(), book.RegisteredNodes[0].AdminKey.String())
	require.Len(t, book.RegisteredNodes[0].ServiceEndpoints, 1)

	ep, ok := book.RegisteredNodes[0].ServiceEndpoints[0].(*MirrorNodeServiceEndpoint)
	require.True(t, ok, "expected *MirrorNodeServiceEndpoint")
	require.Equal(t, net.IPv4(10, 0, 0, 1).To4(), net.IP(ep.GetIPAddress()))
	require.Equal(t, uint32(443), ep.GetPort())
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

	// Verify the receipt contains a registered node ID
	require.NotNil(t, receipt.RegisteredNodeId, "registeredNodeId should be set on the receipt")

	time.Sleep(time.Second * 5)

	// Query the registered node from the mirror node and verify the endpoint
	book, err := NewRegisteredNodeAddressBookQuery().
		SetRegisteredNodeId(*receipt.RegisteredNodeId).
		Execute(client)
	require.NoError(t, err)
	require.Len(t, book.RegisteredNodes, 1)
	require.Equal(t, *receipt.RegisteredNodeId, book.RegisteredNodes[0].RegisteredNodeID)
	require.Equal(t, adminKey.PublicKey().String(), book.RegisteredNodes[0].AdminKey.String())
	require.Len(t, book.RegisteredNodes[0].ServiceEndpoints, 1)

	ep, ok := book.RegisteredNodes[0].ServiceEndpoints[0].(*RpcRelayServiceEndpoint)
	require.True(t, ok, "expected *RpcRelayServiceEndpoint")
	require.Equal(t, "rpc.example.com", ep.GetDomainName())
	require.Equal(t, uint32(8545), ep.GetPort())
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

	// Verify the receipt contains a registered node ID
	require.NotNil(t, receipt.RegisteredNodeId, "registeredNodeId should be set on the receipt")

	time.Sleep(time.Second * 5)

	// Query the registered node from the mirror node and verify the endpoint
	book, err := NewRegisteredNodeAddressBookQuery().
		SetRegisteredNodeId(*receipt.RegisteredNodeId).
		Execute(client)
	require.NoError(t, err)
	require.Len(t, book.RegisteredNodes, 1)
	require.Equal(t, *receipt.RegisteredNodeId, book.RegisteredNodes[0].RegisteredNodeID)
	require.Equal(t, adminKey.PublicKey().String(), book.RegisteredNodes[0].AdminKey.String())
	require.Len(t, book.RegisteredNodes[0].ServiceEndpoints, 1)

	ep, ok := book.RegisteredNodes[0].ServiceEndpoints[0].(*GeneralServiceEndpoint)
	require.True(t, ok, "expected *GeneralServiceEndpoint")
	require.Equal(t, "general.example.com", ep.GetDomainName())
	require.Equal(t, uint32(9000), ep.GetPort())
	require.Equal(t, "custom-service", ep.GetDescription())
}

func TestIntegrationRegisteredNodeCreateTransactionMixedEndpointsSucceeds(t *testing.T) {
	t.Parallel()

	// Set the network
	client, err := ClientForNetworkV2(map[string]AccountID{"localhost:50211": {Account: 3}})
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{"localhost:5600"})

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// Generate admin key
	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Build one endpoint per kind: block node, mirror node, RPC relay, general service.
	blockEndpoint := (&BlockNodeServiceEndpoint{}).SetIPAddress(net.IPv4(192, 168, 1, 1).To4()).SetPort(50211).AddEndpointApi(BlockNodeApiPublish)
	mirrorEndpoint := (&MirrorNodeServiceEndpoint{}).SetIPAddress(net.IPv4(10, 0, 0, 1).To4()).SetPort(443)
	rpcEndpoint := (&RpcRelayServiceEndpoint{}).SetDomainName("rpc.example.com").SetPort(8545)
	generalEndpoint := (&GeneralServiceEndpoint{}).SetDomainName("general.example.com").SetPort(9000).SetDescription("custom-service")

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

	// Verify the receipt contains a registered node ID
	require.NotNil(t, receipt.RegisteredNodeId, "registeredNodeId should be set on the receipt")

	time.Sleep(time.Second * 5)

	// Query the registered node from the mirror node and verify all four endpoints
	// come back in the order they were sent.
	book, err := NewRegisteredNodeAddressBookQuery().SetRegisteredNodeId(*receipt.RegisteredNodeId).Execute(client)
	require.NoError(t, err)
	require.Len(t, book.RegisteredNodes, 1)
	require.Equal(t, *receipt.RegisteredNodeId, book.RegisteredNodes[0].RegisteredNodeID)
	require.Equal(t, adminKey.PublicKey().String(), book.RegisteredNodes[0].AdminKey.String())
	require.Len(t, book.RegisteredNodes[0].ServiceEndpoints, 4)

	bn, ok := book.RegisteredNodes[0].ServiceEndpoints[0].(*BlockNodeServiceEndpoint)
	require.True(t, ok, "expected *BlockNodeServiceEndpoint at index 0")
	require.Equal(t, net.IPv4(192, 168, 1, 1).To4(), net.IP(bn.GetIPAddress()))
	require.Equal(t, uint32(50211), bn.GetPort())
	require.Equal(t, []BlockNodeApi{BlockNodeApiPublish}, bn.GetEndpointApis())

	mn, ok := book.RegisteredNodes[0].ServiceEndpoints[1].(*MirrorNodeServiceEndpoint)
	require.True(t, ok, "expected *MirrorNodeServiceEndpoint at index 1")
	require.Equal(t, net.IPv4(10, 0, 0, 1).To4(), net.IP(mn.GetIPAddress()))
	require.Equal(t, uint32(443), mn.GetPort())

	rr, ok := book.RegisteredNodes[0].ServiceEndpoints[2].(*RpcRelayServiceEndpoint)
	require.True(t, ok, "expected *RpcRelayServiceEndpoint at index 2")
	require.Equal(t, "rpc.example.com", rr.GetDomainName())
	require.Equal(t, uint32(8545), rr.GetPort())

	gs, ok := book.RegisteredNodes[0].ServiceEndpoints[3].(*GeneralServiceEndpoint)
	require.True(t, ok, "expected *GeneralServiceEndpoint at index 3")
	require.Equal(t, "general.example.com", gs.GetDomainName())
	require.Equal(t, uint32(9000), gs.GetPort())
	require.Equal(t, "custom-service", gs.GetDescription())
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

	// Verify the receipt contains a registered node ID
	require.NotNil(t, receipt.RegisteredNodeId, "registeredNodeId should be set on the receipt")

	time.Sleep(time.Second * 5)

	// Query the registered node from the mirror node and verify the description
	book, err := NewRegisteredNodeAddressBookQuery().
		SetRegisteredNodeId(*receipt.RegisteredNodeId).
		Execute(client)
	require.NoError(t, err)
	require.Len(t, book.RegisteredNodes, 1)
	require.Equal(t, *receipt.RegisteredNodeId, book.RegisteredNodes[0].RegisteredNodeID)
	require.Equal(t, description, book.RegisteredNodes[0].Description)
	require.Equal(t, adminKey.PublicKey().String(), book.RegisteredNodes[0].AdminKey.String())
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
