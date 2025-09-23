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

func TestIntegrationAccountUpdateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetExpirationTime(time.Now().Add(time.Hour * 24 * 92)).
		SetKey(newKey2.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)
	tx.Sign(newKey2)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)

	require.NoError(t, err)

	txDelete.Sign(newKey2)

	resp, err = txDelete.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)

	require.NoError(t, err)

}

func TestIntegrationAccountUpdateTransactionNoSigning(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	_, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetKey(newKey2.PublicKey()).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, newKey.PublicKey().String(), info.Key.String())

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)

	require.NoError(t, err)

	txDelete.Sign(newKey)

	resp, err = txDelete.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountUpdateTransactionAccountIDNotSet(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewAccountUpdateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Contains(t, err.Error(), "exceptional precheck status ACCOUNT_ID_DOES_NOT_EXIST")
	}
}

func TestIntegrationAccountUpdateOptionalFields(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a new account
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	// Test setting the optional fields
	testMemo := "test memo"
	maxTokens := int32(100)
	declineReward := true

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountMemo(testMemo).
		SetMaxAutomaticTokenAssociations(maxTokens).
		SetDeclineStakingReward(declineReward).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the account was updated correctly
	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, testMemo, info.AccountMemo)
	assert.Equal(t, maxTokens, int32(info.MaxAutomaticTokenAssociations))
	assert.Equal(t, declineReward, info.StakingInfo.DeclineStakingReward)

	// Test setting to zero/empty values
	emptyMemo := ""
	zeroTokens := int32(0)
	acceptReward := false

	tx, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountMemo(emptyMemo).
		SetMaxAutomaticTokenAssociations(zeroTokens).
		SetDeclineStakingReward(acceptReward).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the account was updated correctly with zero/empty values
	info, err = NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, emptyMemo, info.AccountMemo)
	assert.Equal(t, zeroTokens, int32(info.MaxAutomaticTokenAssociations))
	assert.Equal(t, acceptReward, info.StakingInfo.DeclineStakingReward)

	// Clean up by deleting the test account
	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)
	require.NoError(t, err)

	txDelete.Sign(newKey)

	resp, err = txDelete.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationAccountUpdateSelectiveFieldChanges(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a new account
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	// Set initial values for all fields
	initialMemo := "initial memo"
	initialMaxTokens := int32(100)
	initialDeclineReward := false

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountMemo(initialMemo).
		SetMaxAutomaticTokenAssociations(initialMaxTokens).
		SetDeclineStakingReward(initialDeclineReward).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify initial values
	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, initialMemo, info.AccountMemo)
	assert.Equal(t, initialMaxTokens, int32(info.MaxAutomaticTokenAssociations))
	assert.Equal(t, initialDeclineReward, info.StakingInfo.DeclineStakingReward)

	// Test 1: Update ONLY memo, verify others unchanged
	updatedMemo := "updated memo"

	tx, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountMemo(updatedMemo).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify only memo changed
	info, err = NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, updatedMemo, info.AccountMemo)                               // Updated
	assert.Equal(t, initialMaxTokens, int32(info.MaxAutomaticTokenAssociations)) // Unchanged
	assert.Equal(t, initialDeclineReward, info.StakingInfo.DeclineStakingReward) // Unchanged

	// Test 2: Update ONLY maxAutomaticTokenAssociations, verify others unchanged
	updatedMaxTokens := int32(50)

	tx, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxAutomaticTokenAssociations(updatedMaxTokens).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify only maxAutomaticTokenAssociations changed
	info, err = NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, updatedMemo, info.AccountMemo)                               // Unchanged from last update
	assert.Equal(t, updatedMaxTokens, int32(info.MaxAutomaticTokenAssociations)) // Updated
	assert.Equal(t, initialDeclineReward, info.StakingInfo.DeclineStakingReward) // Unchanged

	// Test 3: Update ONLY declineReward, verify others unchanged
	updatedDeclineReward := true

	tx, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetDeclineStakingReward(updatedDeclineReward).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify only declineReward changed
	info, err = NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, updatedMemo, info.AccountMemo)                               // Unchanged from last update
	assert.Equal(t, updatedMaxTokens, int32(info.MaxAutomaticTokenAssociations)) // Unchanged from last update
	assert.Equal(t, updatedDeclineReward, info.StakingInfo.DeclineStakingReward) // Updated

	// Clean up by deleting the test account
	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)
	require.NoError(t, err)

	txDelete.Sign(newKey)

	resp, err = txDelete.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

// HIP-1195 hooks

func TestIntegrationAccountUpdateTransactionCanExecuteWithHook(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{})))

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		AddHook(*hookDetail).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationAccountUpdateTransactionAddDuplicateHook(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{})))

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetHooks([]HookCreationDetails{*hookDetail, *hookDetail}).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.ErrorContains(t, err, "exceptional precheck status HOOK_ID_REPEATED_IN_CREATION_DETAILS")
}

func TestIntegrationAccountUpdateTransactionAddExisingHook(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{})))

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		AddHook(*hookDetail).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		AddHook(*hookDetail).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: HOOK_ID_IN_USE")
}

func TestIntegrationAccountUpdateTransactionUpdateAddHookWithInitialStorageUpdates(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().
			SetStorageUpdates([]LambdaStorageUpdate{*NewLambdaStorageUpdate().SetStorageSlot(*NewLambdaStorageSlot().SetKey([]byte{0x01}).SetValue([]byte{0x02}))}).
			SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{})))

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		AddHook(*hookDetail).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationAccountUpdateTransactionCannotAddHookThatIsInUse(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{})))

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetHooks([]HookCreationDetails{*hookDetail}).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		AddHook(*hookDetail).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: HOOK_ID_IN_USE")
}

func TestIntegrationAccountUpdateTransactionCanDeleteHook(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{})))

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetHooks([]HookCreationDetails{*hookDetail}).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		DeleteHook(1).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationAccountUpdateTransactionCanotDeleteNonExistantHook(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{})))

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetHooks([]HookCreationDetails{*hookDetail}).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		DeleteHook(123).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: HOOK_NOT_FOUND")
}

func TestIntegrationAccountUpdateTransactionCanotAddAndDeleteHookAtTheSameTime(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{})))

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		DeleteHook(1).
		SetHooks([]HookCreationDetails{*hookDetail}).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: HOOK_NOT_FOUND")
}

func TestIntegrationAccountUpdateTransactionCanotDeleteDeletedHook(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(ContractID{})))

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		AddHook(*hookDetail).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	require.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		DeleteHook(1).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err = NewAccountUpdateTransaction().
		SetAccountID(accountID).
		DeleteHook(1).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx.Sign(newKey)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: HOOK_DELETED")
}
