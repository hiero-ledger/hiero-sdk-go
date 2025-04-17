//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationBatchTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key).
		SetInitialBalance(NewHbar(1)).
		Batchify(env.Client, env.OperatorKey)
	require.NoError(t, err)

	batchTransaction := NewBatchTransaction().
		AddInnerTransaction(tx)

	resp, err := batchTransaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	receipt, err := batchTransaction.GetInnerTransactionIDs()[0].GetReceipt(env.Client)
	require.NoError(t, err)
	accountID := receipt.AccountID
	require.NotNil(t, accountID)

	info, err := NewAccountInfoQuery().
		SetAccountID(*accountID).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID.String(), info.AccountID.String())
}

func TestIntegrationBatchTransactionToFromBytes(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key).
		SetInitialBalance(NewHbar(1)).
		Batchify(env.Client, env.OperatorKey)
	require.NoError(t, err)

	batchTransaction := NewBatchTransaction().
		AddInnerTransaction(tx)

	batchTransactionBytes, err := batchTransaction.ToBytes()
	require.NoError(t, err)

	batchTransactionFromBytes, err := TransactionFromBytes(batchTransactionBytes)
	require.NoError(t, err)

	txInterface := batchTransactionFromBytes.(TransactionInterface)
	resp, err := TransactionExecute(txInterface, env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	receipt, err := batchTransaction.GetInnerTransactionIDs()[0].GetReceipt(env.Client)
	require.NoError(t, err)
	accountID := receipt.AccountID
	require.NotNil(t, accountID)

	info, err := NewAccountInfoQuery().
		SetAccountID(*accountID).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, accountID.String(), info.AccountID.String())
}

func TestIntegrationBatchTransactionLarge(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	batchTransaction := NewBatchTransaction()

	// Add 25 transactions (maximum limit)
	for i := 0; i < 25; i++ {
		key, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)

		tx, err := NewAccountCreateTransaction().
			SetKeyWithoutAlias(key).
			SetInitialBalance(NewHbar(1)).
			Batchify(env.Client, env.OperatorKey)
		require.NoError(t, err)

		batchTransaction.AddInnerTransaction(tx)
	}

	resp, err := batchTransaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify all inner transactions succeeded
	for _, txID := range batchTransaction.GetInnerTransactionIDs() {
		receipt, err := NewTransactionReceiptQuery().
			SetTransactionID(txID).
			Execute(env.Client)
		require.NoError(t, err)
		assert.Equal(t, StatusSuccess, receipt.Status)
	}
}

func TestIntegrationBatchTransactionEmpty(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	batchTransaction := NewBatchTransaction()
	_, err := batchTransaction.Execute(env.Client)
	assert.ErrorContains(t, err, "exceptional precheck status BATCH_LIST_EMTPY")
}

func TestIntegrationBatchTransactionBlacklisted(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Test FreezeTransaction
	freezeTx, err := NewFreezeTransaction().
		SetStartTime(time.Now()).
		SetFreezeType(FreezeTypeFreezeOnly).
		Batchify(env.Client, env.OperatorKey)
	require.NoError(t, err)

	batchTx := NewBatchTransaction()
	err = batchTx.AddInnerTransaction(freezeTx).freezeError
	assert.Equal(t, errTransactionTypeNotAllowed, err)

	// Test BatchTransaction
	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key).
		SetInitialBalance(NewHbar(1)).
		Batchify(env.Client, env.OperatorKey)
	require.NoError(t, err)

	innerBatchTx := NewBatchTransaction().
		AddInnerTransaction(tx)

	err = NewBatchTransaction().AddInnerTransaction(innerBatchTx).freezeError
	assert.Equal(t, errTransactionTypeNotAllowed, err)
}

