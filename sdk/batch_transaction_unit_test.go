//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	privateKeyED25519, _ = PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	privateKeyECDSA, _   = PrivateKeyFromStringECDSA("7f109a9e3b0d8ecfba9cc23a3614433ce0fa7ddcc80f2a8f10b222179a5a80d6")
	validStart           = time.Unix(1554158542, 0)
)

func spawnTestTransactionAccountCreate() *AccountCreateTransaction {
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	accountID3, _ := AccountIDFromString("0.0.3")
	transactionID := TransactionIDGenerate(accountID2)

	tx := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID).
		SetKeyWithoutAlias(privateKeyED25519).
		SetInitialBalance(NewHbar(450)).
		SetAccountMemo("some memo").
		SetReceiverSignatureRequired(true).
		SetAutoRenewPeriod(10 * time.Hour).
		SetStakedAccountID(accountID3).
		SetMaxAutomaticTokenAssociations(100).
		SetMaxTransactionFee(NewHbar(100000)).
		SetBatchKey(privateKeyECDSA)

	tx2, err := tx.Freeze()
	if err != nil {
		panic(err)
	}
	tx2.Sign(privateKeyED25519)

	return tx2
}

func spawnTestBatchTransaction() *BatchTransaction {
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)
	batchKey, err := PrivateKeyGenerateEd25519()
	if err != nil {
		panic(err)
	}

	tx := NewBatchTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	innerTx1 := spawnTestTransactionAccountCreate()
	innerTx2 := spawnTestTransactionAccountCreate()
	innerTx3 := spawnTestTransactionAccountCreate()

	tx.AddInnerTransaction(innerTx1)
	tx.AddInnerTransaction(innerTx2)
	tx.AddInnerTransaction(innerTx3)

	tx2, err := tx.Freeze()
	if err != nil {
		panic(err)
	}
	tx2.Sign(batchKey)

	return tx2
}

func TestUnitBatchTransactionToFromBytes(t *testing.T) {
	t.Parallel()

	tx := spawnTestBatchTransaction()
	txBytes, err := tx.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	assert.Equal(t, tx.buildProtoBody(), txFromBytes.(BatchTransaction).buildProtoBody())
}

func TestUnitBatchTransactionToFromBytesNoSetters(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	txBytes, err := tx.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	assert.Equal(t, tx.buildProtoBody(), txFromBytes.(BatchTransaction).buildProtoBody())
}

func TestUnitBatchTransactionGetInnerTransactions(t *testing.T) {
	t.Parallel()

	tx := spawnTestBatchTransaction()
	innerTxs := tx.GetInnerTransactions()

	assert.NotNil(t, innerTxs)
	assert.Len(t, innerTxs, 3)
}

func TestUnitBatchTransactionSetInnerTransactions(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	newInnerTxs := []TransactionInterface{
		spawnTestTransactionAccountCreate(),
		spawnTestTransactionAccountCreate(),
	}

	tx2 := tx.SetInnerTransactions(newInnerTxs)
	assert.Nil(t, tx2.freezeError)

	assert.NotNil(t, tx2.GetInnerTransactions())
	assert.Len(t, tx2.GetInnerTransactions(), 2)
	assert.Equal(t, newInnerTxs, tx2.GetInnerTransactions())
}

func TestUnitBatchTransactionAddInnerTransaction(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	newTx := spawnTestTransactionAccountCreate()

	tx.AddInnerTransaction(newTx)

	assert.NotNil(t, tx.GetInnerTransactions())
	assert.Len(t, tx.GetInnerTransactions(), 1)
	assert.Contains(t, tx.GetInnerTransactions(), newTx)
}

func TestUnitBatchTransactionGetInnerTransactionIDs(t *testing.T) {
	t.Parallel()

	tx := spawnTestBatchTransaction()
	accountID, _ := AccountIDFromString("0.0.5006")
	expectedID := TransactionIDGenerate(accountID)

	ids := tx.GetInnerTransactionIDs()

	assert.NotNil(t, ids)
	assert.Len(t, ids, 3)
	for _, id := range ids {
		assert.Equal(t, expectedID.AccountID, id.AccountID)
	}
}

