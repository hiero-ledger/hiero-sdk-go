//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

// setupLocalClient creates a client pointed at the local network for e2e tests.
func setupRegisteredNodeUpdateLocalClient(t *testing.T) *Client {
	t.Helper()

	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	return client
}

// createRegisteredNode creates a registered node with the given admin key and returns the registeredNodeId.
func createRegisteredNode(t *testing.T, client *Client, adminKey PrivateKey) uint64 {
	t.Helper()

	endpoint := &BlockNodeServiceEndpoint{}
	endpoint.SetIPAddress(net.IPv4(10, 0, 0, 1).To4()).SetPort(8080)
	endpoint.AddEndpointApi(BlockNodeApiStatus)

	createTx, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetDescription("test node").
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	createResp, err := createTx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	createReceipt, err := createResp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	require.Greater(t, createReceipt.RegisteredNodeId, uint64(0), "registeredNodeId should be non-zero")
	t.Log(createReceipt.RegisteredNodeId)
	fmt.Println(createReceipt)
	return createReceipt.RegisteredNodeId
}

func TestIntegrationRegisteredNodeUpdateTransactionUpdateDescription(t *testing.T) {
	t.Parallel()

	client := setupRegisteredNodeUpdateLocalClient(t)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	registeredNodeId := createRegisteredNode(t, client, adminKey)
	require.Greater(t, registeredNodeId, uint64(0))

	updateTx, err := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(registeredNodeId).
		SetDescription("updated description").
		FreezeWith(client)
	require.NoError(t, err)

	updateResp, err := updateTx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	_, err = updateResp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// TODO: query the registered node and verify description == "updated description"
}

func TestIntegrationRegisteredNodeUpdateTransactionReplaceEndpoints(t *testing.T) {
	t.Parallel()

	client := setupRegisteredNodeUpdateLocalClient(t)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	registeredNodeId := createRegisteredNode(t, client, adminKey)
	require.Greater(t, registeredNodeId, uint64(0))

	newEndpoint := &BlockNodeServiceEndpoint{}
	newEndpoint.SetIPAddress(net.IPv4(172, 16, 0, 1).To4()).SetPort(9090)
	newEndpoint.AddEndpointApi(BlockNodeApiPublish)

	updateTx, err := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(registeredNodeId).
		SetServiceEndpoints([]RegisteredServiceEndpoint{newEndpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	updateResp, err := updateTx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	_, err = updateResp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// TODO: query the registered node and verify endpoints were replaced
}

func TestIntegrationRegisteredNodeUpdateTransactionUpdateAdminKeyBothSign(t *testing.T) {
	t.Parallel()

	client := setupRegisteredNodeUpdateLocalClient(t)

	oldAdminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	registeredNodeId := createRegisteredNode(t, client, oldAdminKey)
	require.Greater(t, registeredNodeId, uint64(0))

	newAdminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	updateTx, err := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(registeredNodeId).
		SetAdminKey(newAdminKey.PublicKey()).
		FreezeWith(client)
	require.NoError(t, err)

	updateResp, err := updateTx.
		Sign(oldAdminKey).
		Sign(newAdminKey).
		Execute(client)
	require.NoError(t, err)

	_, err = updateResp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// TODO: query the registered node and verify admin key was updated to newAdminKey
}

func TestIntegrationRegisteredNodeUpdateTransactionUpdateAdminKeyOnlyOldSigns(t *testing.T) {
	t.Parallel()

	client := setupRegisteredNodeUpdateLocalClient(t)

	oldAdminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	registeredNodeId := createRegisteredNode(t, client, oldAdminKey)
	require.Greater(t, registeredNodeId, uint64(0))

	newAdminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	updateTx, err := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(registeredNodeId).
		SetAdminKey(newAdminKey.PublicKey()).
		FreezeWith(client)
	require.NoError(t, err)

	updateResp, err := updateTx.
		Sign(oldAdminKey).
		Execute(client)
	require.NoError(t, err)

	_, err = updateResp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationRegisteredNodeUpdateTransactionFailsIfNonExistentId(t *testing.T) {
	t.Parallel()

	client := setupRegisteredNodeUpdateLocalClient(t)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	updateTx, err := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(9999999).
		SetAdminKey(adminKey).
		SetDescription("should fail").
		FreezeWith(client)
	require.NoError(t, err)

	updateResp, err := updateTx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	_, err = updateResp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "INVALID_REGISTERED_NODE_ID")
}

func TestIntegrationRegisteredNodeUpdateTransactionReplaceDomainEndpoint(t *testing.T) {
	t.Parallel()

	client := setupRegisteredNodeUpdateLocalClient(t)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	registeredNodeId := createRegisteredNode(t, client, adminKey)
	require.Greater(t, registeredNodeId, uint64(0))

	domainEndpoint := &BlockNodeServiceEndpoint{}
	domainEndpoint.SetDomainName("node.example.com").SetPort(443)
	domainEndpoint.AddEndpointApi(BlockNodeApiStatus)

	updateTx, err := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(registeredNodeId).
		SetServiceEndpoints([]RegisteredServiceEndpoint{domainEndpoint}).
		FreezeWith(client)
	require.NoError(t, err)

	updateResp, err := updateTx.Sign(adminKey).Execute(client)
	require.NoError(t, err)

	_, err = updateResp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// TODO: query the registered node and verify endpoint uses domain name "node.example.com"
}
