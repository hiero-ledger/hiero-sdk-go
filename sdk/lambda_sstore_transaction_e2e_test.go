//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationLambdaSStoreUpdatesStorageWithValidSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	initialSlot := NewLambdaStorageSlot().
		SetKey([]byte{0x01}).
		SetValue([]byte{0x01})

	lambdaHook := NewLambdaEvmHook().
		SetContractId(hookContractId).
		AddStorageUpdate(initialSlot)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(3).
		SetLambdaEvmHook(*lambdaHook).
		SetAdminKey(accountKey.PublicKey())

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		AddHook(*hookDetails).
		SetInitialBalance(NewHbar(10)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	ownerId := *receipt.AccountID

	entityId := NewHookEntityIdWithAccountId(ownerId)
	hookId := NewHookId(*entityId, 3)

	update := NewLambdaStorageSlot().
		SetKey([]byte{0x01}).
		SetValue([]byte{0x02})

	frozenTxn, err := NewLambdaSStoreTransaction().
		SetHookId(*hookId).
		AddStorageUpdate(update).
		FreezeWith(env.Client)

	resp, err = frozenTxn.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, receipt.Status)
}

func TestIntegrationLambdaSStoreFailsWithoutProperSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	initialSlot := NewLambdaStorageSlot().
		SetKey([]byte{0x01}).
		SetValue([]byte{0x01})

	lambdaHook := NewLambdaEvmHook().
		SetContractId(hookContractId).
		AddStorageUpdate(initialSlot)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(3).
		SetLambdaEvmHook(*lambdaHook).
		SetAdminKey(accountKey.PublicKey())

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		AddHook(*hookDetails).
		SetInitialBalance(NewHbar(10)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	ownerId := *receipt.AccountID

	entityId := NewHookEntityIdWithAccountId(ownerId)
	hookId := NewHookId(*entityId, 3)

	update := NewLambdaStorageSlot().
		SetKey([]byte{0x01}).
		SetValue([]byte{0x02})

	unauthorizedKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	frozenTxn, err := NewLambdaSStoreTransaction().
		SetHookId(*hookId).
		AddStorageUpdate(update).
		FreezeWith(env.Client)

	resp, err = frozenTxn.Sign(unauthorizedKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "INVALID_SIGNATURE")
}

func TestIntegrationLambdaSStoreFailsWithNonExistentHookId(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	signerKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(signerKey).
		SetInitialBalance(NewHbar(10)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	signerId := *receipt.AccountID

	entityId := NewHookEntityIdWithAccountId(signerId)
	hookId := NewHookId(*entityId, 9999)

	update := NewLambdaStorageSlot().
		SetKey([]byte{0x0A}).
		SetValue([]byte{0x0B})

	frozenTxn, err := NewLambdaSStoreTransaction().
		SetHookId(*hookId).
		AddStorageUpdate(update).
		FreezeWith(env.Client)

	resp, err = frozenTxn.Sign(signerKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HOOK_NOT_FOUND")
}

func TestIntegrationLambdaSStoreTooManyStorageUpdatesFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	lambdaHook := NewLambdaEvmHook().
		SetContractId(hookContractId)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*lambdaHook).
		SetAdminKey(accountKey.PublicKey())

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		AddHook(*hookDetails).
		SetInitialBalance(NewHbar(10)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	ownerId := *receipt.AccountID

	entityId := NewHookEntityIdWithAccountId(ownerId)
	hookId := NewHookId(*entityId, 1)

	slot := NewLambdaStorageSlot().
		SetKey([]byte{0x01, 0x02, 0x03, 0x04}).
		SetValue([]byte{0x05, 0x06, 0x07, 0x08})

	var updates []LambdaStorageUpdate
	for i := 0; i < 256; i++ {
		updates = append(updates, slot)
	}

	frozenTxn, err := NewLambdaSStoreTransaction().
		SetHookId(*hookId).
		SetStorageUpdates(updates).
		FreezeWith(env.Client)

	resp, err = frozenTxn.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "TOO_MANY_LAMBDA_STORAGE_UPDATES")
}