func TestUnitBatchTransactionChainedSetters(t *testing.T) {
	t.Parallel()

	accountID, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	tx := NewBatchTransaction().
		SetNodeAccountIDs([]AccountID{accountID}).
		SetTransactionID(transactionID).
		AddInnerTransaction(spawnTestTransactionAccountCreate())

	tx2, err := tx.Freeze()
	require.NoError(t, err)

	assert.Len(t, tx2.GetInnerTransactions(), 1)
	assert.Len(t, tx2.GetNodeAccountIDs(), 1)
	assert.NotNil(t, tx2.GetTransactionID())
}

func TestUnitBatchTransactionRejectFreezeTransaction(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	freezeTx := NewFreezeTransaction().
		SetStartTime(time.Now()).
		SetFreezeType(FreezeTypeFreezeOnly).
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	freezeTx2, err := freezeTx.Freeze()
	require.NoError(t, err)

	tx2 := tx.AddInnerTransaction(freezeTx2)
	assert.NotNil(t, tx2.freezeError)
	assert.Equal(t, errTransactionTypeNotAllowed, tx2.freezeError)
}

func TestUnitBatchTransactionRejectBatchTransaction(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	innerBatchTx := NewBatchTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs([]AccountID{accountID1, accountID2})

	innerBatchTx2, _ := innerBatchTx.Freeze()

	tx2 := tx.AddInnerTransaction(innerBatchTx2)
	assert.NotNil(t, tx2.freezeError)
	assert.Equal(t, errTransactionTypeNotAllowed, tx2.freezeError)
}

func TestUnitBatchTransactionRejectBlacklistedTransactionInList(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	validTx := spawnTestTransactionAccountCreate()
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	freezeTx := NewFreezeTransaction().
		SetStartTime(time.Now()).
		SetFreezeType(FreezeTypeFreezeOnly).
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	freezeTx2, err := freezeTx.Freeze()
	require.NoError(t, err)

	tx.SetInnerTransactions([]TransactionInterface{validTx, freezeTx2})
	assert.NotNil(t, tx.freezeError)
	assert.Equal(t, errTransactionTypeNotAllowed, tx.freezeError)
}

func TestUnitBatchTransactionRejectNullTransaction(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction().AddInnerTransaction(nil)
	assert.NotNil(t, tx.freezeError)
	assert.Equal(t, errInnerTransactionNil, tx.freezeError)
}

func TestUnitBatchTransactionRejectUnfrozenTransaction(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	unfrozenTx := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	tx2 := tx.AddInnerTransaction(unfrozenTx)
	assert.NotNil(t, tx2.freezeError)
	assert.Equal(t, errInnerTransactionShouldBeFrozen, tx2.freezeError)
}

func TestUnitBatchTransactionRejectTransactionAfterFreeze(t *testing.T) {
	t.Parallel()

	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	tx := NewBatchTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	tx2, _ := tx.Freeze()
	tx3 := tx2.AddInnerTransaction(spawnTestTransactionAccountCreate())

	assert.NotNil(t, tx3.freezeError)
	assert.Equal(t, errTransactionIsFrozen, tx3.freezeError)
}

func TestUnitBatchTransactionRejectTransactionListAfterFreeze(t *testing.T) {
	t.Parallel()

	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	tx := NewBatchTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	tx2, _ := tx.Freeze()
	tx2.SetInnerTransactions([]TransactionInterface{spawnTestTransactionAccountCreate()})

	assert.NotNil(t, tx.freezeError)
}

func TestUnitBatchTransactionAllowEmptyTransactionList(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	tx2 := tx.SetInnerTransactions([]TransactionInterface{})
	assert.Nil(t, tx2.freezeError)

	assert.NotNil(t, tx2.GetInnerTransactions())
	assert.Empty(t, tx2.GetInnerTransactions())
}

