//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationLambdaSStoreTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{Contract: 1})))

	account, accountKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.AddHook(*hookDetail)
	})
	require.NoError(t, err)

	frozenTxn, err := NewLambdaSStoreTransaction().
		SetHookId(*NewHookId().SetEntityId(*NewHookEntityId().SetAccountId(account)).SetHookId(1)).
		AddStorageUpdate(*NewLambdaStorageUpdate().SetStorageSlot(*NewLambdaStorageSlot().SetKey([]byte{0x01}).SetValue([]byte{0x02}))).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozenTxn.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationLambdaSStoreTransactionInvalidSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{Contract: 1}))).
		SetAdminKey(adminKey)

	account, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.AddHook(*hookDetail)
	})
	require.NoError(t, err)

	resp, err := NewLambdaSStoreTransaction().
		SetHookId(*NewHookId().SetEntityId(*NewHookEntityId().SetAccountId(account)).SetHookId(1)).
		AddStorageUpdate(*NewLambdaStorageUpdate().SetStorageSlot(*NewLambdaStorageSlot().SetKey([]byte{0x01}).SetValue([]byte{0x02}))).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationLambdaSStoreTransactionNonExistentHook(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewLambdaSStoreTransaction().
		SetHookId(*NewHookId().SetEntityId(*NewHookEntityId().SetAccountId(AccountID{Account: 1})).SetHookId(1)).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: HOOK_NOT_FOUND")
}

func TestIntegrationLambdaSStoreTransactionClearStorageUpdate(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{Contract: 1})))

	account, accountKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.AddHook(*hookDetail)
	})
	require.NoError(t, err)

	frozenTxn, err := NewLambdaSStoreTransaction().
		SetHookId(*NewHookId().SetEntityId(*NewHookEntityId().SetAccountId(account)).SetHookId(1)).
		AddStorageUpdate(*NewLambdaStorageUpdate().SetStorageSlot(*NewLambdaStorageSlot().SetKey([]byte{0x01}).SetValue([]byte{0x02}))).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozenTxn.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	frozenTxn, err = NewLambdaSStoreTransaction().
		SetHookId(*NewHookId().SetEntityId(*NewHookEntityId().SetAccountId(account)).SetHookId(1)).
		ClearStorageUpdate([]byte{0x01}).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenTxn.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationLambdaSStoreTransactionTooManyStorageUpdates(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{Contract: 1})))

	account, accountKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.AddHook(*hookDetail)
	})
	require.NoError(t, err)

	storageUpdates := make([]LambdaStorageUpdate, 11)
	for i := 0; i < 11; i++ {
		storageUpdates[i] = *NewLambdaStorageUpdate().SetStorageSlot(*NewLambdaStorageSlot().SetKey([]byte{byte(0x01)}).SetValue([]byte{byte(0x02)}))
	}

	frozenTxn, err := NewLambdaSStoreTransaction().
		SetHookId(*NewHookId().SetEntityId(*NewHookEntityId().SetAccountId(account)).SetHookId(1)).
		SetStorageUpdates(storageUpdates).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozenTxn.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: TOO_MANY_LAMBDA_STORAGE_UPDATES")
}
