//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"
	"time"

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

// func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountId(t *testing.T) {
// 	// Set the network
// 	network := make(map[string]AccountID)
// 	network["localhost:50211"] = AccountID{Account: 1038}
// 	client, err := ClientForNetworkV2(network)
// 	require.NoError(t, err)
// 	defer client.Close()
// 	mirror := []string{"localhost:5600"}
// 	client.SetMirrorNetwork(mirror)

// 	// Set the operator to be account 0.0.2
// 	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
// 	require.NoError(t, err)
// 	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

// 	resp, err := NewNodeUpdateTransaction().
// 		SetNodeID(0).
// 		SetDescription("testUpdated").
// 		SetAccountID(AccountID{Account: 3}).
// 		Execute(client)

// 	require.NoError(t, err)
// 	_, err = resp.SetValidateStatus(true).GetReceipt(client)
// 	require.NoError(t, err)
// }

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
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_NODE_ACCOUNT_ID")
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountIdToDeletedAccountId(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: unskip when services implements check for this")

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

	frozen, err := NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetAccountID(newAccount).
		FreezeWith(client)

	resp, err = frozen.Sign(newAccountKey).Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "exceptional receipt status: ACCOUNT_DELETED")
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountINoBalance(t *testing.T) {
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

	newAccountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		Execute(client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
	newAccount := *receipt.AccountID

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	frozen, err := NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetAccountID(newAccount).
		FreezeWith(client)

	resp, err = frozen.Sign(newAccountKey).Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "exceptional receipt status: NODE_ACCOUNT_HAS_ZERO_BALANCE")
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountUpdateAddressbookAndRetry(t *testing.T) {

	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = AccountID{Account: 3}
	network["localhost:51211"] = AccountID{Account: 4}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	defer client.Close()
	mirror := []string{"localhost:5600"}
	client.SetMirrorNetwork(mirror)

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// create the account that will be the node account id
	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(originalOperatorKey.PublicKey()).
		SetInitialBalance(HbarFrom(1, HbarUnit("hbar"))).
		Execute(client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
	newNodeAccountID := *receipt.AccountID

	// update node account id
	// 0.0.3 -> 0.0.1003
	resp, err = NewNodeUpdateTransaction().
		SetNodeID(0).
		SetDescription("testUpdated").
		SetAccountID(newNodeAccountID).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// wait for mirror node to import data
	time.Sleep(time.Second * 10)

	newAccountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	// submit to node 3 and node 4, node 3 fails, node 4 succeeds
	resp, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		Execute(client)
	require.NoError(t, err)

	// verify address book has been updated
	key1 := newNodeAccountID
	key2 := AccountID{Account: 4}
	require.Equal(t, newNodeAccountID.String(), client.network.addressBook[key1].AccountID.String())
	require.Equal(t, AccountID{Account: 4}.String(), client.network.addressBook[key2].AccountID.String())

	// this transactin should succeed
	resp, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		SetNodeAccountIDs([]AccountID{newNodeAccountID}).
		Execute(client)
	require.NoError(t, err)
	receipt, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// revert the node account id
	resp, err = NewNodeUpdateTransaction().
		SetNodeID(0).
		SetNodeAccountIDs([]AccountID{newNodeAccountID}).
		SetDescription("testUpdated").
		SetAccountID(AccountID{Account: 3}).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}
