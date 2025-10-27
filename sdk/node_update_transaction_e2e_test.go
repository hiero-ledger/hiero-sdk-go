//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestIntegrationNodeUpdateTransactionCanExecute(t *testing.T) {
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

	resp, err := NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetDeclineReward(true).
		SetGrpcWebProxyEndpoint(Endpoint{
			domainName: "testWebUpdated.com",
			port:       123456,
		}).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

func TestIntegrationNodeUpdateTransactionDeleteGrpcWebProxyEndpoint(t *testing.T) {
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

	resp, err := NewNodeUpdateTransaction().
		SetNodeID(0).
		DeleteGrpcWebProxyEndpoint().
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountId(t *testing.T) {
	t.Skip()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	defer client.Close()
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	resp, err := NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetAccountID(AccountID{Account: 10}).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountIdToTheSameAccount(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	defer client.Close()
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	resp, err := NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetAccountID(AccountID{Account: 3}).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountIdInvalidSignature(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	defer client.Close()
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	newOperatorKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	newBalance := NewHbar(2)
	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newOperatorKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
	operator := *receipt.AccountID

	client.SetOperator(operator, newOperatorKey)

	resp, err = NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetAccountID(AccountID{Account: 3}).
		Execute(client)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountIdToNonExistentAccountId(t *testing.T) {
	t.Parallel()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	defer client.Close()
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	resp, err := NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetAccountID(AccountID{Account: 9999999}).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	// TODO this should be INVALID_ACCOUNT_ID
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_NODE_ACCOUNT_ID")
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountIdToDeletedAccountId(t *testing.T) {
	t.Parallel()
	t.Skip()

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	defer client.Close()
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	newAccountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		Execute(client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
	newAccount := *receipt.AccountID

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(newAccount).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	require.NoError(t, err)

	resp, err = tx.Sign(newAccountKey).Execute(client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	resp, err = NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetAccountID(newAccount).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "exceptional receipt status: ACCOUNT_DELETED")
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountIdRetry(t *testing.T) {
	t.Skip()
	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	defer client.Close()
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	resp, err := NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetAccountID(AccountID{Account: 10}).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	newAccountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	resp, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		Execute(client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}