func TestUnitBatchTransactionPreserveTransactionOrder(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	tx1 := spawnTestTransactionAccountCreate()
	tx2 := spawnTestTransactionAccountCreate()
	tx3 := spawnTestTransactionAccountCreate()

	transactions := []TransactionInterface{tx1, tx2, tx3}
	batchTx := tx.SetInnerTransactions(transactions)
	assert.Nil(t, batchTx.freezeError)

	assert.Equal(t, []TransactionInterface{tx1, tx2, tx3}, batchTx.GetInnerTransactions())
}

func TestUnitBatchTransactionRejectTransactionWithoutBatchKey(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	txWithoutBatchKey := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	txWithoutBatchKey2, _ := txWithoutBatchKey.Freeze()

	tx2 := tx.AddInnerTransaction(txWithoutBatchKey2)
	assert.NotNil(t, tx2.freezeError)
	assert.Equal(t, errBatchKeyNotSet, tx2.freezeError)
}

func TestUnitBatchTransactionValidateAllTransactionsInList(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	validTx := spawnTestTransactionAccountCreate()
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	txWithoutBatchKey := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	txWithoutBatchKey2, _ := txWithoutBatchKey.Freeze()

	tx.SetInnerTransactions([]TransactionInterface{validTx, txWithoutBatchKey2})
	assert.NotNil(t, tx.freezeError)
	assert.Equal(t, errBatchKeyNotSet, tx.freezeError)
}

func TestUnitBatchTransactionValidateMultipleConditions(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	// Test unfrozen transaction with no batch key
	unfrozenTxWithoutBatchKey := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	tx2 := tx.AddInnerTransaction(unfrozenTxWithoutBatchKey)
	assert.NotNil(t, tx2.freezeError)
	assert.Equal(t, errInnerTransactionShouldBeFrozen, tx2.freezeError)

	// Test frozen transaction with no batch key
	frozenTxWithoutBatchKey, err := unfrozenTxWithoutBatchKey.Freeze()
	require.NoError(t, err)
	tx3 := tx.AddInnerTransaction(frozenTxWithoutBatchKey)
	assert.NotNil(t, tx3.freezeError)
	assert.Equal(t, errBatchKeyNotSet, tx3.freezeError)

	// Test blacklisted transaction with batch key
	blacklistedTx := NewFreezeTransaction().
		SetStartTime(time.Now()).
		SetFreezeType(FreezeTypeFreezeOnly).
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID).
		SetBatchKey(privateKeyECDSA)
	require.NoError(t, err)

	blacklistedTx2, err := blacklistedTx.Freeze()
	require.NoError(t, err)
	tx4 := tx.AddInnerTransaction(blacklistedTx2)
	assert.NotNil(t, tx4.freezeError)
	assert.Equal(t, errTransactionTypeNotAllowed, tx4.freezeError)
}

func TestUnitBatchTransactionValidateTransactionStateInOrder(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	transaction := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID)

	// First check should be for frozen state
	tx2 := tx.AddInnerTransaction(transaction)
	assert.NotNil(t, tx2.freezeError)
	assert.Equal(t, errInnerTransactionShouldBeFrozen, tx2.freezeError)

	// After freezing, next check should be for batch key
	frozenTx, _ := transaction.Freeze()
	tx3 := tx.AddInnerTransaction(frozenTx)
	assert.NotNil(t, tx3.freezeError)
	assert.Equal(t, errBatchKeyNotSet, tx3.freezeError)
}

func TestUnitBatchTransactionAcceptValidTransaction(t *testing.T) {
	t.Parallel()

	tx := NewBatchTransaction()
	accountID1, _ := AccountIDFromString("0.0.5005")
	accountID2, _ := AccountIDFromString("0.0.5006")
	transactionID := TransactionIDGenerate(accountID2)

	validTx := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{accountID1, accountID2}).
		SetTransactionID(transactionID).
		SetBatchKey(privateKeyECDSA)

	validTx2, err := validTx.Freeze()
	require.NoError(t, err)
	tx2 := tx.AddInnerTransaction(validTx2)

	assert.Nil(t, tx2.freezeError)
	assert.NotNil(t, tx2.GetInnerTransactions())
	assert.Len(t, tx2.GetInnerTransactions(), 1)
	assert.Contains(t, tx2.GetInnerTransactions(), validTx2)
}
