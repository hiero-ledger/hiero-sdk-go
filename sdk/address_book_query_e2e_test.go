//go:build all || DAB

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Updates local previewnet.pb file on run
func TestIntegrationAddressBookQueryPreviewnet(t *testing.T) {
	client, err := ClientFromConfig([]byte(`{"network":"previewnet"}`))
	require.NoError(t, err)
	client.SetMirrorNetwork(previewnetMirror)

	previewnet, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		SetMaxAttempts(5).
		Execute(client)
	require.NoError(t, err)
	require.Greater(t, len(previewnet.NodeAddresses), 0)

	filePreviewnet, err := os.OpenFile("addressbook/previewnet.pb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	_, err = filePreviewnet.Write(previewnet.ToBytes())
	require.NoError(t, err)

	err = filePreviewnet.Close()
	require.NoError(t, err)
}

// Updates local testnet.pb file on run
func TestIntegrationAddressBookQueryTestnet(t *testing.T) {
	client, err := ClientFromConfig([]byte(`{"network":"testnet"}`))
	require.NoError(t, err)
	client.SetMirrorNetwork(testnetMirror)

	testnet, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		SetMaxAttempts(5).
		Execute(client)

	require.NoError(t, err)
	require.Greater(t, len(testnet.NodeAddresses), 0)

	fileTestnet, err := os.OpenFile("addressbook/testnet.pb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	_, err = fileTestnet.Write(testnet.ToBytes())
	require.NoError(t, err)

	err = fileTestnet.Close()
	require.NoError(t, err)
}

// Updates local mainnet.pb file on run
func TestIntegrationAddressBookQueryMainnet(t *testing.T) {
	client, err := ClientFromConfig([]byte(`{"network":"mainnet"}`))
	require.NoError(t, err)
	client.SetMirrorNetwork(mainnetMirror)

	mainnet, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		SetMaxAttempts(5).
		Execute(client)
	require.NoError(t, err)
	require.Greater(t, len(mainnet.NodeAddresses), 0)

	fileMainnet, err := os.OpenFile("addressbook/mainnet.pb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	_, err = fileMainnet.Write(mainnet.ToBytes())
	require.NoError(t, err)

	err = fileMainnet.Close()
	require.NoError(t, err)
}

func TestIntegrationAddressBookQueryLocal(t *testing.T) {
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

	client._UpdateAddressBook()

	newAccountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	// submit to node 3
	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newAccountKey.PublicKey()).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		Execute(client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

}
