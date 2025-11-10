//go:build all || dab

package hiero

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

var (
	nodeIDToUpdate        = uint64(0)
	originalNodeAccountId = AccountID{Account: 3}
)

func TestIntegrationNodeUpdateTransactionCanExecute(t *testing.T) {
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
		SetNodeID(nodeIDToUpdate).
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
		SetNodeID(nodeIDToUpdate).
		DeleteGrpcWebProxyEndpoint().
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountIdToTheSameAccount(t *testing.T) {
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
		SetNodeID(nodeIDToUpdate).
		SetAccountID(originalNodeAccountId).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

func TestIntegrationNodeUpdateTransactionChangeNodeAccountIdMissingAdminSig(t *testing.T) {
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
		SetNodeID(nodeIDToUpdate).
		SetAccountID(operator).
		Execute(client)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationNodeUpdateTransactionChangeNodeAccountIdMissingAccountSig(t *testing.T) {
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
	nodeAccountId := *receipt.AccountID

	resp, err = NewNodeUpdateTransaction().
		SetNodeID(nodeIDToUpdate).
		SetAccountID(nodeAccountId).
		Execute(client)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationNodeUpdateTransactionChangeNodeAccountIdToNonExistentAccountId(t *testing.T) {
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
		SetNodeID(nodeIDToUpdate).
		SetAccountID(AccountID{Account: 9999999}).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountIdToDeletedAccountId(t *testing.T) {
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
		SetNodeID(nodeIDToUpdate).
		SetAccountID(newAccount).
		FreezeWith(client)

	resp, err = frozen.Sign(newAccountKey).Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.ErrorContains(t, err, "exceptional receipt status: ACCOUNT_DELETED")
}

func TestIntegrationNodeUpdateTransactionChangeNodeAccountINoBalance(t *testing.T) {
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
		SetNodeID(nodeIDToUpdate).
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
	network["localhost:50211"] = originalNodeAccountId
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
	newAccountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
	newNodeAccountID := *receipt.AccountID

	// update node account id
	frozen, err := NewNodeUpdateTransaction().
		SetNodeID(nodeIDToUpdate).
		SetAccountID(newNodeAccountID).
		FreezeWith(client)

	require.NoError(t, err)
	resp, err = frozen.Sign(newAccountKey).Execute(client)
	require.NoError(t, err)

	// wait for mirror node to import data
	time.Sleep(time.Second * 5)

	newAccountKey, err = PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// submit to the updated node, retries
	resp, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		SetNodeAccountIDs([]AccountID{originalNodeAccountId, {Account: 4}}).
		Execute(client)
	require.NoError(t, err)
	receipt, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// verify address book has been updated
	key1 := newNodeAccountID
	key2 := AccountID{Account: 4}

	node1, ok := client.network._GetNodeForAccountID(key1)
	require.True(t, ok)
	require.Equal(t, newNodeAccountID.String(), node1.accountID.String())
	node2, ok := client.network._GetNodeForAccountID(key2)
	require.True(t, ok)
	require.Equal(t, AccountID{Account: 4}.String(), node2.accountID.String())

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
		SetNodeID(nodeIDToUpdate).
		SetNodeAccountIDs([]AccountID{newNodeAccountID}).
		SetAccountID(originalNodeAccountId).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountWithoutMirrorNodeSetup(t *testing.T) {
	// Set the network
	network := make(map[string]AccountID)
	network["localhost:50211"] = originalNodeAccountId
	network["localhost:51211"] = AccountID{Account: 4}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	defer client.Close()

	// Set the operator to be account 0.0.2
	originalOperatorKey, err := PrivateKeyFromStringEd25519("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 2}, originalOperatorKey)

	// create the account that will be the node account id
	newAccountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
	newNodeAccountID := *receipt.AccountID

	// update node account id
	frozen, err := NewNodeUpdateTransaction().
		SetNodeID(nodeIDToUpdate).
		SetAccountID(newNodeAccountID).
		FreezeWith(client)

	require.NoError(t, err)
	resp, err = frozen.Sign(newAccountKey).Execute(client)
	require.NoError(t, err)

	newAccountKey, err = PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// submit to the updated node, retries
	resp, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		SetNodeAccountIDs([]AccountID{originalNodeAccountId, {Account: 4}}).
		Execute(client)
	require.NoError(t, err)
	receipt, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	key1 := originalNodeAccountId
	key2 := AccountID{Account: 4}

	// verify address book has NOT been updated
	node1, ok := client.network._GetNodeForAccountID(key1)
	require.True(t, ok)
	require.Equal(t, AccountID{Account: 8}.String(), node1.accountID.String())
	node2, ok := client.network._GetNodeForAccountID(key2)
	require.True(t, ok)
	require.Equal(t, AccountID{Account: 4}.String(), node2.accountID.String())

	// this transactin should succeed again because we will retry
	resp, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		Execute(client)
	require.NoError(t, err)
	receipt, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	// revert the node account id
	resp, err = NewNodeUpdateTransaction().
		SetNodeID(nodeIDToUpdate).
		SetAccountID(originalNodeAccountId).
		Execute(client)

	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}