func TestIntegrationBatchTransactionInvalidBatchKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	invalidKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key).
		SetInitialBalance(NewHbar(1)).
		Batchify(env.Client, invalidKey)
	require.NoError(t, err)

	batchTransaction := NewBatchTransaction().
		AddInnerTransaction(tx)

	resp, err := batchTransaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationBatchTransactionDifferentBatchKeys(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	batchKey1, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	batchKey2, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	batchKey3, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create and prepare first transfer
	key1, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	tx1, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key1).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)
	receipt1, err := tx1.GetReceipt(env.Client)
	require.NoError(t, err)
	account1 := receipt1.AccountID
	require.NotNil(t, account1)

	transfer1, err := NewTransferTransaction().
		AddHbarTransfer(env.OperatorID, NewHbar(1)).
		AddHbarTransfer(*account1, NewHbar(-1)).
		SetTransactionID(TransactionIDGenerate(*account1)).
		SetBatchKey(batchKey1).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer1.Sign(key1)

	// Create and prepare second transfer
	key2, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	tx2, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key2).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)
	receipt2, err := tx2.GetReceipt(env.Client)
	require.NoError(t, err)
	account2 := receipt2.AccountID
	require.NotNil(t, account2)

	transfer2, err := NewTransferTransaction().
		AddHbarTransfer(env.OperatorID, NewHbar(1)).
		AddHbarTransfer(*account2, NewHbar(-1)).
		SetTransactionID(TransactionIDGenerate(*account2)).
		SetBatchKey(batchKey2).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer2.Sign(key2)

	// Create and prepare third transfer
	key3, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	tx3, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key3).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)
	receipt3, err := tx3.GetReceipt(env.Client)
	require.NoError(t, err)
	account3 := receipt3.AccountID
	require.NotNil(t, account3)

	transfer3, err := NewTransferTransaction().
		AddHbarTransfer(env.OperatorID, NewHbar(1)).
		AddHbarTransfer(*account3, NewHbar(-1)).
		SetTransactionID(TransactionIDGenerate(*account3)).
		SetBatchKey(batchKey3).
		FreezeWith(env.Client)
	require.NoError(t, err)
	transfer3.Sign(key3)

	// Create and execute batch transaction
	batchTx, err := NewBatchTransaction().
		AddInnerTransaction(transfer1).
		AddInnerTransaction(transfer2).
		AddInnerTransaction(transfer3).
		FreezeWith(env.Client)
	require.NoError(t, err)

	batchTx.Sign(batchKey1)
	batchTx.Sign(batchKey2)
	batchTx.Sign(batchKey3)

	resp, err := batchTx.Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, receipt.Status)
}

func TestIntegrationBatchTransactionPartialFailure(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Get initial balance
	info, err := NewAccountInfoQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)
	initialBalance := info.Balance

	// Create transactions
	key1, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	tx1, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key1).
		SetInitialBalance(NewHbar(1)).
		Batchify(env.Client, env.OperatorKey)
	require.NoError(t, err)

	key2, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	tx2, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key2).
		SetInitialBalance(NewHbar(1)).
		Batchify(env.Client, env.OperatorKey)
	require.NoError(t, err)

	key3, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	tx3, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key3).
		SetReceiverSignatureRequired(true).
		SetInitialBalance(NewHbar(1)).
		Batchify(env.Client, env.OperatorKey)
	require.NoError(t, err)

	// Execute batch transaction
	batchTx := NewBatchTransaction().
		AddInnerTransaction(tx1).
		AddInnerTransaction(tx2).
		AddInnerTransaction(tx3)

	resp, err := batchTx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)

	// Check final balance
	info, err = NewAccountInfoQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)
	finalBalance := info.Balance

	assert.Less(t, finalBalance.AsTinybar(), initialBalance.AsTinybar())
}

func TestIntegrationBatchTransactionBatchifiedOutsideBatch(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx, err := NewTopicCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetTopicMemo("[e2e::TopicCreateTransaction]").
		Batchify(env.Client, key)
	require.NoError(t, err)

	_, err = tx.Execute(env.Client)
	assert.Equal(t, errBatchedAndNotBatchTransaction, err)
}
