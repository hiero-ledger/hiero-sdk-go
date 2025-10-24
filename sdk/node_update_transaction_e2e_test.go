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

func TestIntegrationNodeUpdateTransactionCanChangeNodeAccountId(t *testing.T) {
	// TODO
	t.Skip()
	t.Parallel()

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
	// client.SetLogger(NewLogger("Go-SDK", LoggerLevelTrace))

	time.Sleep(2 * time.Second)

	client._UpdateAddressBook()

	// resp, err = NewNodeUpdateTransaction().
	// 	SetNodeID(0).
	// 	SetDescription("testUpdated").
	// 	SetServiceEndpoints([]hiero.Endpoint{serviceEndpoint}).
	// 	SetAccountID(AccountID{Account: 12}).
	// 	Execute(client)

	// require.NoError(t, err)
	// _, err = resp.SetValidateStatus(true).GetReceipt(client)
	// require.NoError(t, err)

	// newKey, err := PrivateKeyGenerateEd25519()
	// require.NoError(t, err)

	// newBalance := NewHbar(2)

	// resp, err = NewAccountCreateTransaction().
	// 	SetKeyWithoutAlias(newKey).
	// 	SetInitialBalance(newBalance).
	// 	SetMaxAutomaticTokenAssociations(100).
	// 	Execute(client)

	// require.NoError(t, err)

	// receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	// require.NoError(t, err)
	// fmt.Println(receipt.AccountID)

}
