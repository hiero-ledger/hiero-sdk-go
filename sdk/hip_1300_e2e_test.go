//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	maximumTransactionSize = 130000
	systemAccountID        = "0.0.2"
	systemAccountKey       = "302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137"
)

// createSystemAccountClient creates a client configured with the system account (0.0.2)
// This is needed for transactions that exceed the normal 6KB limit
func createSystemAccountClient(t *testing.T) *Client {
	env := NewIntegrationTestEnv(t)
	client := env.Client

	sysAccID, err := AccountIDFromString(systemAccountID)
	require.NoError(t, err)

	sysAccKey, err := PrivateKeyFromString(systemAccountKey)
	require.NoError(t, err)

	client.SetOperator(sysAccID, sysAccKey)

	return client
}

// TestIntegrationHIP1300TransactionWithMoreThan6KBDataWithSignatures tests that a system account
// can create a transaction with more than 6KB of data with signatures
func TestIntegrationHIP1300TransactionWithMoreThan6KBDataWithSignatures(t *testing.T) {
	t.Parallel()

	client := createSystemAccountClient(t)
	defer client.Close()

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create a transaction and freeze it
	transaction, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		FreezeWith(client)
	require.NoError(t, err)

	// Keep signing until we exceed the maximum transaction size
	for {
		txBytes, err := transaction.ToBytes()
		require.NoError(t, err)

		if len(txBytes) >= maximumTransactionSize {
			break
		}

		signingKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)

		transaction = transaction.Sign(signingKey)
	}

	// Execute the large transaction - should succeed with system account
	resp, err := transaction.Execute(client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

// TestIntegrationHIP1300TransactionWithMoreThan6KBDataInFile tests that a system account
// can create a file with more than 6KB of data
func TestIntegrationHIP1300TransactionWithMoreThan6KBDataInFile(t *testing.T) {
	t.Parallel()

	client := createSystemAccountClient(t)
	defer client.Close()

	// Create a file with 10KB of data
	contents := make([]byte, 1024*10)
	for i := range contents {
		contents[i] = 1
	}

	transaction, err := NewFileCreateTransaction().
		SetContents(contents).
		FreezeWith(client)
	require.NoError(t, err)

	// Execute the large file creation - should succeed with system account
	resp, err := transaction.Execute(client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

// TestIntegrationHIP1300TransactionWithMoreThan6KBDataInFileNormalAccount tests that a normal account
// cannot create a file with more than 6KB of data
func TestIntegrationHIP1300TransactionWithMoreThan6KBDataInFileNormalAccount(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a file with 10KB of data
	contents := make([]byte, 1024*10)
	for i := range contents {
		contents[i] = 1
	}

	transaction, err := NewFileCreateTransaction().
		SetContents(contents).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	// Execute the large file creation - should fail with TRANSACTION_OVERSIZE
	resp, err := transaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.Error(t, err)
	require.ErrorContains(t, err, "TRANSACTION_OVERSIZE")
}

// TestIntegrationHIP1300TransactionWithMoreThan6KBDataWithSignaturesNormalAccount tests that a normal account
// cannot create a transaction with more than 6KB of data with signatures
func TestIntegrationHIP1300TransactionWithMoreThan6KBDataWithSignaturesNormalAccount(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	regularUserKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create a transaction and freeze it
	transaction, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(regularUserKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	// Keep signing until we exceed the maximum transaction size
	for {
		txBytes, err := transaction.ToBytes()
		require.NoError(t, err)

		if len(txBytes) >= maximumTransactionSize {
			break
		}

		signingKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)

		transaction = transaction.Sign(signingKey)
	}

	// Execute the large transaction - should fail with TRANSACTION_OVERSIZE
	resp, err := transaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.Error(t, err)
	require.ErrorContains(t, err, "TRANSACTION_OVERSIZE")
}

