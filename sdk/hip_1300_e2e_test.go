//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func setupSystemAccountOperator(t *testing.T, client *Client) {
	accountID, err := AccountIDFromString("0.0.2")
	require.NoError(t, err)
	privateKey, err := PrivateKeyFromString("302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137")
	require.NoError(t, err)
	client.SetOperator(accountID, privateKey)
}

func setupRegularAccountOperator(t *testing.T, env IntegrationTestEnv) {
	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key).
		SetInitialBalance(NewHbar(10)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(*receipt.AccountID, key)
}

// TestIntegrationHIP1300TransactionWithMoreThan6KBDataWithSignatures tests that a system account
// can create a transaction with more than 6KB of data with signatures
func TestIntegrationHIP1300TransactionWithMoreThan6KBDataWithSignatures(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	setupSystemAccountOperator(t, env.Client)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transaction, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	transaction, err = transaction.SignWithOperator(env.Client)
	require.NoError(t, err)

	// Add signatures to exceed 6KB
	for i := 0; i < 70; i++ {
		signingKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)
		transaction = transaction.Sign(signingKey)
	}

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.Greater(t, len(txBytes), 6144, "Transaction should exceed 6KB")

	resp, err := transaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

// TestIntegrationHIP1300TransactionWithMoreThan6KBDataInFile tests that a system account
// can create a file with more than 6KB of data
func TestIntegrationHIP1300TransactionWithMoreThan6KBDataInFile(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	setupSystemAccountOperator(t, env.Client)

	contents := make([]byte, 1024*10)
	for i := range contents {
		contents[i] = 1
	}

	transaction, err := NewFileCreateTransaction().
		SetContents(contents).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := transaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

// TestIntegrationHIP1300TransactionWithMoreThan6KBDataInFileNormalAccount tests that a normal account
// cannot create a file with more than 6KB of data
func TestIntegrationHIP1300TransactionWithMoreThan6KBDataInFileNormalAccount(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	setupRegularAccountOperator(t, env)

	contents := make([]byte, 1024*10)
	for i := range contents {
		contents[i] = 1
	}

	transaction, err := NewFileCreateTransaction().
		SetContents(contents).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = transaction.Execute(env.Client)
	require.Error(t, err)
	require.ErrorContains(t, err, "TRANSACTION_OVERSIZE")
}

// TestIntegrationHIP1300TransactionWithMoreThan6KBDataWithSignaturesNormalAccount tests that a normal account
// cannot create a transaction with more than 6KB of data with signatures
func TestIntegrationHIP1300TransactionWithMoreThan6KBDataWithSignaturesNormalAccount(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	setupRegularAccountOperator(t, env)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transaction, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	transaction, err = transaction.SignWithOperator(env.Client)
	require.NoError(t, err)

	// Add signatures to exceed 6KB
	for i := 0; i < 70; i++ {
		signingKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)
		transaction = transaction.Sign(signingKey)
	}

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.Greater(t, len(txBytes), 6144, "Transaction should exceed 6KB")

	_, err = transaction.Execute(env.Client)
	require.Error(t, err)
	require.ErrorContains(t, err, "TRANSACTION_OVERSIZE")
}
